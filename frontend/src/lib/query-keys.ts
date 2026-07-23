/**
 * Central registry of React Query cache keys (DEVELOPMENT_RULES §4).
 */
export const queryKeys = {
  auth: {
    me: ["auth", "me"] as const,
  },
  public: {
    meta: ["public", "meta"] as const,
    tree: ["public", "tree"] as const,
    person: (id: string) => ["public", "person", id] as const,
  },
  access: {
    me: ["access", "me"] as const,
  },
  studio: {
    canEdit: ["studio", "can-edit"] as const,
  },
  admin: {
    accessRequests: ["admin", "access-requests"] as const,
    users: (search: string, page: number) => ["admin", "users", search, page] as const,
  },
  transfer: {
    export: (familyId: string) => ["transfer", "export", familyId] as const,
  },
  families: {
    all: ["families"] as const,
    list: (page: number, perPage: number) => ["families", "list", page, perPage] as const,
    detail: (id: string) => ["families", "detail", id] as const,
    collaborators: (id: string) => ["families", "collaborators", id] as const,
    tree: (id: string) => ["families", "tree", id] as const,
    persons: (id: string, search: string, page: number) =>
      ["families", "persons", id, search, page] as const,
  },
  persons: {
    detail: (id: string) => ["persons", "detail", id] as const,
    ancestors: (id: string, depth: number) => ["persons", "ancestors", id, depth] as const,
    descendants: (id: string, depth: number) => ["persons", "descendants", id, depth] as const,
  },
} as const;
