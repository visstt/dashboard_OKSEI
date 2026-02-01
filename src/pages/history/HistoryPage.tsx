import { useState, useEffect, useMemo } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import {
  CardWrapper,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/card";

import loadAttendance from "@/features/attendanceConverter/xlsxAttendanceConverter";
import type {
  AttendanceRecord,
  Department,
  FlatStudent,
} from "@/features/attendanceConverter/attendance";

export function HistoryPage() {
  const [data, setData] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadAttendance()
      .then((jsonData) => {
        setData(jsonData);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error loading data:", err);
        setLoading(false);
      });
  }, []);

  const stats = useMemo(() => {
    if (!data.length) return null;

    const allStudents: FlatStudent[] = [];
    const missedByDate: Record<string, number> = {};
    const missedByGroup: Record<string, number> = {};
    const missedByDepartment: Record<string, number> = {};

    data.forEach((dept) => {
      missedByDepartment[dept.department] = 0;
      dept.groups?.forEach((group) => {
        missedByGroup[group.group] = 0;
        group.students?.forEach((student) => {
          const totalMissed =
            student.attendance?.reduce(
              (sum: number, att: AttendanceRecord) => sum + att.missed,
              0
            ) || 0;
          allStudents.push({
            name: student.student,
            group: group.group,
            department: dept.department,
            totalMissed,
            attendance: student.attendance || [],
          });
          missedByGroup[group.group] += totalMissed;
          missedByDepartment[dept.department] += totalMissed;

          student.attendance?.forEach((att) => {
            missedByDate[att.date] = (missedByDate[att.date] || 0) + att.missed;
          });
        });
      });
    });

    const totalMissed = allStudents.reduce((sum, s) => sum + s.totalMissed, 0);

    const topOffenders = [...allStudents]
      .sort((a, b) => b.totalMissed - a.totalMissed)
      .slice(0, 10);

    const dateData = Object.entries(missedByDate)
      .sort(([a], [b]) => new Date(a).getTime() - new Date(b).getTime())
      .map(([date, missed]) => ({
        date: new Date(date).toLocaleDateString("ru-RU", {
          day: "2-digit",
          month: "2-digit",
        }),
        missed,
      }));

    const groupData = Object.entries(missedByGroup)
      .map(([group, missed]) => ({ group, missed }))
      .sort((a, b) => b.missed - a.missed)
      .slice(0, 10);

    const deptData = Object.entries(missedByDepartment).map(
      ([department, missed]) => ({
        department:
          department.length > 30
            ? department.substring(0, 30) + "..."
            : department,
        missed,
      })
    );

    return {
      totalStudents: allStudents.length,
      totalMissed,
      topOffenders,
      dateData,
      groupData,
      deptData,
      allStudents,
    };
  }, [data]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-lg">Загрузка данных...</div>
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-lg">Нет данных для отображения</div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Графики */}
      <div className="space-y-4">
        <div className="grid grid-cols-1 gap-6">
          <CardWrapper>
            <CardHeader>
              <CardTitle>История пропусков</CardTitle>
              <CardDescription>Динамика пропусков по датам</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={stats.dateData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                  <XAxis dataKey="date" stroke="#64748b" fontSize={12} />
                  <YAxis stroke="#64748b" fontSize={12} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "white",
                      border: "1px solid #e5e7eb",
                      borderRadius: "8px",
                    }}
                  />
                  <Line
                    type="monotone"
                    dataKey="missed"
                    stroke="#D26A69"
                    strokeWidth={2}
                    dot={{ fill: "#D26A69", r: 4 }}
                    activeDot={{ r: 6 }}
                    name="Пропущено часов"
                  />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </CardWrapper>
        </div>
      </div>
    </div>
  );
}
