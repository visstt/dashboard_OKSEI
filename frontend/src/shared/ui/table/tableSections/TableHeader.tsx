import { cn } from "@/shared/lib/cn";
import type { TableSectionProps } from "../types";

export const TableHeader = ({ className, ...props }: TableSectionProps) => (
  <thead className={cn("[&_tr]:border-b", className)} {...props} />
);
