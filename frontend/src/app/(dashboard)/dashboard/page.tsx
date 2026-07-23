"use client";

import Link from "next/link";

import { useMe } from "@/hooks/use-auth";
import { useStudioCanEdit } from "@/hooks/use-public";

/** Quiet account hub for signed-in users (no public marketing of auth). */
export default function DashboardPage() {
  const { data: user } = useMe();
  const canEdit = useStudioCanEdit();

  return (
    <div className="mx-auto max-w-lg space-y-6">
      <div>
        <p className="text-sm font-semibold uppercase tracking-wider text-teal-800/70">Account</p>
        <h1 className="font-display mt-1 text-4xl text-teal-950">
          Hello{user ? `, ${user.name}` : ""}
        </h1>
        <p className="mt-3 text-lg text-slate-600">Choose where to go next.</p>
      </div>

      <div className="grid gap-3">
        <Link
          href="/"
          className="rounded-2xl bg-white/90 px-5 py-4 text-lg font-semibold text-teal-950 ring-1 ring-teal-900/10 hover:ring-teal-600/40"
        >
          View the public tree
        </Link>
        {canEdit.data?.can_edit ? (
          <Link
            href="/studio"
            className="rounded-2xl bg-teal-800 px-5 py-4 text-lg font-semibold text-white hover:bg-teal-900"
          >
            Open studio — add members
          </Link>
        ) : (
          <Link
            href="/request-access"
            className="rounded-2xl bg-teal-800 px-5 py-4 text-lg font-semibold text-white hover:bg-teal-900"
          >
            Ask for edit access
          </Link>
        )}
        {user?.role === "admin" && (
          <>
            <Link
              href="/admin/requests"
              className="rounded-2xl bg-white/90 px-5 py-4 text-lg font-semibold text-teal-950 ring-1 ring-teal-900/10 hover:ring-teal-600/40"
            >
              Review edit requests
            </Link>
            <Link
              href="/admin/users"
              className="rounded-2xl bg-white/90 px-5 py-4 text-lg font-semibold text-teal-950 ring-1 ring-teal-900/10 hover:ring-teal-600/40"
            >
              Manage accounts
            </Link>
          </>
        )}
      </div>
    </div>
  );
}
