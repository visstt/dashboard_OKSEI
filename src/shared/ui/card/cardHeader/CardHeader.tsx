import { cn } from "@/shared/lib/cn";
import type { CardProps } from "../types";

export const CardHeader = ({ className, ...props }: CardProps) => (
  <div className={cn("flex flex-col space-y-1.5 p-6", className)} {...props} />
);
