"use client";

import { useMemo, useState } from "react";

import type { Person } from "@/types/person";
import { displayName } from "@/utils/tree-layout";
import { cn } from "@/utils/cn";

interface PersonSearchPickerProps {
  id?: string;
  people: Person[];
  value: string;
  onChange: (personId: string) => void;
  error?: string;
}

/**
 * Large search + tap-to-pick list for finding a person among many names.
 */
export function PersonSearchPicker({
  id = "person-search",
  people,
  value,
  onChange,
  error,
}: PersonSearchPickerProps) {
  const [query, setQuery] = useState("");

  const selected = useMemo(
    () => people.find((p) => p.id === value) ?? null,
    [people, value],
  );

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    const sorted = people
      .slice()
      .sort((a, b) => displayName(a).localeCompare(displayName(b), undefined, { sensitivity: "base" }));
    if (!q) {
      return sorted;
    }
    return sorted.filter((p) => displayName(p).toLowerCase().includes(q));
  }, [people, query]);

  return (
    <div className="mt-2 space-y-3">
      {selected && (
        <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-teal-800/10 px-4 py-3 ring-1 ring-teal-800/20">
          <p className="text-lg text-teal-950">
            Selected: <strong>{displayName(selected)}</strong>
          </p>
          <button
            type="button"
            className="text-base font-semibold text-teal-800 underline-offset-2 hover:underline"
            onClick={() => {
              onChange("");
              setQuery("");
            }}
          >
            Clear
          </button>
        </div>
      )}

      <div>
        <label htmlFor={id} className="sr-only">
          Search names
        </label>
        <input
          id={id}
          type="search"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Type to search a name…"
          autoComplete="off"
          className={cn(
            "block w-full rounded-xl border-0 bg-white px-4 py-3.5 text-lg text-slate-900 shadow-sm",
            "ring-1 ring-inset ring-slate-300 placeholder:text-slate-400",
            "focus:ring-2 focus:ring-inset focus:ring-teal-700",
            error && "ring-red-400 focus:ring-red-500",
          )}
        />
        <p className="mt-1.5 text-sm text-slate-500">
          {filtered.length === people.length
            ? `${people.length} people`
            : `${filtered.length} of ${people.length} match`}
        </p>
      </div>

      <ul
        role="listbox"
        aria-label="Matching names"
        className="max-h-64 overflow-y-auto rounded-xl bg-white ring-1 ring-slate-200"
      >
        {filtered.length === 0 ? (
          <li className="px-4 py-4 text-base text-slate-500">No names match that search.</li>
        ) : (
          filtered.map((p) => {
            const name = displayName(p);
            const isSelected = p.id === value;
            return (
              <li key={p.id} role="option" aria-selected={isSelected}>
                <button
                  type="button"
                  className={cn(
                    "w-full border-b border-slate-100 px-4 py-3.5 text-left text-lg last:border-b-0",
                    isSelected
                      ? "bg-teal-800 font-semibold text-white"
                      : "text-slate-900 hover:bg-teal-50",
                  )}
                  onClick={() => onChange(p.id)}
                >
                  {name}
                </button>
              </li>
            );
          })
        )}
      </ul>

      {error && <p className="text-sm text-red-600">{error}</p>}
    </div>
  );
}
