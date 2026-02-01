import { cn } from "@/shared/lib/cn";
import type { CardProps } from "./types";

export const CardContent = ({ className, ...props }: CardProps) => (
  <div className={cn("p-6 pt-0", className)} {...props} />
);
