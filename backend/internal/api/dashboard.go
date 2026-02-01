package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"dashboard/internal/services"
)

// DashboardHandler обрабатывает запросы дашборда
type DashboardHandler struct {
	attendanceService *services.AttendanceService
	alertsThreshold  int
}

// NewDashboardHandler создаёт новый handler дашборда
func NewDashboardHandler(attendanceService *services.AttendanceService, alertsThreshold int) *DashboardHandler {
	return &DashboardHandler{
		attendanceService: attendanceService,
		alertsThreshold:   alertsThreshold,
	}
}

// List возвращает список записей посещаемости с фильтрацией
// GET /api/attendance
func (h *DashboardHandler) List(c *gin.Context) {
	_, flat, err := h.attendanceService.LoadFromJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot load attendance"})
		return
	}

	params := services.ParseFilterParams(c.Request)
	filtered := h.attendanceService.Filter(flat, params)

	// Проверяем алерты
	services.CheckAlerts(filtered, h.alertsThreshold)

	c.JSON(http.StatusOK, filtered)
}

// Summary возвращает сводку по посещаемости
// GET /api/attendance/summary
func (h *DashboardHandler) Summary(c *gin.Context) {
	departments, flat, err := h.attendanceService.LoadFromJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot load attendance"})
		return
	}

	params := services.ParseFilterParams(c.Request)
	filtered := h.attendanceService.Filter(flat, params)
	summary := h.attendanceService.BuildSummary(departments, filtered)

	c.JSON(http.StatusOK, summary)
}

// DrillDepartments возвращает drill-down по отделениям
// GET /api/attendance/drill/departments
func (h *DashboardHandler) DrillDepartments(c *gin.Context) {
	departments, flat, err := h.attendanceService.LoadFromJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot load attendance"})
		return
	}

	params := services.ParseFilterParams(c.Request)
	filtered := h.attendanceService.Filter(flat, params)
	result := h.attendanceService.BuildDrillDepartments(departments, filtered)

	c.JSON(http.StatusOK, result)
}

// DrillGroups возвращает drill-down по группам
// GET /api/attendance/drill/groups?department=...
func (h *DashboardHandler) DrillGroups(c *gin.Context) {
	department := strings.TrimSpace(c.Query("department"))
	if department == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "department parameter required"})
		return
	}

	departments, flat, err := h.attendanceService.LoadFromJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot load attendance"})
		return
	}

	params := services.ParseFilterParams(c.Request)
	filtered := h.attendanceService.Filter(flat, params)
	result := h.attendanceService.BuildDrillGroups(departments, filtered, department)

	c.JSON(http.StatusOK, result)
}

// DrillStudents возвращает drill-down по студентам
// GET /api/attendance/drill/students?department=...&group=...
func (h *DashboardHandler) DrillStudents(c *gin.Context) {
	department := strings.TrimSpace(c.Query("department"))
	group := strings.TrimSpace(c.Query("group"))

	if department == "" || group == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "department and group parameters required"})
		return
	}

	_, flat, err := h.attendanceService.LoadFromJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot load attendance"})
		return
	}

	params := services.ParseFilterParams(c.Request)
	filtered := h.attendanceService.Filter(flat, params)
	result := h.attendanceService.BuildDrillStudents(filtered, department, group)

	c.JSON(http.StatusOK, result)
}
