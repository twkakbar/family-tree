import type { PointerEvent } from "react";

import type { Person, TreeNodePosition } from "@/types/person";
import { displayName } from "@/utils/tree-layout";
import { cn } from "@/utils/cn";

interface PersonTreeNodeProps {
  node: TreeNodePosition;
  canEdit: boolean;
  isSource: boolean;
  isHoverTarget: boolean;
  isSelected: boolean;
  isLinking: boolean;
  width: number;
  height: number;
  onPointerDown: (event: PointerEvent) => void;
  onSelect: () => void;
}

function personMeta(person: Person): string {
  const parts: string[] = [];
  if (person.gender !== "unknown") {
    parts.push(person.gender);
  }
  if (person.birth_date) {
    parts.push(`b. ${person.birth_date.slice(0, 4)}`);
  }
  return parts.join(" · ");
}

/** A single person card on the family tree canvas. */
export function PersonTreeNode({
  node,
  canEdit,
  isSource,
  isHoverTarget,
  isSelected,
  isLinking,
  width,
  height,
  onPointerDown,
  onSelect,
}: PersonTreeNodeProps) {
  const { person } = node;
  const meta = personMeta(person);
  const fullName = displayName(person);

  return (
    <div
      data-person-node
      role="button"
      tabIndex={0}
      title={fullName}
      className={cn(
        "group absolute select-none rounded-xl border bg-white text-left shadow-sm transition-all",
        canEdit && !isLinking && "cursor-grab active:cursor-grabbing active:scale-[0.98]",
        isSource && "z-10 border-emerald-500 shadow-md ring-2 ring-emerald-500/30",
        isHoverTarget && "z-10 border-emerald-400 bg-emerald-50 shadow-md ring-2 ring-emerald-400/40",
        isSelected && !isLinking && "border-emerald-600 ring-2 ring-emerald-500/20",
        !isSource && !isHoverTarget && !isSelected && "border-slate-200 hover:border-emerald-300 hover:shadow",
      )}
      style={{
        width,
        height,
        left: node.x,
        top: node.y,
        touchAction: canEdit ? "none" : "auto",
      }}
      onPointerDown={onPointerDown}
      onKeyDown={(event) => {
        if (event.key === "Enter" || event.key === " ") {
          event.preventDefault();
          onSelect();
        }
      }}
    >
      {canEdit && !isLinking && (
        <span
          className="pointer-events-none absolute -right-1.5 top-1/2 flex size-5 -translate-y-1/2 items-center justify-center rounded-full border border-emerald-200 bg-emerald-50 text-emerald-600 opacity-0 shadow-sm transition-opacity group-hover:opacity-100"
          aria-hidden
        >
          <LinkIcon />
        </span>
      )}

      <div className="flex h-full min-h-0 flex-col px-3 py-2">
        <div className="min-h-0 flex-1 overflow-hidden">
          <p className="truncate text-sm font-semibold leading-5 text-slate-900">{person.first_name}</p>
          {person.last_name ? (
            <p className="truncate text-xs leading-4 text-slate-700">{person.last_name}</p>
          ) : null}
        </div>
        <p className="mt-1 shrink-0 truncate text-xs leading-4 text-slate-500">{meta || "Member"}</p>
      </div>
    </div>
  );
}

function LinkIcon() {
  return (
    <svg viewBox="0 0 16 16" className="size-3" fill="none" stroke="currentColor" strokeWidth="1.75">
      <path d="M6 8h4M8 6v4" strokeLinecap="round" />
      <circle cx="4" cy="8" r="2" />
      <circle cx="12" cy="8" r="2" />
    </svg>
  );
}
