"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { queryKeys } from "@/lib/query-keys";
import {
  adminUsersService,
  type AdminUpdateUserInput,
} from "@/services/admin-users.service";

export function useAdminUsers(search: string, page: number, enabled = true) {
  return useQuery({
    queryKey: queryKeys.admin.users(search, page),
    queryFn: () => adminUsersService.list(search, page),
    enabled,
  });
}

export function useUpdateAdminUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: AdminUpdateUserInput }) =>
      adminUsersService.update(id, input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
  });
}

export function useDeleteAdminUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminUsersService.remove(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
  });
}

export function useIssuePasswordReset() {
  return useMutation({
    mutationFn: (id: string) => adminUsersService.issuePasswordReset(id),
  });
}
