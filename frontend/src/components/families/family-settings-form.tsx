"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useDeleteFamily, useUpdateFamily } from "@/hooks/use-families";
import { updateFamilySchema, type UpdateFamilyFormValues } from "@/lib/validations/family";
import { ApiError } from "@/types/api";
import type { Family } from "@/types/family";

interface FamilySettingsFormProps {
  family: Family;
}

/** Owner-only family rename/description and delete controls. */
export function FamilySettingsForm({ family }: FamilySettingsFormProps) {
  const updateFamily = useUpdateFamily(family.id);
  const deleteFamily = useDeleteFamily();

  const form = useForm<UpdateFamilyFormValues>({
    resolver: zodResolver(updateFamilySchema),
    defaultValues: { name: family.name, description: family.description },
  });

  const apiError =
    (updateFamily.error instanceof ApiError && !updateFamily.isPending
      ? updateFamily.error.message
      : null) ??
    (deleteFamily.error instanceof ApiError && !deleteFamily.isPending
      ? deleteFamily.error.message
      : null);

  if (family.role !== "owner") {
    return null;
  }

  return (
    <Card>
      <CardTitle>Family settings</CardTitle>
      <CardDescription>Update details or delete this family.</CardDescription>

      <form
        className="mt-6 space-y-4"
        noValidate
        onSubmit={form.handleSubmit((values) =>
          updateFamily.mutate({
            name: values.name,
            description: values.description || undefined,
          }),
        )}
      >
        {apiError && <Alert tone="error">{apiError}</Alert>}

        <div>
          <Label htmlFor="edit-family-name">Family name</Label>
          <Input
            id="edit-family-name"
            error={form.formState.errors.name?.message}
            {...form.register("name")}
          />
        </div>

        <div>
          <Label htmlFor="edit-family-description">Description</Label>
          <Textarea
            id="edit-family-description"
            error={form.formState.errors.description?.message}
            {...form.register("description")}
          />
        </div>

        <div className="flex flex-wrap gap-3">
          <Button type="submit" loading={updateFamily.isPending}>
            Save changes
          </Button>
          <Button
            type="button"
            variant="danger"
            loading={deleteFamily.isPending}
            onClick={() => {
              if (window.confirm("Delete this family and all of its data? This cannot be undone.")) {
                deleteFamily.mutate(family.id);
              }
            }}
          >
            Delete family
          </Button>
        </div>
      </form>
    </Card>
  );
}
