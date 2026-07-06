"use client";

import { useLogout, useMe } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import { Logo } from "./logo";

/** Top bar of the authenticated app shell. */
export function Header() {
  const { data: user } = useMe();
  const logout = useLogout();

  return (
    <header className="border-b border-slate-200 bg-white">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6">
        <Logo />
        <div className="flex items-center gap-4">
          {user && (
            <div className="hidden text-right sm:block">
              <p className="text-sm font-medium text-slate-900">{user.name}</p>
              <p className="text-xs text-slate-500">{user.email}</p>
            </div>
          )}
          <Button
            variant="secondary"
            size="sm"
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
