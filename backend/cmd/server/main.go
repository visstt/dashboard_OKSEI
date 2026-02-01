package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dashboard/internal/api"
	"dashboard/internal/config"
	"dashboard/internal/database"
	"dashboard/internal/middleware"
	"dashboard/internal/scheduler"
	"dashboard/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// @title Dashboard Backend API
// @version 1.0
// @description API для управления дашбордом посещаемости студентов
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@dashboard.local

// @host localhost:8080
// @BasePath /api
func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Server] Ошибка загрузки конфигурации: %v", err)
	}

	log.Printf("[Server] Запуск бэкенд сервера...")
	log.Printf("[Server] Корневая директория: %s", cfg.ProjectRoot)
	log.Printf("[Server] Интервал обновления: %v", cfg.RefreshInterval)

	// Устанавливаем режим работы Gin (release для продакшена)
	gin.SetMode(gin.ReleaseMode)

	// Подключаемся к БД
	if cfg.DatabaseURL != "" {
		if err := database.Connect(cfg.DatabaseURL); err != nil {
			log.Printf("[Server] Предупреждение: не удалось подключиться к БД: %v", err)
			log.Println("[Server] Продолжаем работу без БД (данные не будут сохраняться)")
		} else {
			// Инициализируем схему БД
			if err := database.InitSchema(); err != nil {
				log.Printf("[Server] Предупреждение: не удалось инициализировать схему БД: %v", err)
			}
			// Закрываем подключение при завершении
			defer database.Close()
		}
	} else {
		log.Println("[Server] DATABASE_URL не указан, работаем без БД")
	}

	// Инициализируем загрузчик БД
	dbLoader := database.NewLoader(database.DB)

	// Инициализируем планировщик
	sched := scheduler.NewScheduler(
		cfg.ProjectRoot,
		cfg.AttendanceInput,
		cfg.AttendanceOutput,
		cfg.StatementInput,
		cfg.StatementOutput,
		cfg.PythonScript,
	)

	// Инициализируем сервисы
	attendanceService := services.NewAttendanceService(cfg.AttendanceOutput)

	// Инициализируем handlers
	ginHandler := api.NewGinHandler(sched, dbLoader)
	authHandler := api.NewAuthHandler(cfg)
	dashboardHandler := api.NewDashboardHandler(attendanceService, cfg.AbsenceThreshold)

	// Настраиваем Gin router (используем gin.New() вместо gin.Default() чтобы избежать дублирования middleware)
	router := gin.New()

	// Подключаем middleware
	router.Use(middleware.SetupCORS())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Передаём пути в контекст для handlers
	router.Use(func(c *gin.Context) {
		c.Set("attendance_output", cfg.AttendanceOutput)
		c.Set("statement_output", cfg.StatementOutput)
		c.Next()
	})

	// API эндпоинты
	apiGroup := router.Group("/api")
	{
		// Публичные эндпоинты (без авторизации)
		apiGroup.POST("/login", authHandler.Login)
		apiGroup.GET("/health", ginHandler.HealthCheck)

		// Защищённые эндпоинты (требуют JWT)
		protected := apiGroup.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			// Эндпоинты дашборда (доступны всем авторизованным)
			protected.GET("/attendance", dashboardHandler.List)
			protected.GET("/attendance/summary", dashboardHandler.Summary)
			protected.GET("/attendance/drill/departments", dashboardHandler.DrillDepartments)
			protected.GET("/attendance/drill/groups", dashboardHandler.DrillGroups)
			protected.GET("/attendance/drill/students", dashboardHandler.DrillStudents)

			// Админские эндпоинты (только для admin)
			adminGroup := protected.Group("/admin")
			adminGroup.Use(middleware.RequireRole("admin"))
			{
				adminGroup.POST("/refresh-data", ginHandler.RefreshData)
				adminGroup.GET("/refresh-status", ginHandler.GetRefreshStatus)
				adminGroup.GET("/refresh-history", ginHandler.GetRefreshHistory)
			}
		}
	}

	// Swagger документация (упрощённый вариант через CDN)
	router.GET("/swagger/*path", serveSwagger)

	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		log.Printf("[Server] HTTP сервер запущен на http://%s", serverAddr)
		log.Printf("[Server] Swagger UI доступен на http://%s/swagger/", serverAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server] Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Настраиваем cron задачу
	c := cron.New()

	// Формируем cron выражение из интервала
	cronExpr := formatCronInterval(cfg.RefreshInterval)
	_, err = c.AddFunc(cronExpr, func() {
		log.Println("[Server] Запуск автоматического обновления данных...")
		if err := sched.RefreshData(); err != nil {
			log.Printf("[Server] Ошибка обновления данных: %v", err)
		}
		// После обновления загружаем в БД
		if database.DB != nil {
			if err := dbLoader.LoadAttendance(cfg.AttendanceOutput); err != nil {
				log.Printf("[Server] Предупреждение при загрузке посещаемости в БД: %v", err)
			}
			if err := dbLoader.LoadStatement(cfg.StatementOutput); err != nil {
				log.Printf("[Server] Предупреждение при загрузке ведомости в БД: %v", err)
			}
		}
	})
	if err != nil {
		log.Fatalf("[Server] Ошибка настройки cron: %v", err)
	}

	// Запускаем обновление сразу при старте
	log.Println("[Server] Первоначальное обновление данных...")
	if err := sched.RefreshData(); err != nil {
		log.Printf("[Server] Предупреждение при первоначальном обновлении: %v", err)
	}

	// Загружаем в БД после первоначального обновления
	if database.DB != nil {
		if err := dbLoader.LoadAttendance(cfg.AttendanceOutput); err != nil {
			log.Printf("[Server] Предупреждение при загрузке посещаемости в БД: %v", err)
		}
		if err := dbLoader.LoadStatement(cfg.StatementOutput); err != nil {
			log.Printf("[Server] Предупреждение при загрузке ведомости в БД: %v", err)
		}
	}

	// Запускаем планировщик
	c.Start()
	log.Printf("[Server] Планировщик запущен. Обновление данных каждые %v.", cfg.RefreshInterval)
	log.Println("[Server] Нажмите Ctrl+C для остановки...")

	// Обработка сигналов для корректного завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Блокируем выполнение до получения сигнала
	<-sigChan
	log.Println("[Server] Получен сигнал завершения. Остановка сервера...")
	c.Stop()

	// Останавливаем HTTP сервер
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("[Server] Ошибка остановки HTTP сервера: %v", err)
	}

	log.Println("[Server] Сервер остановлен.")
}

