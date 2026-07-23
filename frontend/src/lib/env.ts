/**
 * The single place reading environment variables in the frontend
 * (DEVELOPMENT_RULES §4). NEXT_PUBLIC_* values are inlined at build time;
 * defaults keep every environment working without configuration.
 */
export const env = {
  apiBaseUrl: process.env.NEXT_PUBLIC_API_BASE_URL ?? "/api/v1",
  appName: process.env.NEXT_PUBLIC_APP_NAME ?? "Kesultanan Aceh",
} as const;
