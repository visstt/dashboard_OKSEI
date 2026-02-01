import { cn } from "@/shared/lib/cn";
import type { SeparatorProps } from "./types";
import { separatorOrientations } from "./separatorOrientations";

export const Separator = ({
  className,
  orientation = "horizontal",
  decorative = true,
  ...props
}: SeparatorProps) => (
  <div
    role={decorative ? "none" : "separator"}
    aria-orientation={orientation}
    className={cn(
      "shrink-0 bg-slate-200",
      separatorOrientations[orientation],
      className
    )}
    {...props}
  />
);
