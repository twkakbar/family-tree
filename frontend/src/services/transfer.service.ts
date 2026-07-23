import { apiClient, unwrap } from "@/lib/api-client";
import type { ApiEnvelope } from "@/types/api";

export interface TreeTransferFamily {
  name: string;
  description: string;
}

export interface TreeTransferPerson {
  id: string;
  first_name: string;
  last_name: string;
  gender: string;
  birth_date?: string | null;
  death_date?: string | null;
  bio: string;
}

export interface TreeTransferEdge {
  from_person_id: string;
  to_person_id: string;
  type: "parent" | "spouse";
}

export interface TreeTransfer {
  format_version: number;
  exported_at: string;
  family: TreeTransferFamily;
  persons: TreeTransferPerson[];
  relationships: TreeTransferEdge[];
}

export interface ImportResult {
  persons_created: number;
  relationships_created: number;
}

export const transferService = {
  async exportJSON(familyId: string): Promise<TreeTransfer> {
    const res = await apiClient.get<ApiEnvelope<TreeTransfer>>(`/families/${familyId}/export`, {
      params: { format: "json" },
    });
    return unwrap(res.data);
  },

  async importJSON(familyId: string, doc: TreeTransfer): Promise<ImportResult> {
    const res = await apiClient.post<ApiEnvelope<ImportResult>>(`/families/${familyId}/import`, doc);
    return unwrap(res.data);
  },
};
