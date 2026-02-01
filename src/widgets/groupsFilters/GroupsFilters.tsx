import { FilterDropdown } from "@/entities/filterDropdown";
import type { GroupsFiltersProps } from "./types";
import { cn } from "@/shared/lib/cn";

export default function GroupsFilters({
  onGroupSelect,
  groupsList,
  groupSearchValue,
  onChangeGroupSearch,
  onDepartmentSelect,
  departmentsList,
  className,
}: GroupsFiltersProps) {
  return (
    <div className={cn("space-x-6", className)}>
      <FilterDropdown
        label="Группы"
        list={groupsList}
        onSelect={onGroupSelect}
        buttonSize={{ height: 35, width: 133 }}
        popupSize={{ height: 189, width: 150 }}
        allSelection={!groupSearchValue ? "Все группы" : undefined}
        onChangeSearch={onChangeGroupSearch}
        searchInput
      />

      <FilterDropdown
        label="Отделения"
        buttonSize={{ height: 35, width: 133 }}
        popupSize={{ height: 189, width: 273 }}
        list={departmentsList}
        onSelect={onDepartmentSelect}
        allSelection="Все отделения"
      />
    </div>
  );
}
