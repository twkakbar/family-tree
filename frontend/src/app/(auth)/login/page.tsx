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
import { useLogin } from "@/hooks/use-auth";
import { loginSchema, type LoginFormValues } from "@/lib/validations/auth";
import { ApiError } from "@/types/api";

export default function LoginPage() {
  return (
    // useSearchParams requires a Suspense boundary during prerendering.
    <Suspense>
      <LoginForm />
    </Suspense>
  );
}

function LoginForm() {
  const searchParams = useSearchParams();
  const redirectTo = safeInternalPath(searchParams.get("next")) ?? "/dashboard";
  const login = useLogin(redirectTo);

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  const apiError =
    login.error instanceof ApiError && !login.isPending ? login.error.message : null;

  return (
    <Card>
      <CardTitle>Welcome back</CardTitle>
      <CardDescription>Sign in to continue to your family trees.</CardDescription>

      <form
        className="mt-6 space-y-4"
        noValidate
        onSubmit={form.handleSubmit((values) => login.mutate(values))}
      >
        {apiError && <Alert tone="error">{apiError}</Alert>}

        <div>
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            autoComplete="email"
            placeholder="you@example.com"
            error={form.formState.errors.email?.message}
            {...form.register("email")}
          />
        </div>

        <div>
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            autoComplete="current-password"
            placeholder="••••••••"
            error={form.formState.errors.password?.message}
            {...form.register("password")}
          />
        </div>

        <Button type="submit" className="w-full" loading={login.isPending}>
          Sign in
        </Button>
      </form>

      <p className="mt-6 text-center text-sm text-slate-500">
        No account yet?{" "}
        <Link href="/register" className="font-semibold text-emerald-600 hover:text-emerald-500">
          Create one
        </Link>
      </p>
    </Card>
  );
}

/** Only allow same-app relative redirects (no open-redirect). */
function safeInternalPath(value: string | null): string | null {
  if (value && value.startsWith("/") && !value.startsWith("//")) {
    return value;
  }
  return null;
}
