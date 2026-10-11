import { apiClient, unwrap } from "@/lib/api-client";
import type { ApiEnvelope } from "@/types/api";
import type { AuthResult, LoginInput, RegisterInput, VerificationRequired } from "@/types/auth";

/**
 * Auth API calls (docs/api/auth.md). The refresh token never appears here:
 * it lives in an httpOnly cookie managed by the backend.
 */
export const authService = {
  async register(input: RegisterInput): Promise<AuthResult | VerificationRequired> {
    const res = await apiClient.post<ApiEnvelope<AuthResult | VerificationRequired>>(
      "/auth/register",
      input,
    );
    return unwrap(res.data);
  },

  async login(input: LoginInput): Promise<AuthResult> {
    const res = await apiClient.post<ApiEnvelope<AuthResult>>("/auth/login", input);
    return unwrap(res.data);
  },

  async logout(): Promise<void> {
    await apiClient.post("/auth/logout");
  },

  async forgotPassword(email: string): Promise<{
    requested: boolean;
    message: string;
    reset_token?: string;
  }> {
    const res = await apiClient.post<
      ApiEnvelope<{ requested: boolean; message: string; reset_token?: string }>
    >("/auth/forgot-password", { email });
    return unwrap(res.data);
  },

  async resetPassword(token: string, newPassword: string): Promise<void> {
    await apiClient.post("/auth/reset-password", {
      token,
      new_password: newPassword,
    });
  },

  async verifyEmail(token: string): Promise<void> {
    await apiClient.post("/auth/verify-email", { token });
  },

  async resendVerification(email: string): Promise<void> {
    await apiClient.post("/auth/resend-verification", { email });
  },
};
