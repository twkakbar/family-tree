import { apiClient, unwrap } from "@/lib/api-client";
import { ApiError, type ApiEnvelope, type PaginationMeta } from "@/types/api";
import type { Role, User } from "@/types/user";

export interface AdminUserListResult {
  items: User[];
  meta: PaginationMeta;
}

export interface AdminUpdateUserInput {
  name?: string;
  role?: Role;
  email_verified?: boolean;
}

export interface PasswordResetIssueResult {
  reset_token: string;
  expires_in: number;
  hint: string;
}

/** Super-admin user management (docs/api/admin-users.md). */
export const adminUsersService = {
  async list(search = "", page = 1, perPage = 20): Promise<AdminUserListResult> {
    const res = await apiClient.get<ApiEnvelope<User[]>>("/admin/users", {
      params: { search: search || undefined, page, per_page: perPage },
    });
    const envelope = res.data;
    if (!envelope.success || envelope.data === undefined || !envelope.meta) {
      throw new ApiError(res.status, envelope.error);
    }
    return { items: envelope.data, meta: envelope.meta };
  },

  async update(id: string, input: AdminUpdateUserInput): Promise<User> {
    const res = await apiClient.patch<ApiEnvelope<User>>(`/admin/users/${id}`, input);
    return unwrap(res.data);
  },

  async remove(id: string): Promise<void> {
    await apiClient.delete(`/admin/users/${id}`);
  },

  async issuePasswordReset(id: string): Promise<PasswordResetIssueResult> {
    const res = await apiClient.post<ApiEnvelope<PasswordResetIssueResult>>(
      `/admin/users/${id}/password-reset`,
    );
    return unwrap(res.data);
  },
};
