import { cn } from "@/shared/lib/cn";
import type { TableCellHeadProps } from "../types";

export const TableHead = ({ className, ...props }: TableCellHeadProps) => (
  <th
    className={cn(
      "h-12 px-4 text-left align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0",
      className
    )}
    {...props}
  />
);
