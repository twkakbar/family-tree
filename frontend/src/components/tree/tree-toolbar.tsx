"use client";

import type { ReactNode } from "react";

import { Button } from "@/components/ui/button";
import { cn } from "@/utils/cn";

interface TreeToolbarProps {
  canEdit: boolean;
  scale: number;
  onZoomIn: () => void;
  onZoomOut: () => void;
  onResetView: () => void;
  className?: string;
}

/** Compact controls and hints for the tree canvas. */
export function TreeToolbar({
  canEdit,
  scale,
  onZoomIn,
  onZoomOut,
  onResetView,
  className,
}: TreeToolbarProps) {
  return (
    <div
      className={cn(
        "flex flex-col gap-3 border-b border-slate-200 px-3 py-3 sm:flex-row sm:items-center sm:justify-between sm:px-4",
        className,
      )}
    >
      <div className="min-w-0">
        <h2 className="text-base font-semibold text-slate-900">Family tree</h2>
        <p className="mt-0.5 hidden text-sm text-slate-500 sm:block">
          {canEdit
            ? "Pan the background, scroll to zoom, drag one member onto another to connect them."
            : "Pan the background, scroll to zoom, tap a member for details."}
        </p>
        <div className="mt-2 flex flex-wrap gap-1.5 sm:hidden">
          <HintChip>Pan</HintChip>
          <HintChip>Pinch / scroll zoom</HintChip>
          {canEdit && <HintChip accent>Drag to connect</HintChip>}
        </div>
      </div>

      <div className="flex shrink-0 items-center gap-1 self-end sm:self-auto">
        <Button type="button" variant="secondary" size="sm" onClick={onZoomOut} aria-label="Zoom out">
          −
        </Button>
        <span className="min-w-12 text-center text-xs font-medium tabular-nums text-slate-600">
          {Math.round(scale * 100)}%
        </span>
        <Button type="button" variant="secondary" size="sm" onClick={onZoomIn} aria-label="Zoom in">
          +
        </Button>
        <Button type="button" variant="ghost" size="sm" onClick={onResetView}>
          Reset
        </Button>
      </div>
    </div>
  );
}

function HintChip({ children, accent = false }: { children: ReactNode; accent?: boolean }) {
  return (
    <span
      className={cn(
        "rounded-full px-2 py-0.5 text-[11px] font-medium",
        accent ? "bg-emerald-50 text-emerald-700" : "bg-slate-100 text-slate-600",
      )}
    >
      {children}
    </span>
  );
}
