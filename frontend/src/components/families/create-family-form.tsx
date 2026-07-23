"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useCreateFamily } from "@/hooks/use-families";
import { createFamilySchema, type CreateFamilyFormValues } from "@/lib/validations/family";
import { ApiError } from "@/types/api";

/** Inline form for creating a new family from the dashboard. */
export function CreateFamilyForm() {
  const createFamily = useCreateFamily();
  const form = useForm<CreateFamilyFormValues>({
    resolver: zodResolver(createFamilySchema),
    defaultValues: { name: "", description: "" },
  });

  const apiError =
    createFamily.error instanceof ApiError && !createFamily.isPending
      ? createFamily.error.message
      : null;

  return (
    <Card>
      <CardTitle>Create a family</CardTitle>
      <CardDescription>Start a new tree and invite collaborators later.</CardDescription>

      <form
        className="mt-6 space-y-4"
        noValidate
        onSubmit={form.handleSubmit((values) =>
          createFamily.mutate({
            name: values.name,
            description: values.description || undefined,
          }),
        )}
      >
        {apiError && <Alert tone="error">{apiError}</Alert>}

        <div>
          <Label htmlFor="family-name">Family name</Label>
          <Input
            id="family-name"
            placeholder="The Does"
            error={form.formState.errors.name?.message}
            {...form.register("name")}
          />
        </div>

        <div>
          <Label htmlFor="family-description">Description (optional)</Label>
          <Textarea
            id="family-description"
            placeholder="A short note about this branch of the tree"
            error={form.formState.errors.description?.message}
            {...form.register("description")}
          />
        </div>

        <Button type="submit" loading={createFamily.isPending}>
          Create family
        </Button>
      </form>
    </Card>
  );
}
