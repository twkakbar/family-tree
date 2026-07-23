import { apiClient, unwrap } from "@/lib/api-client";
import { ApiError, type ApiEnvelope } from "@/types/api";
import type {
  CreatePersonInput,
  CreateRelationshipInput,
  FamilyTreeData,
  Person,
  PersonDetail,
  PersonListResult,
  PersonWithDepth,
  Relationship,
  UpdatePersonInput,
} from "@/types/person";

export const personsService = {
  async listByFamily(
    familyId: string,
    search = "",
    page = 1,
    perPage = 50,
  ): Promise<PersonListResult> {
    const res = await apiClient.get<ApiEnvelope<Person[]>>(`/families/${familyId}/persons`, {
      params: { search: search || undefined, page, per_page: perPage },
    });
    const envelope = res.data;
    if (!envelope.success || envelope.data === undefined || !envelope.meta) {
      throw new ApiError(res.status, envelope.error);
    }
    return { items: envelope.data, meta: envelope.meta };
  },

  async getTree(familyId: string): Promise<FamilyTreeData> {
    const res = await apiClient.get<ApiEnvelope<FamilyTreeData>>(`/families/${familyId}/tree`);
    return unwrap(res.data);
  },

  async getById(personId: string): Promise<PersonDetail> {
    const res = await apiClient.get<ApiEnvelope<PersonDetail>>(`/persons/${personId}`);
    return unwrap(res.data);
  },

  async create(familyId: string, input: CreatePersonInput): Promise<Person> {
    const res = await apiClient.post<ApiEnvelope<Person>>(`/families/${familyId}/persons`, input);
    return unwrap(res.data);
  },

  async update(personId: string, input: UpdatePersonInput): Promise<Person> {
    const res = await apiClient.patch<ApiEnvelope<Person>>(`/persons/${personId}`, input);
    return unwrap(res.data);
  },

  async remove(personId: string): Promise<void> {
    await apiClient.delete(`/persons/${personId}`);
  },

  async uploadPhoto(personId: string, file: File): Promise<Person> {
    const body = new FormData();
    body.append("photo", file);
    const res = await apiClient.put<ApiEnvelope<Person>>(`/persons/${personId}/photo`, body, {
      headers: { "Content-Type": "multipart/form-data" },
      transformRequest: [
        (data, headers) => {
          // Let the browser set multipart boundary (default JSON Content-Type breaks uploads).
          if (data instanceof FormData && headers) {
            delete (headers as Record<string, unknown>)["Content-Type"];
          }
          return data;
        },
      ],
    });
    return unwrap(res.data);
  },

  async createRelationship(familyId: string, input: CreateRelationshipInput): Promise<Relationship> {
    const res = await apiClient.post<ApiEnvelope<Relationship>>(
      `/families/${familyId}/relationships`,
      input,
    );
    return unwrap(res.data);
  },

  async removeRelationship(relationshipId: string): Promise<void> {
    await apiClient.delete(`/relationships/${relationshipId}`);
  },

  async getAncestors(personId: string, depth = 5): Promise<PersonWithDepth[]> {
    const res = await apiClient.get<ApiEnvelope<PersonWithDepth[]>>(`/persons/${personId}/ancestors`, {
      params: { depth },
    });
    return unwrap(res.data);
  },

  async getDescendants(personId: string, depth = 5): Promise<PersonWithDepth[]> {
    const res = await apiClient.get<ApiEnvelope<PersonWithDepth[]>>(
      `/persons/${personId}/descendants`,
      { params: { depth } },
    );
    return unwrap(res.data);
  },
};
