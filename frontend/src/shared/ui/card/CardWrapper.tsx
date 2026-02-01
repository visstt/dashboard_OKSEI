import { cn } from "@/shared/lib/cn";
import type { CardProps } from "./types";

export const CardWrapper = ({ className, ...props }: CardProps) => (
  <div
    className={cn(
      "rounded-lg border bg-card text-card-foreground shadow-sm",
      className
    )}
    {...props}
  />
);
