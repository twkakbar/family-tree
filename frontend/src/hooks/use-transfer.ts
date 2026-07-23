"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { queryKeys } from "@/lib/query-keys";
import { transferService, type TreeTransfer } from "@/services/transfer.service";

export function useExportTree(familyId: string) {
  return useMutation({
    mutationFn: () => transferService.exportJSON(familyId),
  });
}

export function useImportTree(familyId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (doc: TreeTransfer) => transferService.importJSON(familyId, doc),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.families.tree(familyId) });
      queryClient.invalidateQueries({ queryKey: ["families", "persons", familyId] });
      queryClient.invalidateQueries({ queryKey: queryKeys.public.tree });
      queryClient.invalidateQueries({ queryKey: queryKeys.public.meta });
    },
  });
}
