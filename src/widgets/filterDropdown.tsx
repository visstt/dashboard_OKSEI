import { useEffect, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { FiltersPopup } from "@/components/ui/filtersPopup";
import { LucideArrowDownWideNarrow } from "lucide-react";

interface SizeData {
  width: number | string;
  height: number | string;
}

interface FilterDropdownProps {
  label: string;
  list: string[];
  buttonSize?: SizeData;
  popupSize?: SizeData;
  allSelection?: string;
  searchInput?: boolean;
  onChangeSearch?: (value: string) => void;
  onSelect?: (value: string) => void;
}

export const FilterDropdown = ({
  label,
  list,
  searchInput,
  popupSize,
  buttonSize,
  allSelection,
  onChangeSearch,
  onSelect,
}: FilterDropdownProps) => {
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (!rootRef.current?.contains(e.target as Node)) {
        setOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  return (
    <div ref={rootRef} className="relative inline-block">
      <Button style={buttonSize} onClick={() => setOpen((v) => !v)}>
        {label}
        <LucideArrowDownWideNarrow
          style={open ? { transform: "scaleY(-1)" } : {}}
          height={18}
          width={18}
          className="ml-1.5 transition-all duration-200 ease-out"
        />
      </Button>

      <FiltersPopup
        list={list}
        isOpen={open}
        maxHeight={popupSize?.height || 189}
        width={popupSize?.width || 150}
        onSelect={(value) => {
          onSelect?.(value);
          setOpen(false);
        }}
        allSelection={allSelection}
        onChangeSearch={(value) => onChangeSearch?.(value)}
        searchInput={searchInput}
      />
    </div>
  );
};
