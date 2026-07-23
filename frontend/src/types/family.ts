export type FamilyRole = "owner" | "editor" | "viewer";

export interface Family {
  id: string;
  name: string;
  description: string;
  role: FamilyRole;
  persons_count: number;
  created_at: string;
  updated_at: string;
}

export interface Collaborator {
  user_id: string;
  name: string;
  email: string;
  role: FamilyRole;
  created_at: string;
}

export interface CreateFamilyInput {
  name: string;
  description?: string;
}

export interface UpdateFamilyInput {
  name: string;
  description?: string;
}

export interface InviteCollaboratorInput {
  email: string;
  role: "editor" | "viewer";
}

export interface UpdateCollaboratorInput {
  role: FamilyRole;
}

export interface FamilyListResult {
  items: Family[];
  meta: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}
