"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

import { queryKeys } from "@/lib/query-keys";
import { familiesService } from "@/services/families.service";
import type {
  CreateFamilyInput,
  InviteCollaboratorInput,
  UpdateCollaboratorInput,
  UpdateFamilyInput,
} from "@/types/family";

const defaultPageSize = 20;

export function useFamilies(page = 1, perPage = defaultPageSize) {
  return useQuery({
    queryKey: queryKeys.families.list(page, perPage),
    queryFn: () => familiesService.list(page, perPage),
  });
}

export function useFamily(id: string) {
  return useQuery({
    queryKey: queryKeys.families.detail(id),
    queryFn: () => familiesService.getById(id),
    enabled: Boolean(id),
  });
}

export function useFamilyCollaborators(familyId: string) {
  return useQuery({
    queryKey: queryKeys.families.collaborators(familyId),
    queryFn: () => familiesService.listCollaborators(familyId),
    enabled: Boolean(familyId),
  });
}

export function useCreateFamily() {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: (input: CreateFamilyInput) => familiesService.create(input),
    onSuccess: (family) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.families.all });
      router.push(`/families/${family.id}`);
    },
  });
}

export function useUpdateFamily(familyId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateFamilyInput) => familiesService.update(familyId, input),
    onSuccess: (family) => {
      queryClient.setQueryData(queryKeys.families.detail(familyId), family);
      queryClient.invalidateQueries({ queryKey: queryKeys.families.all });
    },
  });
}

export function useDeleteFamily() {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: (familyId: string) => familiesService.remove(familyId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.families.all });
      router.push("/dashboard");
    },
  });
}

export function useInviteCollaborator(familyId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: InviteCollaboratorInput) =>
      familiesService.inviteCollaborator(familyId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.families.collaborators(familyId) });
    },
  });
}

export function useUpdateCollaborator(familyId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userId, input }: { userId: string; input: UpdateCollaboratorInput }) =>
      familiesService.updateCollaborator(familyId, userId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.families.collaborators(familyId) });
    },
  });
}

export function useRemoveCollaborator(familyId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userId: string) => familiesService.removeCollaborator(familyId, userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.families.collaborators(familyId) });
    },
  });
}
