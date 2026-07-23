import { apiClient, unwrap } from "@/lib/api-client";
import type { ApiEnvelope } from "@/types/api";
import type { FamilyTreeData, PersonDetail, PersonWithDepth } from "@/types/person";

export interface RoyalMeta {
  id: string;
  name: string;
  description: string;
}

export const publicService = {
  async getMeta(): Promise<RoyalMeta> {
    const res = await apiClient.get<ApiEnvelope<RoyalMeta>>("/public/meta");
    return unwrap(res.data);
  },

  async getTree(): Promise<FamilyTreeData> {
    const res = await apiClient.get<ApiEnvelope<FamilyTreeData>>("/public/tree");
    return unwrap(res.data);
  },

  async getPerson(id: string): Promise<PersonDetail> {
    const res = await apiClient.get<ApiEnvelope<PersonDetail>>(`/public/persons/${id}`);
    return unwrap(res.data);
  },

  async getAncestors(id: string, depth = 5): Promise<PersonWithDepth[]> {
    const res = await apiClient.get<ApiEnvelope<PersonWithDepth[]>>(
      `/public/persons/${id}/ancestors`,
      { params: { depth } },
    );
    return unwrap(res.data);
  },

  async getDescendants(id: string, depth = 5): Promise<PersonWithDepth[]> {
    const res = await apiClient.get<ApiEnvelope<PersonWithDepth[]>>(
      `/public/persons/${id}/descendants`,
      { params: { depth } },
    );
    return unwrap(res.data);
  },
};
