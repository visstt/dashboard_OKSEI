import { cn } from "@/shared/lib/cn";
import type { PieDiagramProps } from "../types";
import { pluralizeRu } from "@/shared/lib/pluralizeRu";

export const PieLegend = ({ data, valueLabel, className }: PieDiagramProps) => {
  return (
    <div className={cn("flex flex-col justify-center space-y-3", className)}>
      {data.map((item, index) => (
        <div key={index} className="flex items-center gap-3">
          <div
            className="w-4 h-4 rounded"
            style={{
              backgroundColor: item.color,
            }}
          />
          <div className="flex-1">
            <div className="text-sm font-medium text-black">{item.name}</div>
            <div className="text-xs text-gray-600">
              {item.value}{" "}
              {valueLabel &&
                (typeof valueLabel === "string"
                  ? valueLabel
                  : pluralizeRu(item.value, valueLabel))}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
