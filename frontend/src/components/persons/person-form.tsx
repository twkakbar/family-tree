"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useCreatePerson } from "@/hooks/use-persons";
import { personSchema, type PersonFormValues } from "@/lib/validations/person";
import { ApiError } from "@/types/api";
import type { Family } from "@/types/family";

interface PersonFormProps {
  family: Family;
}

export function PersonForm({ family }: PersonFormProps) {
  const createPerson = useCreatePerson(family.id);
  const canEdit = family.role === "owner" || family.role === "editor";

  const form = useForm<PersonFormValues>({
    resolver: zodResolver(personSchema),
    defaultValues: {
      first_name: "",
      last_name: "",
      gender: "male",
      birth_date: "",
      death_date: "",
      bio: "",
    },
  });

  if (!canEdit) {
    return null;
  }

  const apiError =
    createPerson.error instanceof ApiError && !createPerson.isPending
      ? createPerson.error.message
      : null;

  return (
    <Card>
      <CardTitle>Add member</CardTitle>
      <CardDescription>Create a person record in this family tree.</CardDescription>

      <form
        className="mt-6 space-y-4"
        noValidate
        onSubmit={form.handleSubmit((values) =>
          createPerson.mutate({
            first_name: values.first_name,
            last_name: values.last_name || undefined,
            gender: values.gender,
            birth_date: values.birth_date || undefined,
            death_date: values.death_date || undefined,
            bio: values.bio || undefined,
          }, {
            onSuccess: () => form.reset({
              first_name: "",
              last_name: "",
              gender: "male",
              birth_date: "",
              death_date: "",
              bio: "",
            }),
          }),
        )}
      >
        {apiError && <Alert tone="error">{apiError}</Alert>}

        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <Label htmlFor="first-name">First name</Label>
            <Input id="first-name" error={form.formState.errors.first_name?.message} {...form.register("first_name")} />
          </div>
          <div>
            <Label htmlFor="last-name">Last name</Label>
            <Input id="last-name" error={form.formState.errors.last_name?.message} {...form.register("last_name")} />
          </div>
        </div>

        <div className="grid gap-4 sm:grid-cols-3">
          <div>
            <Label htmlFor="gender">Gender</Label>
            <select
              id="gender"
              className="block w-full rounded-lg border-0 px-3 py-2 text-sm text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 focus:ring-2 focus:ring-inset focus:ring-teal-700"
              {...form.register("gender")}
            >
              <option value="male">Male</option>
              <option value="female">Female</option>
            </select>
          </div>
          <div>
            <Label htmlFor="birth-date">Birth date</Label>
            <Input id="birth-date" type="date" error={form.formState.errors.birth_date?.message} {...form.register("birth_date")} />
          </div>
          <div>
            <Label htmlFor="death-date">Death date</Label>
            <Input id="death-date" type="date" error={form.formState.errors.death_date?.message} {...form.register("death_date")} />
          </div>
        </div>

        <div>
          <Label htmlFor="bio">Bio</Label>
          <Textarea id="bio" error={form.formState.errors.bio?.message} {...form.register("bio")} />
        </div>

        <Button type="submit" loading={createPerson.isPending}>
          Add member
        </Button>
      </form>
    </Card>
  );
}
