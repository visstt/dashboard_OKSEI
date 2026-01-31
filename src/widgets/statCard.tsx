import type { HTMLAttributes } from "react";
import type { LucideIcon } from "lucide-react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { pluralizeRu, type PluralizeRuFormats } from "@/lib/pluralizeRu";

interface StatCardProps extends HTMLAttributes<HTMLDivElement> {
  title: string;
  value: number;
  icon: LucideIcon;
  label?: string | PluralizeRuFormats;
}

export const StatCard = ({
  title,
  value,
  label,
  icon: Icon,
  className,
}: StatCardProps) => {
  return (
    <Card className={className}>
      <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
        <CardTitle className="text-sm font-medium text-black">
          {title}
        </CardTitle>
        <Icon className="h-4 w-4 text-black" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold text-black">{value}</div>
        <p className="text-xs text-gray-600 mt-1">
          {label &&
            (typeof label === "string" ? label : pluralizeRu(value, label))}
        </p>
      </CardContent>
    </Card>
  );
};
