package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"dashboard/internal/database"
	"dashboard/internal/scheduler"
)

// GinHandler содержит обработчики API для Gin
type GinHandler struct {
	scheduler         *scheduler.Scheduler
	dbLoader          *database.Loader
	lastRefresh       time.Time
	refreshInProgress bool
}

func NewGinHandler(scheduler *scheduler.Scheduler, dbLoader *database.Loader) *GinHandler {
	return &GinHandler{
		scheduler: scheduler,
		dbLoader:  dbLoader,
	}
}

// RefreshData запускает ручное обновление данных
// @Summary Ручное обновление данных
// @Description Запускает конвертацию Excel файлов в JSON и загрузку в БД
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Данные успешно обновлены"
// @Failure 409 {object} map[string]string "Обновление уже выполняется"
// @Failure 500 {object} map[string]string "Ошибка обновления данных"
// @Router /admin/refresh-data [post]
func (h *GinHandler) RefreshData(c *gin.Context) {
	// Проверяем, не выполняется ли уже обновление
	if h.refreshInProgress {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Обновление уже выполняется",
		})
		return
	}

	h.refreshInProgress = true
	defer func() {
		h.refreshInProgress = false
	}()

	log.Println("[API] Запуск ручного обновления данных...")

	// Запускаем обновление
	if err := h.scheduler.RefreshData(); err != nil {
		log.Printf("[API] Ошибка обновления данных: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка обновления данных",
			"details": err.Error(),
		})
		return
	}

	// Загружаем в БД (пути из конфига)
	attendancePath := c.GetString("attendance_output")
	statementPath := c.GetString("statement_output")
	
	if attendancePath != "" {
		if err := h.dbLoader.LoadAttendance(attendancePath); err != nil {
			log.Printf("[API] Предупреждение при загрузке посещаемости в БД: %v", err)
		}
	}
	if statementPath != "" {
		if err := h.dbLoader.LoadStatement(statementPath); err != nil {
			log.Printf("[API] Предупреждение при загрузке ведомости в БД: %v", err)
		}
	}

	h.lastRefresh = time.Now()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Данные успешно обновлены",
		"time":    h.lastRefresh.Format(time.RFC3339),
	})
}

// GetRefreshStatus возвращает статус последнего обновления
// @Summary Статус обновления данных
// @Description Возвращает информацию о последнем обновлении данных
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{} "Статус обновления"
// @Router /admin/refresh-status [get]
func (h *GinHandler) GetRefreshStatus(c *gin.Context) {
	status := gin.H{
		"in_progress": h.refreshInProgress,
	}

	if !h.lastRefresh.IsZero() {
		status["last_refresh"] = h.lastRefresh.Format(time.RFC3339)
		status["last_refresh_ago"] = time.Since(h.lastRefresh).String()
	} else {
		status["last_refresh"] = nil
		status["last_refresh_ago"] = nil
	}

	c.JSON(http.StatusOK, status)
}

// GetRefreshHistory возвращает историю обновлений
// @Summary История обновлений
// @Description Возвращает историю обновлений данных
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{} "История обновлений"
// @Router /admin/refresh-history [get]
func (h *GinHandler) GetRefreshHistory(c *gin.Context) {
	// TODO: Реальная история из БД
	history := []gin.H{
		{
			"time":    h.lastRefresh.Format(time.RFC3339),
			"status":  "success",
			"message": "Обновление выполнено успешно",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}

// HealthCheck проверяет работоспособность сервера
// @Summary Health Check
// @Description Проверяет работоспособность сервера
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "Сервер работает"
// @Router /health [get]
func (h *GinHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "dashboard-backend",
	})
}
