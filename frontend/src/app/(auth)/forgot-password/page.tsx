"use client";

import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForgotPassword } from "@/hooks/use-auth";
import {
  forgotPasswordSchema,
  type ForgotPasswordFormValues,
} from "@/lib/validations/auth";
import { ApiError } from "@/types/api";

export default function ForgotPasswordPage() {
  const forgot = useForgotPassword();
  const form = useForm<ForgotPasswordFormValues>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: { email: "" },
  });

  const apiError =
    forgot.error instanceof ApiError && !forgot.isPending ? forgot.error.message : null;

  return (
    <Card className="border-0 bg-white/95 shadow-xl shadow-black/20 ring-teal-900/10">
      <CardTitle className="font-display text-2xl text-teal-950">Reset password</CardTitle>
      <CardDescription className="text-base">
        Enter your email. Until email delivery ships, ask an admin for a reset link, or use
        the token shown here in development.
      </CardDescription>

      <form
        className="mt-6 space-y-4"
        noValidate
        onSubmit={form.handleSubmit((values) => forgot.mutate(values.email))}
      >
        {apiError && <Alert tone="error">{apiError}</Alert>}

        {forgot.isSuccess && (
          <Alert tone="success">
            {forgot.data.message}
            {forgot.data.reset_token ? (
              <span className="mt-2 block break-all font-mono text-sm">
                Dev token: {forgot.data.reset_token}
              </span>
            ) : null}
          </Alert>
        )}

        <div>
          <Label htmlFor="email" className="text-base">
            Email
          </Label>
          <Input
            id="email"
            type="email"
            autoComplete="email"
            className="py-3 text-base"
            placeholder="you@example.com"
            error={form.formState.errors.email?.message}
            {...form.register("email")}
          />
        </div>

        <Button type="submit" size="lg" className="w-full py-3.5 text-base" loading={forgot.isPending}>
          Request reset
        </Button>
      </form>

      <p className="mt-6 text-center text-base text-slate-500">
        <Link href="/login" className="font-semibold text-teal-800 hover:text-teal-700">
          Back to sign in
        </Link>
        {forgot.data?.reset_token ? (
          <>
            {" · "}
            <Link
              href={`/reset-password?token=${encodeURIComponent(forgot.data.reset_token)}`}
              className="font-semibold text-teal-800 hover:text-teal-700"
            >
              Continue with token
            </Link>
          </>
        ) : null}
      </p>
    </Card>
  );
}
