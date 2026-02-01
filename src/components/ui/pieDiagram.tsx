import { useEffect, useState, type HTMLAttributes } from "react";
import { pluralizeRu, type PluralizeRuFormats } from "@/lib/pluralizeRu";
import { Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";

interface PieDiagramData {
  name: string;
  color: string;
  value: number;
}

interface PieDiagramProps extends HTMLAttributes<HTMLDivElement> {
  data: PieDiagramData[];
  valueLabel?: string | PluralizeRuFormats;
}

export const PieDiagram = ({
  data,
  valueLabel,
  className,
}: PieDiagramProps) => {
  const [animation, setAnimation] = useState(true);

  useEffect(() => {
    const t = setTimeout(() => setAnimation(false), 2200);

    return () => clearTimeout(t);
  }, []);

  return (
    <div className={`grid grid-cols-1 md:grid-cols-2 gap-6 ${className}`}>
      <ResponsiveContainer width="100%" height={300}>
        <PieChart>
          <Pie
            data={data.map((item) => ({
              ...item,
              fill: item.color,
            }))}
            cx="50%"
            cy="50%"
            labelLine={false}
            label={({ percent = 0 }) =>
              percent === 0
                ? "0%"
                : Math.round(percent * 100) === 0
                  ? "<1%"
                  : `${Math.round(percent * 100)}%`
            }
            outerRadius={100}
            animationDuration={1200}
            isAnimationActive={animation}
          />
          <Tooltip />
        </PieChart>
      </ResponsiveContainer>
      <div className="flex flex-col justify-center space-y-3">
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
    </div>
  );
};
