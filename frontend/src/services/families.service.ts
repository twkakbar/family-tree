import { apiClient, unwrap } from "@/lib/api-client";
import { ApiError, type ApiEnvelope } from "@/types/api";
import type {
  Collaborator,
  CreateFamilyInput,
  Family,
  FamilyListResult,
  InviteCollaboratorInput,
  UpdateCollaboratorInput,
  UpdateFamilyInput,
} from "@/types/family";

/**
 * Families API (docs/api/families.md).
 */
export const familiesService = {
  async list(page = 1, perPage = 20): Promise<FamilyListResult> {
    const res = await apiClient.get<ApiEnvelope<Family[]>>("/families", {
      params: { page, per_page: perPage },
    });
    const envelope = res.data;
    if (!envelope.success || envelope.data === undefined || !envelope.meta) {
      throw new ApiError(res.status, envelope.error);
    }
    return { items: envelope.data, meta: envelope.meta };
  },

  async getById(id: string): Promise<Family> {
    const res = await apiClient.get<ApiEnvelope<Family>>(`/families/${id}`);
    return unwrap(res.data);
  },

  async create(input: CreateFamilyInput): Promise<Family> {
    const res = await apiClient.post<ApiEnvelope<Family>>("/families", input);
    return unwrap(res.data);
  },

  async update(id: string, input: UpdateFamilyInput): Promise<Family> {
    const res = await apiClient.patch<ApiEnvelope<Family>>(`/families/${id}`, input);
    return unwrap(res.data);
  },

  async remove(id: string): Promise<void> {
    await apiClient.delete(`/families/${id}`);
  },

  async listCollaborators(familyId: string): Promise<Collaborator[]> {
    const res = await apiClient.get<ApiEnvelope<Collaborator[]>>(
      `/families/${familyId}/collaborators`,
    );
    return unwrap(res.data);
  },

  async inviteCollaborator(familyId: string, input: InviteCollaboratorInput): Promise<Collaborator> {
    const res = await apiClient.post<ApiEnvelope<Collaborator>>(
      `/families/${familyId}/collaborators`,
      input,
    );
    return unwrap(res.data);
  },

  async updateCollaborator(
    familyId: string,
    userId: string,
    input: UpdateCollaboratorInput,
  ): Promise<Collaborator> {
    const res = await apiClient.patch<ApiEnvelope<Collaborator>>(
      `/families/${familyId}/collaborators/${userId}`,
      input,
    );
    return unwrap(res.data);
  },

  async removeCollaborator(familyId: string, userId: string): Promise<void> {
    await apiClient.delete(`/families/${familyId}/collaborators/${userId}`);
  },
};
