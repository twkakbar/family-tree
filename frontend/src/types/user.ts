/** Global account role. Per-family roles arrive with Phase 2. */
export type Role = "admin" | "user";

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  email_verified: boolean;
  created_at: string;
}
