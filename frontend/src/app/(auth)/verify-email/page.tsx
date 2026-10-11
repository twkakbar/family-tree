"use client";

import { Suspense } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { useVerifyEmail } from "@/hooks/use-auth";
import { ApiError } from "@/types/api";

export default function VerifyEmailPage() {
  return (
    <Suspense>
      <VerifyEmailForm />
    </Suspense>
  );
}

function VerifyEmailForm() {
  const searchParams = useSearchParams();
  const token = searchParams.get("token") ?? "";
  const verify = useVerifyEmail();

  const apiError =
    verify.error instanceof ApiError && !verify.isPending ? verify.error.message : null;

  return (
    <Card className="border-0 bg-white/95 shadow-xl shadow-black/20 ring-teal-900/10">
      <CardTitle className="font-display text-2xl text-teal-950">Verify your email</CardTitle>
      <CardDescription className="text-base">
        Confirm this address so you can sign in.
      </CardDescription>

      <div className="mt-6 space-y-4">
        {apiError && <Alert tone="error">{apiError}</Alert>}
        {verify.isSuccess && (
          <Alert tone="success">Your email is verified. You can sign in now.</Alert>
        )}
        {!token && <Alert tone="error">This verification link is missing its token.</Alert>}

        {!verify.isSuccess && (
          <Button
            type="button"
            size="lg"
            className="w-full py-3.5 text-base"
            loading={verify.isPending}
            disabled={!token}
            onClick={() => verify.mutate(token)}
          >
            Verify email
          </Button>
        )}
      </div>

      <p className="mt-6 text-center text-base text-slate-500">
        <Link href="/login" className="font-semibold text-teal-800 hover:text-teal-700">
          Back to sign in
        </Link>
      </p>
    </Card>
  );
}
