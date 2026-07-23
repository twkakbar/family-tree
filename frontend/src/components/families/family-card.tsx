"use client";

import Link from "next/link";

import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import type { Family } from "@/types/family";

const roleLabels: Record<Family["role"], string> = {
  owner: "Owner",
  editor: "Editor",
  viewer: "Viewer",
};

interface FamilyCardProps {
  family: Family;
}

/** Summary card linking to a family detail page. */
export function FamilyCard({ family }: FamilyCardProps) {
  return (
    <Link href={`/families/${family.id}`} className="block transition hover:-translate-y-0.5">
      <Card className="h-full hover:ring-emerald-200">
        <div className="flex items-start justify-between gap-3">
          <CardTitle>{family.name}</CardTitle>
          <span className="shrink-0 rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-medium text-emerald-700">
            {roleLabels[family.role]}
          </span>
        </div>
        {family.description ? (
          <CardDescription className="line-clamp-2">{family.description}</CardDescription>
        ) : (
          <CardDescription>No description yet.</CardDescription>
        )}
        <p className="mt-4 text-xs text-slate-500">
          {family.persons_count} member{family.persons_count === 1 ? "" : "s"}
        </p>
      </Card>
    </Link>
  );
}
