import type {
  HTMLAttributes,
  TableHTMLAttributes,
  TdHTMLAttributes,
  ThHTMLAttributes,
} from "react";

export type TableSectionProps = HTMLAttributes<HTMLTableSectionElement>;
export type TableCaptionProps = HTMLAttributes<HTMLTableCaptionElement>;
export type TableWrapperProps = TableHTMLAttributes<HTMLTableElement>;
export type TableRowProps = HTMLAttributes<HTMLTableRowElement>;
export type TableCellHeadProps = ThHTMLAttributes<HTMLTableCellElement>;
export type TableCellProps = TdHTMLAttributes<HTMLTableCellElement>;
