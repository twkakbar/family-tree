"use client";

import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useMe } from "@/hooks/use-auth";
import { useMyAccessRequest, useStudioCanEdit, useSubmitAccessRequest } from "@/hooks/use-public";
import { ApiError } from "@/types/api";

const schema = z.object({
  message: z
    .string()
    .trim()
    .min(10, "Tell us a bit more (at least 10 characters)")
    .max(2000),
});

type FormValues = z.infer<typeof schema>;

export default function RequestAccessPage() {
  const { data: user } = useMe();
  const canEdit = useStudioCanEdit();
  const myRequest = useMyAccessRequest();
  const submit = useSubmitAccessRequest();

  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { message: "" },
  });

  if (canEdit.data?.can_edit) {
    return (
      <div className="aceh-panel mx-auto max-w-lg">
        <h1 className="font-display text-3xl text-teal-950">You can already edit</h1>
        <p className="mt-2 text-lg text-slate-600">Open the studio to add members.</p>
        <Link
          href="/studio"
          className="mt-6 inline-flex rounded-xl bg-teal-800 px-5 py-3.5 text-lg font-semibold text-white hover:bg-teal-900"
        >
          Go to studio
        </Link>
      </div>
    );
  }

  const latest = myRequest.isError ? null : myRequest.data;
  const apiError =
    submit.error instanceof ApiError && !submit.isPending ? submit.error.message : null;

  return (
    <div className="mx-auto max-w-lg space-y-6">
      <div>
        <h1 className="font-display text-4xl text-teal-950">Request edit access</h1>
        <p className="mt-3 text-lg text-slate-600">
          Hello{user ? `, ${user.name}` : ""}. Write a short note about how you connect to the Aceh
          royal line.
        </p>
      </div>

      {latest && (
        <Alert tone={latest.status === "rejected" ? "error" : "info"}>
          Your latest request is <strong>{latest.status}</strong>
          {latest.status === "pending" ? " — please wait for an administrator." : "."}
        </Alert>
      )}

      {(!latest || latest.status === "rejected") && (
        <div className="aceh-panel">
          <form
            className="space-y-5"
            noValidate
            onSubmit={form.handleSubmit((values) => submit.mutate(values.message))}
          >
            {apiError && <Alert tone="error">{apiError}</Alert>}
            <div>
              <Label htmlFor="message" className="text-base">
                Your message
              </Label>
              <Textarea
                id="message"
                rows={5}
                className="mt-1 text-lg"
                error={form.formState.errors.message?.message}
                {...form.register("message")}
              />
            </div>
            <Button type="submit" size="lg" className="w-full py-4 text-lg" loading={submit.isPending}>
              Send request
            </Button>
          </form>
        </div>
      )}

      <p className="text-base text-slate-500">
        <Link href="/" className="font-medium text-teal-800 hover:text-teal-700">
          ← Back to the tree
        </Link>
      </p>
    </div>
  );
}
