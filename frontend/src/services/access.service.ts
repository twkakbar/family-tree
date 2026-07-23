import { apiClient, unwrap } from "@/lib/api-client";
import type { ApiEnvelope } from "@/types/api";

export type AccessRequestStatus = "pending" | "approved" | "rejected";

export interface AccessRequest {
  id: string;
  status: AccessRequestStatus;
  message: string;
  created_at: string;
  reviewed_at?: string | null;
}

export interface AccessRequestAdmin extends AccessRequest {
  user_id: string;
  name: string;
  email: string;
}

export interface StudioCanEdit {
  can_edit: boolean;
  family_id: string;
}

export const accessService = {
  async submit(message: string): Promise<AccessRequest> {
    const res = await apiClient.post<ApiEnvelope<AccessRequest>>("/access-requests", { message });
    return unwrap(res.data);
  },

  async myLatest(): Promise<AccessRequest> {
    const res = await apiClient.get<ApiEnvelope<AccessRequest>>("/access-requests/me");
    return unwrap(res.data);
  },

  async canEdit(): Promise<StudioCanEdit> {
    const res = await apiClient.get<ApiEnvelope<StudioCanEdit>>("/studio/can-edit");
    return unwrap(res.data);
  },

  async listPending(): Promise<AccessRequestAdmin[]> {
    const res = await apiClient.get<ApiEnvelope<AccessRequestAdmin[]>>("/admin/access-requests");
    return unwrap(res.data);
  },

  async approve(id: string): Promise<AccessRequest> {
    const res = await apiClient.post<ApiEnvelope<AccessRequest>>(
      `/admin/access-requests/${id}/approve`,
    );
    return unwrap(res.data);
  },

  async reject(id: string): Promise<AccessRequest> {
    const res = await apiClient.post<ApiEnvelope<AccessRequest>>(
      `/admin/access-requests/${id}/reject`,
    );
    return unwrap(res.data);
  },
};
