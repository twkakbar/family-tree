"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

import { authToken } from "@/lib/auth-token";
import { queryKeys } from "@/lib/query-keys";
import { authService } from "@/services/auth.service";
import { usersService } from "@/services/users.service";
import type { AuthResult } from "@/types/auth";

/**
 * Auth hooks: the only way components interact with authentication.
 * Session establishment = store access token in memory + prime the me-query.
 */

export function useMe() {
  return useQuery({
    queryKey: queryKeys.auth.me,
    queryFn: usersService.getMe,
    retry: false,
  });
}

export function useLogin(redirectTo: string = "/dashboard") {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: authService.login,
    onSuccess: (result: AuthResult) => {
      establishSession(queryClient, result);
      router.push(redirectTo);
    },
  });
}

export function useRegister(redirectTo: string = "/dashboard") {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: authService.register,
    onSuccess: (result: AuthResult) => {
      establishSession(queryClient, result);
      router.push(redirectTo);
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: authService.logout,
    onSettled: () => {
      // Even if the API call failed, drop all local session state.
      authToken.clear();
      queryClient.clear();
      router.push("/login");
    },
  });
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: usersService.updateMe,
    onSuccess: (user) => {
      queryClient.setQueryData(queryKeys.auth.me, user);
    },
  });
}

export function useForgotPassword() {
  return useMutation({
    mutationFn: (email: string) => authService.forgotPassword(email),
  });
}

export function useResetPassword() {
  const router = useRouter();

  return useMutation({
    mutationFn: ({ token, password }: { token: string; password: string }) =>
      authService.resetPassword(token, password),
    onSuccess: () => {
      router.push("/login");
    },
  });
}

function establishSession(
  queryClient: ReturnType<typeof useQueryClient>,
  result: AuthResult,
): void {
  authToken.set(result.access_token);
  queryClient.setQueryData(queryKeys.auth.me, result.user);
}
