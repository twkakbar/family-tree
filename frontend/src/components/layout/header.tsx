"use client";

import Link from "next/link";

import { useLogout, useMe } from "@/hooks/use-auth";
import { useStudioCanEdit } from "@/hooks/use-public";
import { Button } from "@/components/ui/button";

/** Top bar of the authenticated app shell — Aceh sea/gold chrome. */
export function Header() {
  const { data: user } = useMe();
  const logout = useLogout();
  const canEdit = useStudioCanEdit();

  return (
    <header className="border-b border-teal-900/10 bg-[#04151a]/90 text-teal-50 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between gap-4 px-4 sm:px-6">
        <nav className="flex items-center gap-4 text-base font-medium sm:gap-5">
          <Link href="/" className="font-display text-xl tracking-wide text-amber-200">
            Kesultanan Aceh
          </Link>
          <Link href="/about" className="hidden text-teal-100/80 hover:text-amber-100 sm:inline">
            About
          </Link>
          {canEdit.data?.can_edit && (
            <Link href="/studio" className="hidden text-teal-100/80 hover:text-amber-100 sm:inline">
              Studio
            </Link>
          )}
          {user?.role === "admin" && (
            <>
              <Link
                href="/admin/requests"
                className="hidden text-teal-100/80 hover:text-amber-100 sm:inline"
              >
                Approvals
              </Link>
              <Link
                href="/admin/users"
                className="hidden text-teal-100/80 hover:text-amber-100 sm:inline"
              >
                Accounts
              </Link>
            </>
          )}
        </nav>
        <div className="flex items-center gap-3 sm:gap-4">
          {user && (
            <div className="hidden text-right sm:block">
              <p className="text-sm font-medium text-amber-50">{user.name}</p>
              <p className="text-xs text-teal-200/60">{user.email}</p>
            </div>
          )}
          <Button
            variant="secondary"
            size="sm"
            className="border-0 bg-white/10 text-teal-50 ring-teal-400/30 hover:bg-white/20"
            loading={logout.isPending}
            onClick={() => logout.mutate()}
          >
            Sign out
          </Button>
        </div>
      </div>
    </header>
  );
}
