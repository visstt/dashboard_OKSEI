import { useState, useEffect, useMemo } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import loadAttendance from "@/lib/xlsxAttendanceConverter";
import type { AttendanceRecord, Department } from "@/types/attendance";
import { PieDiagram } from "@/components/ui/pieDiagram";

export function GroupPage() {
  const { groupName } = useParams();
  const navigate = useNavigate();
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

  const groupData = useMemo(() => {
    if (!data.length || !groupName) return null;

    let groupInfo = null;
    let departmentName = "";

    // Находим группу
    for (const dept of data) {
      const foundGroup = dept.groups?.find((g) => g.group === groupName);
      if (foundGroup) {
        groupInfo = foundGroup;
        departmentName = dept.department;
        break;
      }
    }

    if (!groupInfo) return null;

    // Собираем студентов с пропусками
    const students = groupInfo.students
      ?.map((student) => {
        const totalMissed =
          student.attendance?.reduce(
            (sum: number, att: AttendanceRecord) => sum + att.missed,
            0
          ) || 0;
        return {
          name: student.student,
          totalMissed,
          attendance: student.attendance || [],
        };
      })
      .filter((s) => s.totalMissed > 0)
      .sort((a, b) => b.totalMissed - a.totalMissed);

    const totalMissed =
      students?.reduce((sum, s) => sum + s.totalMissed, 0) || 0;

    return {
      groupName,
      departmentName,
      students,
      totalMissed,
      totalStudents: students?.length || 0,
    };
  }, [data, groupName]);

  if (loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <p className="text-lg text-black">Загрузка данных...</p>
      </div>
    );
  }

  if (!groupData) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-center space-y-4">
          <p className="text-lg text-black">Группа не найдена</p>
          <Button
            onClick={() => navigate("/groups")}
            variant="default"
            size="default"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Вернуться назад
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Шапка */}
      <div className="space-y-4">
        <Button
          variant="outline"
          size="default"
          onClick={() => navigate("/groups")}
          className="mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Назад к списку групп
        </Button>

        <div>
          <h1 className="text-3xl font-semibold tracking-tight text-black">
            Группа {groupData.groupName}
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            {groupData.departmentName}
          </p>
        </div>
      </div>

      {/* Статистика группы */}
      <Card className="lg:col-span-2">
        <CardContent>
          <PieDiagram
            data={[
              {
                name: "По уважительной причине",
                color: "#54EB66",
                value: 132,
              },
              {
                name: "По неуважительной причине",
                color: "#EA5596",
                value: 245,
              },
              {
                name: "Присутствуют",
                color: "#4FB4E5",
                value: 825,
              },
            ]}
            valueLabel={{
              one: "час",
              few: "часа",
              many: "часов",
            }}
          />
        </CardContent>
      </Card>

      {/* Таблица студентов */}
      <Card>
        <CardHeader>
          <CardTitle>Список студентов с пропусками</CardTitle>
          <CardDescription>
            Обучающиеся группы {groupData.groupName} с пропусками по
            неуважительной причине
          </CardDescription>
        </CardHeader>
        <CardContent>
          {groupData.students && groupData.students.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="text-black">#</TableHead>
                  <TableHead className="text-black">ФИО студента</TableHead>
                  <TableHead className="text-right text-black">
                    Пропущено часов
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {groupData.students.map((student, index) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium text-black">
                      {index + 1}
                    </TableCell>
                    <TableCell className="text-black">{student.name}</TableCell>
                    <TableCell className="text-right">
                      <Badge
                        variant="outline"
                        style={{ height: 21, cursor: "default" }}
                        className={
                          student.totalMissed > 20
                            ? "bg-red-600 text-white border-none"
                            : student.totalMissed > 10
                              ? "bg-orange-600 text-white border-none"
                              : "bg-green-500 text-white border-none"
                        }
                      >
                        {student.totalMissed} ч
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <p className="text-center text-muted-foreground py-8">
              В данной группе нет студентов с пропусками
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
