import type { User } from "./user";

/** Payload of register/login/refresh responses (docs/api/auth.md). */
export interface AuthResult {
  access_token: string;
  token_type: "Bearer";
  /** Seconds until the access token expires. */
  expires_in: number;
  user: User;
}

export interface RegisterInput {
  name: string;
  email: string;
  password: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface VerificationRequired {
  verification_required: true;
  message: string;
}

export function isAuthResult(value: AuthResult | VerificationRequired): value is AuthResult {
  return "access_token" in value;
}
