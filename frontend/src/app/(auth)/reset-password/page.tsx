"use client";

import { Suspense } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useResetPassword } from "@/hooks/use-auth";
import {
  resetPasswordSchema,
  type ResetPasswordFormValues,
} from "@/lib/validations/auth";
import { ApiError } from "@/types/api";

export default function ResetPasswordPage() {
  return (
    <Suspense>
      <ResetPasswordForm />
    </Suspense>
  );
}

function ResetPasswordForm() {
  const searchParams = useSearchParams();
  const tokenFromUrl = searchParams.get("token") ?? "";
  const reset = useResetPassword();

  const form = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      token: tokenFromUrl,
      password: "",
      confirmPassword: "",
    },
  });

  const apiError =
    reset.error instanceof ApiError && !reset.isPending ? reset.error.message : null;

  return (
    <Card className="border-0 bg-white/95 shadow-xl shadow-black/20 ring-teal-900/10">
      <CardTitle className="font-display text-2xl text-teal-950">Choose a new password</CardTitle>
      <CardDescription className="text-base">
        Paste the reset token from an admin (or from the forgot-password response in
        development).
      </CardDescription>

      <form
        className="mt-6 space-y-4"
        noValidate
        onSubmit={form.handleSubmit((values) =>
          reset.mutate({ token: values.token, password: values.password }),
        )}
      >
        {apiError && <Alert tone="error">{apiError}</Alert>}

        <div>
          <Label htmlFor="token" className="text-base">
            Reset token
          </Label>
          <Input
            id="token"
            className="py-3 font-mono text-sm"
            placeholder="Paste token"
            error={form.formState.errors.token?.message}
            {...form.register("token")}
          />
        </div>

        <div>
          <Label htmlFor="password" className="text-base">
            New password
          </Label>
          <Input
            id="password"
            type="password"
            autoComplete="new-password"
            className="py-3 text-base"
            placeholder="••••••••"
            error={form.formState.errors.password?.message}
            {...form.register("password")}
          />
        </div>

        <div>
          <Label htmlFor="confirmPassword" className="text-base">
            Confirm password
          </Label>
          <Input
            id="confirmPassword"
            type="password"
            autoComplete="new-password"
            className="py-3 text-base"
            placeholder="••••••••"
            error={form.formState.errors.confirmPassword?.message}
            {...form.register("confirmPassword")}
          />
        </div>

        <Button type="submit" size="lg" className="w-full py-3.5 text-base" loading={reset.isPending}>
          Save new password
        </Button>
      </form>

      <p className="mt-6 text-center text-base text-slate-500">
        <Link href="/login" className="font-semibold text-teal-800 hover:text-teal-700">
          Back to sign in
        </Link>
      </p>
    </Card>
  );
}
