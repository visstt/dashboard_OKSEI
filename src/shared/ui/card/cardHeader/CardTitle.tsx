import { cn } from "@/shared/lib/cn";
import type { CardProps } from "../types";

export const CardTitle = ({ className, ...props }: CardProps) => (
  <h3
    className={cn(
      "text-2xl font-semibold leading-none tracking-tight",
      className
    )}
    {...props}
  />
);
