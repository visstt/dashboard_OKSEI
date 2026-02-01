import type { HTMLAttributes } from "react";

export interface GroupsFiltersProps extends HTMLAttributes<HTMLDivElement> {
  onGroupSelect?: (group: string) => void;
  onDepartmentSelect?: (department: string) => void;
  onChangeGroupSearch?: (value: string) => void;
  groupsList: string[];
  departmentsList: string[];
  groupSearchValue?: string;
}
