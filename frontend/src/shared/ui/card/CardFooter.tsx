import { cn } from "@/shared/lib/cn";
import type { CardProps } from "./types";

export const CardFooter = ({ className, ...props }: CardProps) => (
  <div className={cn("flex items-center p-6 pt-0", className)} {...props} />
);
