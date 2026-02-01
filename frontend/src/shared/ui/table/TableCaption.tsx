import { cn } from "@/shared/lib/cn";
import type { TableCaptionProps } from "./types";

export const TableCaption = ({ className, ...props }: TableCaptionProps) => (
  <caption
    className={cn("mt-4 text-sm text-muted-foreground", className)}
    {...props}
  />
);