// formatCronInterval преобразует time.Duration в cron выражение
// Для упрощения используем @every, но можно сделать более точное расписание
func formatCronInterval(d time.Duration) string {
	// robfig/cron поддерживает @every напрямую
	minutes := int(d.Minutes())
	if minutes < 60 {
		return fmt.Sprintf("@every %dm", minutes)
	}
	hours := minutes / 60
	if hours*60 == minutes {
		return fmt.Sprintf("@every %dh", hours)
	}
	// Если не кратно часам, используем минуты
	return fmt.Sprintf("@every %dm", minutes)
}

// serveSwagger обрабатывает запросы к Swagger UI и JSON
func serveSwagger(c *gin.Context) {
	path := c.Param("path")
	
	// Если запрос к /swagger/doc.json - возвращаем JSON
	if path == "/doc.json" {
		c.Header("Content-Type", "application/json")
	swaggerJSON := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Dashboard Backend API",
			"version": "1.0",
			"description": "API для управления дашбордом посещаемости студентов",
			"contact": {
				"name": "API Support",
				"email": "support@dashboard.local"
			}
		},
		"servers": [
			{
				"url": "http://localhost:8080/api",
				"description": "Локальный сервер"
			}
		],
		"tags": [
			{
				"name": "admin",
				"description": "Административные операции"
			},
			{
				"name": "system",
				"description": "Системные операции"
			}
		],
		"paths": {
			"/admin/refresh-data": {
				"post": {
					"tags": ["admin"],
					"summary": "Ручное обновление данных",
					"description": "Запускает конвертацию Excel файлов в JSON и загрузку в БД",
					"responses": {
						"200": {
							"description": "Данные успешно обновлены",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"status": {"type": "string", "example": "success"},
											"message": {"type": "string", "example": "Данные успешно обновлены"},
											"time": {"type": "string", "format": "date-time"}
										}
									}
								}
							}
						},
						"409": {
							"description": "Обновление уже выполняется",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"error": {"type": "string", "example": "Обновление уже выполняется"}
										}
									}
								}
							}
						},
						"500": {
							"description": "Ошибка обновления данных",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"error": {"type": "string"},
											"details": {"type": "string"}
										}
									}
								}
							}
						}
					}
				}
			},
			"/admin/refresh-status": {
				"get": {
					"tags": ["admin"],
					"summary": "Статус обновления данных",
					"description": "Возвращает информацию о последнем обновлении данных",
					"responses": {
						"200": {
							"description": "Статус обновления",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"in_progress": {"type": "boolean", "example": false},
											"last_refresh": {"type": "string", "format": "date-time", "nullable": true},
											"last_refresh_ago": {"type": "string", "nullable": true}
										}
									}
								}
							}
						}
					}
				}
			},
			"/admin/refresh-history": {
				"get": {
					"tags": ["admin"],
					"summary": "История обновлений",
					"description": "Возвращает историю обновлений данных",
					"responses": {
						"200": {
							"description": "История обновлений",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"history": {
												"type": "array",
												"items": {
													"type": "object",
													"properties": {
														"time": {"type": "string", "format": "date-time"},
														"status": {"type": "string"},
														"message": {"type": "string"}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			},
			"/health": {
				"get": {
					"tags": ["system"],
					"summary": "Health Check",
					"description": "Проверяет работоспособность сервера",
					"responses": {
						"200": {
							"description": "Сервер работает",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"status": {"type": "string", "example": "ok"},
											"service": {"type": "string", "example": "dashboard-backend"}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`
		c.Data(http.StatusOK, "application/json", []byte(swaggerJSON))
		return
	}
	
	// Иначе возвращаем HTML страницу Swagger UI
	if path == "" || path == "/" {
		html := `<!DOCTYPE html>
<html>
<head>
	<title>Dashboard API - Swagger UI</title>
	<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
	<style>
		html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
		*, *:before, *:after { box-sizing: inherit; }
		body { margin:0; background: #fafafa; }
	</style>
</head>
<body>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
	<script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
	<script>
		window.onload = function() {
			const ui = SwaggerUIBundle({
				url: "/swagger/doc.json",
				dom_id: '#swagger-ui',
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIStandalonePreset
				],
				layout: "StandaloneLayout"
			});
		};
	</script>
</body>
</html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	} else {
		c.Status(http.StatusNotFound)
	}
}
