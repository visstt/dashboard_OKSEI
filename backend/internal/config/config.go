package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Config содержит конфигурацию приложения
type Config struct {
	// Интервал обновления данных
	RefreshInterval time.Duration

	// Пути к файлам
	ProjectRoot      string
	AttendanceInput  string
	AttendanceOutput string
	StatementInput   string
	StatementOutput  string
	PythonScript     string

	// Настройки сервера
	ServerPort string
	ServerHost string

	// Настройки БД
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	// JWT авторизация (из attendance-backend)
	JWTSecret string

	// CORS (из attendance-backend)
	CORSOrigins []string

	// Алерты (из attendance-backend)
	AbsenceThreshold int

	// Логин (из attendance-backend)
	LoginUser     string
	LoginPassword string
	LoginRole     string
}

// Load загружает конфигурацию из переменных окружения или использует значения по умолчанию
func Load() (*Config, error) {
	// Получаем корневую директорию проекта
	// Используем рабочую директорию при запуске сервера
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения рабочей директории: %v", err)
	}

	var projectRoot string
	// Если запускаем из backend/, поднимаемся на уровень выше
	if filepath.Base(wd) == "backend" {
		projectRoot = filepath.Dir(wd)
	} else {
		// Ищем директорию с папкой public/ и backend/
		current := wd
		for {
			// Проверяем наличие папок public/ и backend/ в текущей директории
			publicExists := false
			backendExists := false
			
			if _, err := os.Stat(filepath.Join(current, "public")); err == nil {
				publicExists = true
			}
			if _, err := os.Stat(filepath.Join(current, "backend")); err == nil {
				backendExists = true
			}
			
			// Если есть обе папки - это корень проекта
			if publicExists && backendExists {
				projectRoot = current
				break
			}
			
			// Поднимаемся на уровень выше
			parent := filepath.Dir(current)
			if parent == current || parent == "/" {
				// Дошли до корня, не нашли
				break
			}
			current = parent
		}
		
		// Если не нашли, пробуем найти по наличию public/
		if projectRoot == "" {
			current := wd
			for {
				if _, err := os.Stat(filepath.Join(current, "public")); err == nil {
					projectRoot = current
					break
				}
				parent := filepath.Dir(current)
				if parent == current || parent == "/" {
					break
				}
				current = parent
			}
		}
		
		// Если всё ещё не нашли, используем текущую директорию
		if projectRoot == "" {
			projectRoot = wd
		}
	}

	// Проверяем, что директория существует и содержит public/
	if _, err := os.Stat(filepath.Join(projectRoot, "public")); os.IsNotExist(err) {
		return nil, fmt.Errorf("не найдена директория проекта (ожидается папка 'public' в %s)", projectRoot)
	}

	// Интервал обновления (по умолчанию 90 минут)
	refreshInterval := 90 * time.Minute
	if intervalStr := os.Getenv("REFRESH_INTERVAL"); intervalStr != "" {
		if parsed, err := time.ParseDuration(intervalStr); err == nil {
			refreshInterval = parsed
		}
	}

	// Порт сервера (по умолчанию 8080)
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// Хост сервера (по умолчанию localhost)
	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		serverHost = "localhost"
	}

	// Настройки БД
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Формируем URL из отдельных параметров, если DATABASE_URL не указан
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost"
		}
		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "5432"
		}
		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "postgres"
		}
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "dashboard"
		}

		if dbPassword != "" {
			databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				dbUser, dbPassword, dbHost, dbPort, dbName)
		}
	}

	// JWT Secret (из attendance-backend)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-me-in-production"
	}

	// CORS Origins (из attendance-backend)
	corsEnv := strings.TrimSpace(os.Getenv("CORS_ORIGINS"))
	var corsOrigins []string
	if corsEnv == "" || corsEnv == "*" {
		corsOrigins = []string{"*"}
	} else {
		parts := strings.Split(corsEnv, ",")
		trimmed := make([]string, 0, len(parts))
		for _, o := range parts {
			s := strings.TrimSpace(o)
			if s != "" {
				trimmed = append(trimmed, s)
			}
		}
		if len(trimmed) == 0 {
			trimmed = []string{"http://localhost:3000", "http://localhost:5173"}
		}
		corsOrigins = trimmed
	}

	// Absence Threshold (из attendance-backend)
	threshold, _ := strconv.Atoi(os.Getenv("ABSENCE_THRESHOLD"))
	if threshold <= 0 || threshold > 100 {
		threshold = 10
	}

	// Login credentials (из attendance-backend)
	loginUser := os.Getenv("LOGIN_USER")
	if loginUser == "" {
		loginUser = "admin"
	}
	loginPassword := os.Getenv("LOGIN_PASSWORD")
	if loginPassword == "" {
		loginPassword = "admin"
	}
	loginRole := os.Getenv("LOGIN_ROLE")
	if loginRole == "" {
		loginRole = "admin"
	}

	cfg := &Config{
		RefreshInterval:  refreshInterval,
		ProjectRoot:      projectRoot,
		AttendanceInput:  filepath.Join(projectRoot, "Посещаемость.xlsx"),
		AttendanceOutput: filepath.Join(projectRoot, "public", "attendance.json"),
		StatementInput:   filepath.Join(projectRoot, "ведомость.xls"),
		StatementOutput:  filepath.Join(projectRoot, "public", "summary.json"),
		PythonScript:     filepath.Join(projectRoot, "statement-converter", "xls_to_xlsx.py"),
		ServerPort:       serverPort,
		ServerHost:       serverHost,
		DatabaseURL:      databaseURL,
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabaseUser:     os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName:     os.Getenv("DB_NAME"),
		JWTSecret:        jwtSecret,
		CORSOrigins:      corsOrigins,
		AbsenceThreshold: threshold,
		LoginUser:        loginUser,
		LoginPassword:    loginPassword,
		LoginRole:        loginRole,
	}

	return cfg, nil
}
