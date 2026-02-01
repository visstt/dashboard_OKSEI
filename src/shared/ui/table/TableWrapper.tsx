import { cn } from "@/shared/lib/cn";
import type { TableWrapperProps } from "./types";

export const TableWrapper = ({ className, ...props }: TableWrapperProps) => (
  <div className="relative w-full overflow-auto">
    <table
      className={cn("w-full caption-bottom text-sm", className)}
      {...props}
    />
  </div>
);
