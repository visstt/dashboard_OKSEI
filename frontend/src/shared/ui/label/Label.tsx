import { cn } from "@/shared/lib/cn";
import type { LabelProps } from "./types";

export const Label = ({ className, ...props }: LabelProps) => (
  <label
    className={cn(
      "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 text-black",
      className
    )}
    {...props}
  />
);
