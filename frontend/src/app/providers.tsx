"use client";

import { useState, type ReactNode } from "react";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

import { makeQueryClient } from "@/lib/query-client";

/** Client-side context providers wrapping the whole app. */
export function Providers({ children }: { children: ReactNode }) {
  // useState keeps one QueryClient instance across re-renders.
  const [queryClient] = useState(makeQueryClient);

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {/* No-op in production builds. */}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
