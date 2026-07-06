/**
 * Shapes of the backend's JSON envelope (see docs/api/README.md).
 * Every service unwraps these; components never see envelopes.
 */

export interface ApiEnvelope<T> {
  success: boolean;
  data?: T;
  error?: ApiErrorPayload;
  meta?: PaginationMeta;
}

export interface ApiErrorPayload {
  /** Stable machine-readable code, e.g. "VALIDATION_ERROR". */
  code: string;
  message: string;
  /** For VALIDATION_ERROR: map of snake_case field -> message. */
  details?: Record<string, string>;
}

export interface PaginationMeta {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

/** Normalized error thrown by the service layer. */
export class ApiError extends Error {
  readonly code: string;
  readonly status: number;
  readonly details?: Record<string, string>;

  constructor(status: number, payload?: ApiErrorPayload) {
    super(payload?.message ?? "Something went wrong");
    this.name = "ApiError";
    this.status = status;
    this.code = payload?.code ?? "UNKNOWN";
    this.details = payload?.details;
  }
}
