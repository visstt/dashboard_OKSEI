export interface AttendanceRecord {
  date: string;
  missed: number;
}

export interface Student {
  student: string;
  attendance?: AttendanceRecord[];
}

export interface Group {
  group: string;
  students?: Student[];
}

export interface Department {
  department: string;
  groups?: Group[];
}

export interface FlatStudent {
  name: string;
  group: string;
  department: string;
  totalMissed: number;
  attendance: AttendanceRecord[];
}
