import { cn } from "@/shared/lib/cn";
import type { TableRowProps } from "./types";

export const TableRow = ({ className, ...props }: TableRowProps) => (
  <tr
    className={cn(
      "border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted",
      className
    )}
    {...props}
  />
);