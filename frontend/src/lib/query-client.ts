import { QueryClient } from "@tanstack/react-query";

import { ApiError } from "@/types/api";

/**
 * React Query defaults. Auth failures (401/403) never retry: retrying cannot
 * fix them and the api-client already handled token refresh.
 */
export function makeQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 30_000,
        retry: (failureCount, error) => {
          if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
            return false;
          }
          return failureCount < 2;
        },
        refetchOnWindowFocus: false,
      },
    },
  });
}
