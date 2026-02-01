import { cn } from "@/shared/lib/cn";
import type { CardProps } from "../types";

export const CardDescription = ({ className, ...props }: CardProps) => (
  <p className={cn("text-sm text-muted-foreground", className)} {...props} />
);
