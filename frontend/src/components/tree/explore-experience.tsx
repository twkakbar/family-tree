"use client";

import dynamic from "next/dynamic";
import Link from "next/link";
import { useState } from "react";

import { Spinner } from "@/components/ui/spinner";
import { TreeLegend } from "@/components/tree/tree-legend";
import { usePublicPerson, usePublicTree, useRoyalMeta } from "@/hooks/use-public";
import type { Person } from "@/types/person";
import { displayName } from "@/utils/tree-layout-3d";

const RoyalTreeCanvas = dynamic(
  () => import("@/components/tree/royal-tree-canvas").then((m) => m.RoyalTreeCanvas),
  {
    ssr: false,
    loading: () => (
      <div className="aceh-sea-bg grid h-full place-items-center">
        <Spinner className="size-8 text-teal-300" />
      </div>
    ),
  },
);

export function ExploreExperience() {
  const meta = useRoyalMeta();
  const tree = usePublicTree();
  const [selected, setSelected] = useState<Person | null>(null);
  const detail = usePublicPerson(selected?.id ?? null);

  return (
    <div className="aceh-sea-bg relative h-dvh w-full overflow-hidden text-teal-50">
      <header className="pointer-events-none absolute inset-x-0 top-0 z-30 flex items-start justify-between gap-3 p-3 sm:gap-4 sm:p-6">
        <div className="pointer-events-auto max-w-md rounded-2xl bg-slate-950/55 px-4 py-3 shadow-xl ring-1 ring-teal-500/20 backdrop-blur-md">
          <p className="font-display text-lg tracking-wide text-amber-200 sm:text-xl">
            {meta.data?.name ?? "Kesultanan Aceh"}
          </p>
          <p className="mt-1 text-xs text-teal-100/70 sm:text-sm">
            Drag to orbit · scroll to zoom · tap a member
          </p>
          <div className="mt-2">
            <TreeLegend tone="dark" />
          </div>
        </div>
        <Link
          href="/about"
          className="pointer-events-auto rounded-full bg-teal-950/70 px-5 py-2.5 text-base font-medium text-teal-100 ring-1 ring-teal-400/30 backdrop-blur hover:bg-teal-900/80"
        >
          About
        </Link>
      </header>

      <div className="absolute inset-0">
        {tree.isPending ? (
          <div className="grid h-full place-items-center">
            <Spinner className="size-8 text-teal-300" />
          </div>
        ) : tree.data && tree.data.persons.length > 0 ? (
          <RoyalTreeCanvas
            data={tree.data}
            selectedId={selected?.id}
            onSelect={setSelected}
            onDismiss={() => setSelected(null)}
          />
        ) : (
          <div className="grid h-full place-items-center px-6 text-center">
            <div>
              <p className="font-display text-2xl text-amber-200">The tree awaits its first names</p>
              <p className="mt-2 max-w-md text-sm text-teal-100/70">
                Editors will add members from living descendants up to the sultan line.
              </p>
            </div>
          </div>
        )}
      </div>

      {selected && (
        <aside className="absolute inset-x-0 bottom-0 z-40 max-h-[75dvh] overflow-y-auto rounded-t-3xl bg-[#071c22] p-4 pb-[max(1rem,env(safe-area-inset-bottom))] shadow-2xl ring-1 ring-teal-400/30 sm:inset-x-auto sm:bottom-6 sm:right-6 sm:max-h-[min(70vh,28rem)] sm:w-[22rem] sm:rounded-2xl sm:bg-slate-950/95 sm:pb-4">
          <div className="flex items-start justify-between gap-3">
            <div>
              <h2 className="font-display text-xl text-amber-100">{displayName(selected)}</h2>
              <p className="mt-1 text-xs capitalize text-teal-200/70">{selected.gender}</p>
            </div>
            <button
              type="button"
              className="text-sm text-teal-300/80 hover:text-teal-100"
              onClick={() => setSelected(null)}
            >
              Close
            </button>
          </div>
          {(selected.birth_date || selected.death_date) && (
            <p className="mt-3 text-sm text-teal-100/80">
              {selected.birth_date ?? "?"} — {selected.death_date ?? "living"}
            </p>
          )}
          {(detail.data?.photo_url || selected.photo_url) && (
            <div className="mt-4 overflow-hidden rounded-xl bg-teal-950/50 ring-1 ring-teal-400/20">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={detail.data?.photo_url || selected.photo_url}
                alt={displayName(selected)}
                className="mx-auto max-h-72 w-full object-contain"
              />
            </div>
          )}
          {selected.bio && <p className="mt-3 text-sm leading-relaxed text-teal-50/80">{selected.bio}</p>}
          {detail.isPending && (
            <p className="mt-4 text-xs text-teal-300/60">Loading relationships…</p>
          )}
          {detail.data && detail.data.relationships.length > 0 && (
            <p className="mt-4 text-xs text-teal-200/60">
              {detail.data.relationships.length} relationship
              {detail.data.relationships.length === 1 ? "" : "s"} on record
            </p>
          )}
        </aside>
      )}
    </div>
  );
}
