import axios, { AxiosError, type InternalAxiosRequestConfig } from "axios";

import { ApiError, type ApiEnvelope } from "@/types/api";
import { authToken } from "./auth-token";
import { env } from "./env";

/**
 * The single axios instance behind the entire service layer.
 *
 * Responsibilities:
 *  - attach the in-memory access token to every request
 *  - on 401, perform ONE single-flight cookie-based refresh and retry once
 *  - normalize every failure into ApiError
 *
 * Components and hooks never import this directly - services do
 * (DEVELOPMENT_RULES §4).
 */
export const apiClient = axios.create({
  baseURL: env.apiBaseUrl,
  withCredentials: true, // refresh cookie
  timeout: 15_000,
  headers: { "Content-Type": "application/json" },
});

interface RetriableRequestConfig extends InternalAxiosRequestConfig {
  _retried?: boolean;
}

apiClient.interceptors.request.use((config) => {
  const token = authToken.get();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

/** Endpoints that must never trigger a refresh-retry loop. */
const AUTH_PATHS = ["/auth/login", "/auth/register", "/auth/refresh", "/auth/logout"];

const isAuthPath = (url: string | undefined): boolean =>
  AUTH_PATHS.some((p) => url?.includes(p));

/** Single-flight guard: concurrent 401s share one refresh request. */
let refreshPromise: Promise<string> | null = null;

async function refreshAccessToken(): Promise<string> {
  refreshPromise ??= axios
    .post<ApiEnvelope<{ access_token: string }>>(
      `${env.apiBaseUrl}/auth/refresh`,
      null,
      { withCredentials: true, timeout: 15_000 },
    )
    .then((res) => {
      const token = res.data.data?.access_token;
      if (!token) throw new Error("refresh response missing access token");
      authToken.set(token);
      return token;
    })
    .finally(() => {
      refreshPromise = null;
    });
  return refreshPromise;
}

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiEnvelope<unknown>>) => {
    const config = error.config as RetriableRequestConfig | undefined;
    const status = error.response?.status ?? 0;

    const shouldTryRefresh =
      status === 401 && config !== undefined && !config._retried && !isAuthPath(config.url);

    if (shouldTryRefresh) {
      try {
        const token = await refreshAccessToken();
        config._retried = true;
        config.headers.Authorization = `Bearer ${token}`;
        return apiClient(config);
      } catch {
        authToken.clear();
        redirectToLogin();
      }
    }

    throw new ApiError(status, error.response?.data?.error);
  },
);

/** Hard redirect keeps this module free of router coupling. */
function redirectToLogin(): void {
  if (typeof window === "undefined") return;
  const publicPaths = ["/", "/login", "/register"];
  if (!publicPaths.includes(window.location.pathname)) {
    window.location.assign("/login");
  }
}

/** Unwraps the success envelope; services use this for every call. */
export function unwrap<T>(envelope: ApiEnvelope<T>): T {
  if (!envelope.success || envelope.data === undefined) {
    throw new ApiError(0, envelope.error);
  }
  return envelope.data;
}
