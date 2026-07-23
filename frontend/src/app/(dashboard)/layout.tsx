"use client";

import { useEffect, type ReactNode } from "react";
import { useRouter } from "next/navigation";

import { Header } from "@/components/layout/header";
import { Spinner } from "@/components/ui/spinner";
import { useMe } from "@/hooks/use-auth";

/**
 * Authenticated app shell. src/middleware.ts already redirects visitors
 * without a session cookie; this guard covers expired/invalid sessions by
 * actually loading the profile through the API.
 */
export default function DashboardLayout({ children }: { children: ReactNode }) {
  const { data: user, isPending, isError } = useMe();
  const router = useRouter();

  useEffect(() => {
    if (isError) {
      router.replace("/login");
    }
  }, [isError, router]);

  if (isPending || isError || !user) {
    return (
      <div className="aceh-shell grid min-h-dvh place-items-center">
        <Spinner className="size-10 text-teal-700" />
      </div>
    );
  }

  return (
    <div className="aceh-shell">
      <Header />
      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 sm:py-10">{children}</main>
    </div>
  );
}
