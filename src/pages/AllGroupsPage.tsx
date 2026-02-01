import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PieDiagram } from "@/components/ui/pieDiagram";
import { FilterDropdown } from "@/widgets/filterDropdown";
import { useCallback, useEffect, useMemo, useState } from "react";

import loadAttendance from "@/lib/xlsxAttendanceConverter";
import type {
  AttendanceRecord,
  Department,
  FlatStudent,
} from "@/types/attendance";
import { useNavigate } from "react-router-dom";

function searchGroup(group: string, search: string) {
  if (!search) return true;

  const g = group.toLowerCase();
  const s = search.toLowerCase();

  if (/\d$/.test(s)) return g.startsWith(s);
  return g.replace(/\d$/, "").includes(s);
}

const LIMIT = 10;

const PIE_DATA = [
  { name: "По уважительной причине", color: "#54EB66" },
  { name: "По неуважительной причине", color: "#EA5596" },
  { name: "Присутствуют", color: "#4FB4E5" },
] as const;

function AllGroupsPage() {
  const [data, setData] = useState<Department[]>([]);
  const [searchValue, setSearchValue] = useState("");
  const [selectedGroup, setSelectedGroup] = useState("");
  const [selectedDepartment, setSelectedDepartment] = useState("");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const goToGroup = useCallback(
    (group: string) => {
      navigate("/group/" + group);
    },
    [navigate]
  );

  const [limit, setLimit] = useState(LIMIT);

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
      .sort((a, b) => b.missed - a.missed);

    const deptData = Object.entries(missedByDepartment).map(
      ([department, missed]) => ({
        department,
        missed,
      })
    );

    const groupToDepartment: Record<string, string> = {};

    data.forEach((dept) => {
      dept.groups?.forEach((group) => {
        groupToDepartment[group.group] = dept.department;
      });
    });

    return {
      totalStudents: allStudents.length,
      totalMissed,
      groupToDepartment,
      topOffenders,
      dateData,
      groupData,
      deptData,
      allStudents,
    };
  }, [data]);

  const filteredGroups = useMemo(() => {
    if (!stats) return [];

    return stats.groupData
      .map(({ group }) => group)
      .filter((group) => searchGroup(group, searchValue));
  }, [stats, searchValue]);

  const filteredGroupData = useMemo(() => {
    if (!stats) return [];

    return stats.groupData.filter(({ group }) => {
      if (selectedGroup && selectedGroup !== group) return false;

      if (
        selectedDepartment &&
        stats.groupToDepartment[group] !== selectedDepartment
      ) {
        return false;
      }

      return true;
    });
  }, [stats, selectedGroup, selectedDepartment]);

  const paginatedGroups = useMemo(() => {
    return filteredGroupData.slice(0, limit);
  }, [filteredGroupData, limit]);

  const groupCards = useMemo(() => {
    return paginatedGroups.map(({ group }) => (
      <Card key={group} onClick={() => goToGroup(group)}>
        <CardHeader className="h-0 p-0">
          <CardTitle className="relative m-6">{group}</CardTitle>
        </CardHeader>
        <CardContent>
          <PieDiagram
            data={[
              { ...PIE_DATA[0], value: 78 },
              { ...PIE_DATA[1], value: 18 },
              { ...PIE_DATA[2], value: 245 },
            ]}
            valueLabel={{ one: "час", few: "часа", many: "часов" }}
          />
        </CardContent>
      </Card>
    ));
  }, [paginatedGroups, goToGroup]);

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
      <div className="space-y-4">
        <div className="space-y-2">
          <h2 className="text-lg font-semibold text-black">Общая статистика</h2>
          <div className="space-x-6">
            <FilterDropdown
              label="Группы"
              list={filteredGroups}
              onSelect={(group) => {
                setSelectedGroup(group);
                setLimit(LIMIT);
                setSelectedDepartment("");
              }}
              buttonSize={{ height: 35, width: 133 }}
              allSelection={!searchValue ? "Все группы" : undefined}
              onChangeSearch={setSearchValue}
              searchInput
            />

            <FilterDropdown
              label="Отделения"
              buttonSize={{ height: 35, width: 133 }}
              popupSize={{ height: 189, width: 273 }}
              list={stats.deptData.map(({ department }) =>
                department.length > 30
                  ? department.substring(0, 28) + "..."
                  : department
              )}
              onSelect={(dep) => {
                setSelectedDepartment(dep);
                setLimit(LIMIT);
                setSelectedGroup("");
              }}
              allSelection="Все отделения"
            />
          </div>
        </div>
        <div className="grid gap-6 grid-cols-1 2xl:grid-cols-2">
          {groupCards}
          {limit < filteredGroupData.length && (
            <div className="flex justify-center">
              <button
                onClick={() => setLimit((prev) => prev + LIMIT)}
                className="mt-4 rounded-md border px-6 py-2 text-sm font-medium hover:bg-gray-100 transition"
              >
                Показать ещё
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
export default AllGroupsPage;
