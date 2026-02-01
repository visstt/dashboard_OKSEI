import { cn } from "@/shared/lib/cn";
import type { TableSectionProps } from "../types";

export const TableFooter = ({ className, ...props }: TableSectionProps) => (
  <tfoot
    className={cn(
      "border-t bg-muted/50 font-medium [&>tr]:last:border-b-0",
      className
    )}
    {...props}
  />
);
