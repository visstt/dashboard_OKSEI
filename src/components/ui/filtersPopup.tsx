import { cn } from "@/lib/utils";
import { useState } from "react";

interface FiltersPopupProps {
  list: string[];
  isOpen: boolean;
  width: number | string;
  allSelection?: string;
  maxHeight: number | string;
  onChangeSearch?: (value: string) => void;
  onSelect?: (value: string) => void;
  searchInput?: boolean;
}

export const FiltersPopup = ({
  list,
  isOpen,
  width,
  maxHeight,
  searchInput,
  allSelection,
  onChangeSearch,
  onSelect,
}: FiltersPopupProps) => {
  const [overflow, setOverflow] = useState<"scroll" | "hidden">("hidden");

  return (
    <div
      style={{
        width,
        maxHeight: isOpen ? maxHeight : 0,
        overflowY: overflow,
      }}
      className={cn(
        "ml-2 absolute left-0 top-full mt-1.5 z-50 overflow-x-hidden rounded-md border bg-white shadow-lg",
        "transition-[max-height,opacity,transform] duration-200 ease-out",
        isOpen
          ? "opacity-100 translate-y-0"
          : "opacity-0 -translate-y-1 pointer-events-none"
      )}
      onTransitionRun={() => setOverflow("hidden")}
      onTransitionEnd={() => isOpen && setOverflow("scroll")}
    >
      {searchInput && (
        <input
          className="w-full border-b px-3 py-2 text-sm focus:outline-none"
          placeholder="Поиск..."
          onChange={(e) => onChangeSearch?.(e.target.value)}
        />
      )}

      <ul className="py-1 text-sm text-black">
        {allSelection && (
          <li
            onClick={() => onSelect?.("")}
            className="cursor-pointer px-4 py-2 hover:bg-gray-100"
          >
            {allSelection}
          </li>
        )}
        {list.length ? (
          list.map((item) => (
            <li
              key={item}
              onClick={() => onSelect?.(item)}
              className="cursor-pointer px-4 py-2 hover:bg-gray-100"
            >
              {item}
            </li>
          ))
        ) : (
          <div className="px-4 py-2 text-sm text-gray-400">Нет данных</div>
        )}
      </ul>
    </div>
  );
};
