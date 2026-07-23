"use client";

import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useRegister } from "@/hooks/use-auth";
import { registerSchema, type RegisterFormValues } from "@/lib/validations/auth";
import { ApiError } from "@/types/api";

export default function RegisterPage() {
  const register = useRegister();

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { name: "", email: "", password: "", confirmPassword: "" },
  });

  const onSubmit = (values: RegisterFormValues) => {
    register.mutate(
      { name: values.name, email: values.email, password: values.password },
      {
        onError: (error) => {
          // Merge server-side 422 field errors into the form.
          if (error instanceof ApiError && error.details) {
            for (const [field, message] of Object.entries(error.details)) {
              if (field === "name" || field === "email" || field === "password") {
                form.setError(field, { type: "server", message });
              }
            }
          }
        },
      },
    );
  };

  const apiError =
    register.error instanceof ApiError && !register.isPending && !register.error.details
      ? register.error.message
      : null;

  return (
    <Card className="border-0 bg-white/95 shadow-xl shadow-black/20 ring-teal-900/10">
      <CardTitle className="font-display text-2xl text-teal-950">Create your account</CardTitle>
      <CardDescription className="text-base">
        Register so an administrator can grant tree edit access.
      </CardDescription>

      <form className="mt-6 space-y-4" noValidate onSubmit={form.handleSubmit(onSubmit)}>
        {apiError && <Alert tone="error">{apiError}</Alert>}

        <div>
          <Label htmlFor="name" className="text-base">Full name</Label>
          <Input
            id="name"
            autoComplete="name"
            className="py-3 text-base"
            placeholder="Your name"
            error={form.formState.errors.name?.message}
            {...form.register("name")}
          />
        </div>

        <div>
          <Label htmlFor="email" className="text-base">Email</Label>
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

        <div>
          <Label htmlFor="password" className="text-base">Password</Label>
          <Input
            id="password"
            type="password"
            autoComplete="new-password"
            className="py-3 text-base"
            placeholder="At least 8 characters, letters and digits"
            error={form.formState.errors.password?.message}
            {...form.register("password")}
          />
        </div>

        <div>
          <Label htmlFor="confirmPassword" className="text-base">Confirm password</Label>
          <Input
            id="confirmPassword"
            type="password"
            autoComplete="new-password"
            className="py-3 text-base"
            placeholder="Repeat your password"
            error={form.formState.errors.confirmPassword?.message}
            {...form.register("confirmPassword")}
          />
        </div>

        <Button type="submit" size="lg" className="w-full py-3.5 text-base" loading={register.isPending}>
          Create account
        </Button>
      </form>

      <p className="mt-6 text-center text-base text-slate-500">
        Already have an account?{" "}
        <Link href="/login" className="font-semibold text-teal-800 hover:text-teal-700">
          Sign in
        </Link>
      </p>
    </Card>
  );
}
