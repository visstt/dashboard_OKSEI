package services

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"dashboard/internal/models"
)

// AttendanceService предоставляет бизнес-логику для работы с посещаемостью
type AttendanceService struct {
	attendancePath string
}

// NewAttendanceService создаёт новый сервис
func NewAttendanceService(attendancePath string) *AttendanceService {
	return &AttendanceService{
		attendancePath: attendancePath,
	}
}

// LoadFromJSON загружает данные из JSON файла
func (s *AttendanceService) LoadFromJSON() ([]models.DepartmentJSON, []models.FlatRecord, error) {
	raw, err := os.ReadFile(s.attendancePath)
	if err != nil {
		return nil, nil, err
	}

	var departments []models.DepartmentJSON
	if err := json.Unmarshal(raw, &departments); err != nil {
		return nil, nil, err
	}

	flat := models.Flatten(departments)
	return departments, flat, nil
}

// FilterParams параметры фильтрации
type FilterParams struct {
	Department string
	Group       string
	Student     string
	Date        string
	DateFrom    string
	DateTo      string
	Period      string
	Search      string
	MissedMin   int
}

// Filter фильтрует записи по параметрам
func (s *AttendanceService) Filter(records []models.FlatRecord, params FilterParams) []models.FlatRecord {
	today := time.Now().Format("2006-01-02")
	var from, to string

	// Определяем период
	switch params.Period {
	case "7d":
		from = todayAdd(-7)
		to = today
	case "30d":
		from = todayAdd(-30)
		to = today
	case "90d":
		from = todayAdd(-90)
		to = today
	default:
		from = params.DateFrom
		to = params.DateTo
	}

	searchLower := ""
	if params.Search != "" {
		searchLower = strings.ToLower(params.Search)
	}

	out := make([]models.FlatRecord, 0, len(records))
	for _, rec := range records {
		// Фильтр по отделению
		if params.Department != "" && rec.Department != params.Department {
			continue
		}

		// Фильтр по группе
		if params.Group != "" && rec.Group != params.Group {
			continue
		}

		// Фильтр по студенту
		if params.Student != "" && rec.Student != params.Student {
			continue
		}

		// Поиск по подстроке
		if searchLower != "" {
			ok := strings.Contains(strings.ToLower(rec.Department), searchLower) ||
				strings.Contains(strings.ToLower(rec.Group), searchLower) ||
				strings.Contains(strings.ToLower(rec.Student), searchLower)
			if !ok {
				continue
			}
		}

		// Фильтр по минимальному количеству пропусков
		if params.MissedMin >= 0 && rec.Missed < params.MissedMin {
			continue
		}

		// Фильтр по дате
		if from != "" || to != "" {
			if from != "" && rec.Date < from {
				continue
			}
			if to != "" && rec.Date > to {
				continue
			}
		} else {
			if params.Date == "today" {
				if rec.Date != today {
					continue
				}
			} else if params.Date != "" && rec.Date != params.Date {
				continue
			}
		}

		out = append(out, rec)
	}

	return out
}

// ParseFilterParams извлекает параметры фильтрации из HTTP запроса
func ParseFilterParams(r *http.Request) FilterParams {
	q := r.URL.Query()
	missedMin := -1
	if s := q.Get("missed_min"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n >= 0 {
			missedMin = n
		}
	}

	return FilterParams{
		Department: q.Get("department"),
		Group:      q.Get("group"),
		Student:    q.Get("student"),
		Date:       q.Get("date"),
		DateFrom:   q.Get("date_from"),
		DateTo:     q.Get("date_to"),
		Period:     q.Get("period"),
		Search:     strings.TrimSpace(q.Get("q")),
		MissedMin:  missedMin,
	}
}

func todayAdd(days int) string {
	t := time.Now().AddDate(0, 0, days)
	return t.Format("2006-01-02")
}

// SummaryResponse ответ для сводки
type SummaryResponse struct {
	TotalStudents int             `json:"total_students"`
	Present       int             `json:"present"`
	Absent        int             `json:"absent"`
	ByDepartment  []DeptDrillItem `json:"by_department,omitempty"`
}

// DeptDrillItem элемент drill-down по отделению
type DeptDrillItem struct {
	Department  string `json:"department"`
	Total       int    `json:"total"`
	Absent      int    `json:"absent"`
	MissedTotal int    `json:"missed_total"`
}

// GroupDrillItem элемент drill-down по группе
type GroupDrillItem struct {
	Group       string `json:"group"`
	Total       int    `json:"total"`
	Absent      int    `json:"absent"`
	MissedTotal int    `json:"missed_total"`
}

// StudentDrillItem элемент drill-down по студенту
type StudentDrillItem struct {
	Student     string   `json:"student"`
	MissedTotal int      `json:"missed_total"`
	Records     int      `json:"records"`
	Dates       []string `json:"dates,omitempty"`
}

