"use client";

import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { useMe } from "@/hooks/use-auth";

/**
 * Dashboard landing. Phase 1 shows the roadmap; Phase 2 replaces the
 * placeholder cards with the family list (see CONTEXT.md).
 */
const upcoming = [
  { title: "Families", phase: "Phase 2", description: "Create families and invite collaborators." },
  { title: "Family members", phase: "Phase 3", description: "Add persons and define relationships." },
  { title: "Tree view", phase: "Phase 3", description: "Explore your interactive family tree." },
  { title: "Photos & search", phase: "Phase 4", description: "Profile pictures and member search." },
  { title: "Export / import", phase: "Phase 5", description: "Back up or move family trees as JSON." },
];

export default function DashboardPage() {
  const { data: user } = useMe();

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900">
        Welcome{user ? `, ${user.name}` : ""}
      </h1>
      <p className="mt-1 text-slate-600">
        Your account is ready{user?.role === "admin" ? " (administrator)" : ""}. Family tree
        features are on the way — here is what lands next.
      </p>

      <div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {upcoming.map((item) => (
          <Card key={item.title} className="relative opacity-80">
            <span className="absolute right-4 top-4 rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-600">
              {item.phase}
            </span>
            <CardTitle>{item.title}</CardTitle>
            <CardDescription>{item.description}</CardDescription>
          </Card>
        ))}
      </div>
    </div>
  );
}
