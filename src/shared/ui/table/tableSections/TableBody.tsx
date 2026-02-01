import { cn } from "@/shared/lib/cn";
import type { TableSectionProps } from "../types";

export const TableBody = ({ className, ...props }: TableSectionProps) => (
  <tbody className={cn("[&_tr:last-child]:border-0", className)} {...props} />
);
