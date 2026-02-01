import { cn } from "@/shared/lib/cn";
import type { TableCellProps } from "../types";

export const TableCell = ({ className, ...props }: TableCellProps) => (
  <td
    className={cn("p-4 align-middle [&:has([role=checkbox])]:pr-0", className)}
    {...props}
  />
);
