package models

// FlatRecord представляет плоскую запись посещаемости
type FlatRecord struct {
	Department string `json:"department"`
	Group      string `json:"group"`
	Student    string `json:"student"`
	Date       string `json:"date"`
	Missed     int    `json:"missed"`
}

// Flatten преобразует иерархию DepartmentJSON → GroupJSON → StudentJSON → AttendanceRecordJSON
// в плоский список записей для удобной фильтрации и поиска
func Flatten(departments []DepartmentJSON) []FlatRecord {
	var out []FlatRecord
	for _, d := range departments {
		for _, g := range d.Groups {
			for _, s := range g.Students {
				for _, a := range s.Attendance {
					out = append(out, FlatRecord{
						Department: d.Department,
						Group:      g.Group,
						Student:    s.Student,
						Date:       a.Date,
						Missed:     a.Missed,
					})
				}
			}
		}
	}
	return out
}
