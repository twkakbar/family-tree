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
      <div className="grid min-h-dvh place-items-center">
        <Spinner className="size-8 text-emerald-600" />
      </div>
    );
  }

  return (
    <div className="min-h-dvh">
      <Header />
      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6">{children}</main>
    </div>
  );
}
