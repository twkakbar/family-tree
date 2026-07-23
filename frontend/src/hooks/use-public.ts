"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { queryKeys } from "@/lib/query-keys";
import { accessService } from "@/services/access.service";
import { publicService } from "@/services/public.service";

export function useRoyalMeta() {
  return useQuery({
    queryKey: queryKeys.public.meta,
    queryFn: publicService.getMeta,
  });
}

export function usePublicTree() {
  return useQuery({
    queryKey: queryKeys.public.tree,
    queryFn: publicService.getTree,
  });
}

export function usePublicPerson(id: string | null) {
  return useQuery({
    queryKey: queryKeys.public.person(id ?? ""),
    queryFn: () => publicService.getPerson(id!),
    enabled: Boolean(id),
  });
}

export function useStudioCanEdit(enabled = true) {
  return useQuery({
    queryKey: queryKeys.studio.canEdit,
    queryFn: accessService.canEdit,
    enabled,
    retry: false,
  });
}

export function useMyAccessRequest(enabled = true) {
  return useQuery({
    queryKey: queryKeys.access.me,
    queryFn: accessService.myLatest,
    enabled,
    retry: false,
  });
}

export function useSubmitAccessRequest() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (message: string) => accessService.submit(message),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.access.me });
    },
  });
}

export function usePendingAccessRequests(enabled = true) {
  return useQuery({
    queryKey: queryKeys.admin.accessRequests,
    queryFn: accessService.listPending,
    enabled,
  });
}

export function useReviewAccessRequest() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, action }: { id: string; action: "approve" | "reject" }) =>
      action === "approve" ? accessService.approve(id) : accessService.reject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.admin.accessRequests });
      queryClient.invalidateQueries({ queryKey: queryKeys.studio.canEdit });
      queryClient.invalidateQueries({ queryKey: queryKeys.access.me });
    },
  });
}
