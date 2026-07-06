/**
 * In-memory store for the JWT access token.
 *
 * Deliberately NOT localStorage/sessionStorage (XSS-exfiltration risk,
 * DEVELOPMENT_RULES §4/§7). A page reload loses the token; the axios client
 * then silently recovers the session via the httpOnly refresh cookie.
 */
let accessToken: string | null = null;

export const authToken = {
  get: (): string | null => accessToken,
  set: (token: string): void => {
    accessToken = token;
  },
  clear: (): void => {
    accessToken = null;
  },
};
