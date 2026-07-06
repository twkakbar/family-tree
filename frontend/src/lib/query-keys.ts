/**
 * Central registry of React Query cache keys (DEVELOPMENT_RULES §4).
 * Never build query keys inline; add new domains here.
 */
export const queryKeys = {
  auth: {
    me: ["auth", "me"] as const,
  },
  // Phase 2+: families, persons, relationships, tree…
} as const;
