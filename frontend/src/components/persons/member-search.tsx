"use client";

import { useDeferredValue, useState } from "react";

import { Alert } from "@/components/ui/alert";
import { Spinner } from "@/components/ui/spinner";
import { useFamilyPersons } from "@/hooks/use-persons";
import type { Person } from "@/types/person";
import { displayName } from "@/utils/tree-layout";
import { cn } from "@/utils/cn";

interface MemberSearchProps {
  familyId: string;
  onSelect: (person: Person) => void;
  selectedId?: string | null;
}

/** Server-backed name search (pg_trgm) for finding members in a large tree. */
export function MemberSearch({ familyId, onSelect, selectedId }: MemberSearchProps) {
  const [query, setQuery] = useState("");
  const deferredQuery = useDeferredValue(query);
  const results = useFamilyPersons(familyId, deferredQuery.trim(), 1);

  return (
    <div className="aceh-panel space-y-4">
      <div>
        <h2 className="font-display text-2xl text-teal-950">Find a member</h2>
        <p className="mt-1 text-base text-slate-600">
          Type part of a name to search the tree. Tap a result to open details.
        </p>
      </div>

      <input
        type="search"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search by name…"
        className="block w-full rounded-xl border-0 bg-white px-4 py-3.5 text-lg text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-teal-700"
      />

      {deferredQuery.trim().length === 0 ? (
        <p className="text-base text-slate-500">Start typing to search.</p>
      ) : results.isPending ? (
        <div className="grid place-items-center py-6">
          <Spinner className="size-8 text-teal-700" />
        </div>
      ) : results.isError ? (
        <Alert tone="error">Could not search right now. Try again.</Alert>
      ) : results.data && results.data.items.length > 0 ? (
        <ul className="max-h-72 overflow-y-auto rounded-xl bg-white ring-1 ring-slate-200">
          {results.data.items.map((p) => (
            <li key={p.id}>
              <button
                type="button"
                className={cn(
                  "flex w-full items-center gap-3 border-b border-slate-100 px-4 py-3.5 text-left text-lg last:border-b-0",
                  selectedId === p.id
                    ? "bg-teal-800 font-semibold text-white"
                    : "text-slate-900 hover:bg-teal-50",
                )}
                onClick={() => onSelect(p)}
              >
                {p.photo_url ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={p.photo_url}
                    alt=""
                    className="size-10 shrink-0 rounded-full object-cover ring-1 ring-teal-900/10"
                  />
                ) : (
                  <span
                    className={cn(
                      "grid size-10 shrink-0 place-items-center rounded-full text-sm font-bold",
                      selectedId === p.id ? "bg-teal-700 text-amber-100" : "bg-teal-100 text-teal-900",
                    )}
                  >
                    {displayName(p).slice(0, 1).toUpperCase()}
                  </span>
                )}
                <span>{displayName(p)}</span>
              </button>
            </li>
          ))}
        </ul>
      ) : (
        <p className="text-base text-slate-500">No members match “{deferredQuery.trim()}”.</p>
      )}
    </div>
  );
}
