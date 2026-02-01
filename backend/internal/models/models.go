package models

import "time"

// Department модель отделения
type Department struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Group модель группы
type Group struct {
	ID           int       `json:"id" db:"id"`
	DepartmentID int       `json:"department_id" db:"department_id"`
	Name         string    `json:"name" db:"name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Student модель студента
type Student struct {
	ID        int       `json:"id" db:"id"`
	GroupID   int       `json:"group_id" db:"group_id"`
	FullName  string    `json:"full_name" db:"full_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Attendance модель посещаемости
type Attendance struct {
	ID          int       `json:"id" db:"id"`
	StudentID   int       `json:"student_id" db:"student_id"`
	Date        time.Time `json:"date" db:"date"`
	MissedHours int       `json:"missed_hours" db:"missed_hours"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Specialty модель специальности (для summary)
type Specialty struct {
	ID           int       `json:"id" db:"id"`
	DepartmentID int       `json:"department_id" db:"department_id"`
	Name         string    `json:"name" db:"name"`
	TotalMissed  int       `json:"total_missed" db:"total_missed"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SummaryGroup модель группы в ведомости
type SummaryGroup struct {
	ID          int       `json:"id" db:"id"`
	SpecialtyID int       `json:"specialty_id" db:"specialty_id"`
	Name        string    `json:"name" db:"name"`
	TotalMissed int       `json:"total_missed" db:"total_missed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SummaryStudent модель студента в ведомости
type SummaryStudent struct {
	ID            int       `json:"id" db:"id"`
	SummaryGroupID int      `json:"summary_group_id" db:"summary_group_id"`
	FullName      string    `json:"full_name" db:"full_name"`
	MissedTotal   int       `json:"missed_total" db:"missed_total"`
	MissedBad     int       `json:"missed_bad" db:"missed_bad"`
	MissedExcused int       `json:"missed_excused" db:"missed_excused"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// JSON модели для работы с JSON файлами (attendance.json)

// AttendanceRecordJSON запись посещаемости в JSON
type AttendanceRecordJSON struct {
	Date   string `json:"date"`
	Missed int    `json:"missed"`
}

// StudentJSON студент в JSON
type StudentJSON struct {
	Student    string                `json:"student"`
	Attendance []AttendanceRecordJSON `json:"attendance"`
}

// GroupJSON группа в JSON
type GroupJSON struct {
	Group    string        `json:"group"`
	Students []StudentJSON `json:"students"`
}

// DepartmentJSON отделение в JSON
type DepartmentJSON struct {
	Department string      `json:"department"`
	Groups     []GroupJSON `json:"groups"`
}