// BuildSummary строит сводку по данным
func (s *AttendanceService) BuildSummary(departments []models.DepartmentJSON, filtered []models.FlatRecord) SummaryResponse {
	byDept := totalByDept(departments)
	absentSet := make(map[string]struct{})
	deptAbsent := make(map[string]int)
	deptMissed := make(map[string]int)
	deptsInScope := make(map[string]struct{})

	for _, rec := range filtered {
		deptsInScope[rec.Department] = struct{}{}
		k := rec.Department + "\x00" + rec.Group + "\x00" + rec.Student
		if _, ok := absentSet[k]; !ok {
			absentSet[k] = struct{}{}
			deptAbsent[rec.Department]++
		}
		deptMissed[rec.Department] += rec.Missed
	}

	var total int
	if len(deptsInScope) > 0 {
		for d := range deptsInScope {
			total += byDept[d]
		}
	} else {
		for _, n := range byDept {
			total += n
		}
	}

	absent := len(absentSet)
	present := total - absent
	if present < 0 {
		present = 0
	}

	var byDepartment []DeptDrillItem
	iter := byDept
	if len(deptsInScope) > 0 {
		iter = make(map[string]int)
		for d := range deptsInScope {
			iter[d] = byDept[d]
		}
	}

	for dept, tot := range iter {
		byDepartment = append(byDepartment, DeptDrillItem{
			Department:  dept,
			Total:       tot,
			Absent:      deptAbsent[dept],
			MissedTotal: deptMissed[dept],
		})
	}

	return SummaryResponse{
		TotalStudents: total,
		Present:       present,
		Absent:        absent,
		ByDepartment:  byDepartment,
	}
}

// BuildDrillDepartments строит drill-down по отделениям
func (s *AttendanceService) BuildDrillDepartments(departments []models.DepartmentJSON, filtered []models.FlatRecord) []DeptDrillItem {
	byDept := totalByDept(departments)
	deptAbsent := make(map[string]int)
	deptMissed := make(map[string]int)
	seen := make(map[string]map[string]struct{})
	deptsInScope := make(map[string]struct{})

	for _, rec := range filtered {
		deptsInScope[rec.Department] = struct{}{}
		if seen[rec.Department] == nil {
			seen[rec.Department] = make(map[string]struct{})
		}
		k := rec.Group + "\x00" + rec.Student
		if _, ok := seen[rec.Department][k]; !ok {
			seen[rec.Department][k] = struct{}{}
			deptAbsent[rec.Department]++
		}
		deptMissed[rec.Department] += rec.Missed
	}

	iter := byDept
	if len(deptsInScope) > 0 {
		iter = make(map[string]int)
		for d := range deptsInScope {
			iter[d] = byDept[d]
		}
	}

	var out []DeptDrillItem
	for dept, tot := range iter {
		out = append(out, DeptDrillItem{
			Department:  dept,
			Total:       tot,
			Absent:      deptAbsent[dept],
			MissedTotal: deptMissed[dept],
		})
	}
	return out
}

// BuildDrillGroups строит drill-down по группам
func (s *AttendanceService) BuildDrillGroups(departments []models.DepartmentJSON, filtered []models.FlatRecord, department string) []GroupDrillItem {
	byGroup := totalByGroup(departments)
	if byGroup[department] == nil {
		return []GroupDrillItem{}
	}

	grpAbsent := make(map[string]int)
	grpMissed := make(map[string]int)
	seen := make(map[string]map[string]struct{})

	for _, rec := range filtered {
		if rec.Department != department {
			continue
		}
		if seen[rec.Group] == nil {
			seen[rec.Group] = make(map[string]struct{})
		}
		if _, ok := seen[rec.Group][rec.Student]; !ok {
			seen[rec.Group][rec.Student] = struct{}{}
			grpAbsent[rec.Group]++
		}
		grpMissed[rec.Group] += rec.Missed
	}

	var out []GroupDrillItem
	for grp, tot := range byGroup[department] {
		out = append(out, GroupDrillItem{
			Group:       grp,
			Total:       tot,
			Absent:      grpAbsent[grp],
			MissedTotal: grpMissed[grp],
		})
	}
	return out
}

// BuildDrillStudents строит drill-down по студентам
func (s *AttendanceService) BuildDrillStudents(filtered []models.FlatRecord, department, group string) []StudentDrillItem {
	type agg struct {
		missed int
		dates  []string
	}
	m := make(map[string]*agg)

	for _, rec := range filtered {
		if rec.Department != department || rec.Group != group {
			continue
		}
		if m[rec.Student] == nil {
			m[rec.Student] = &agg{dates: []string{}}
		}
		m[rec.Student].missed += rec.Missed
		m[rec.Student].dates = append(m[rec.Student].dates, rec.Date)
	}

	out := make([]StudentDrillItem, 0, len(m))
	for name, a := range m {
		out = append(out, StudentDrillItem{
			Student:     name,
			MissedTotal: a.missed,
			Records:     len(a.dates),
			Dates:       a.dates,
		})
	}
	return out
}

// Вспомогательные функции
func totalByDept(departments []models.DepartmentJSON) map[string]int {
	m := make(map[string]int)
	for _, d := range departments {
		n := 0
		for _, g := range d.Groups {
			n += len(g.Students)
		}
		if n > 0 {
			m[d.Department] = n
		}
	}
	return m
}

func totalByGroup(departments []models.DepartmentJSON) map[string]map[string]int {
	m := make(map[string]map[string]int)
	for _, d := range departments {
		if m[d.Department] == nil {
			m[d.Department] = make(map[string]int)
		}
		for _, g := range d.Groups {
			m[d.Department][g.Group] = len(g.Students)
		}
	}
	return m
}
