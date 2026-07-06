import { apiClient, unwrap } from "@/lib/api-client";
import type { ApiEnvelope } from "@/types/api";
import type { User } from "@/types/user";

/** User profile API calls (docs/api/users.md). */
export const usersService = {
  async getMe(): Promise<User> {
    const res = await apiClient.get<ApiEnvelope<User>>("/users/me");
    return unwrap(res.data);
  },

  async updateMe(input: { name: string }): Promise<User> {
    const res = await apiClient.patch<ApiEnvelope<User>>("/users/me", input);
    return unwrap(res.data);
  },
};
