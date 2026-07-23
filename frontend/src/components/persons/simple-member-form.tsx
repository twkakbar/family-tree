"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { PersonSearchPicker } from "@/components/persons/person-search-picker";
import { useCreatePerson, useCreateRelationship, useFamilyTree } from "@/hooks/use-persons";
import { personSchema } from "@/lib/validations/person";
import { ApiError } from "@/types/api";
import type { Family } from "@/types/family";
import type { Person } from "@/types/person";
import { displayName } from "@/utils/tree-layout";
import { cn } from "@/utils/cn";

const formSchema = personSchema
  .omit({ gender: true })
  .extend({
    gender: z.enum(["male", "female"], { message: "Choose male or female" }),
    link_kind: z.enum(["none", "child_of", "parent_of", "spouse_of"]),
    link_person_id: z.string().optional().or(z.literal("")),
  })
  .superRefine((data, ctx) => {
    if (data.link_kind !== "none" && !data.link_person_id) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Please choose a person from the list",
        path: ["link_person_id"],
      });
    }
  });

type FormValues = z.infer<typeof formSchema>;

interface SimpleMemberFormProps {
  family: Family;
  onCreated?: (person: Person) => void;
}

const genderOptions = [
  { value: "male", label: "Male" },
  { value: "female", label: "Female" },
] as const;

const linkOptions = [
  { value: "none", label: "No link yet — add later" },
  { value: "child_of", label: "Child of someone already in the tree" },
  { value: "parent_of", label: "Parent of someone already in the tree" },
  { value: "spouse_of", label: "Spouse / partner of someone in the tree" },
] as const;

/**
 * Large-type, step-by-step form for adding a person and optionally linking them.
 * Tuned for elderly editors: big controls, plain language, one job at a time.
 */
export function SimpleMemberForm({ family, onCreated }: SimpleMemberFormProps) {
  const tree = useFamilyTree(family.id);
  const createPerson = useCreatePerson(family.id);
  const createRelationship = useCreateRelationship(family.id);
  const [successName, setSuccessName] = useState<string | null>(null);

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      first_name: "",
      last_name: "",
      gender: "male",
      birth_date: "",
      death_date: "",
      bio: "",
      link_kind: "none",
      link_person_id: "",
    },
  });

  const linkKind = form.watch("link_kind");
  const linkPersonId = form.watch("link_person_id") ?? "";
  const people = tree.data?.persons ?? [];

  const apiError =
    (createPerson.error instanceof ApiError && !createPerson.isPending
      ? createPerson.error.message
      : null) ??
    (createRelationship.error instanceof ApiError && !createRelationship.isPending
      ? createRelationship.error.message
      : null);

  const resetForm = () => {
    form.reset({
      first_name: "",
      last_name: "",
      gender: "male",
      birth_date: "",
      death_date: "",
      bio: "",
      link_kind: "none",
      link_person_id: "",
    });
  };

  const onSubmit = form.handleSubmit(async (values) => {
    setSuccessName(null);
    try {
      const person = await createPerson.mutateAsync({
        first_name: values.first_name,
        gender: values.gender,
        birth_date: values.birth_date || undefined,
        death_date: values.death_date || undefined,
        bio: values.bio || undefined,
      });

      if (values.link_kind !== "none" && values.link_person_id) {
        if (values.link_kind === "child_of") {
          await createRelationship.mutateAsync({
            from_person_id: values.link_person_id,
            to_person_id: person.id,
            type: "parent",
          });
        } else if (values.link_kind === "parent_of") {
          await createRelationship.mutateAsync({
            from_person_id: person.id,
            to_person_id: values.link_person_id,
            type: "parent",
          });
        } else if (values.link_kind === "spouse_of") {
          await createRelationship.mutateAsync({
            from_person_id: person.id,
            to_person_id: values.link_person_id,
            type: "spouse",
          });
        }
      }

      setSuccessName(displayName(person));
      onCreated?.(person);
      resetForm();
    } catch {
      // Errors surface via mutation state / apiError.
    }
  });

  const busy = createPerson.isPending || createRelationship.isPending;

  return (
    <div className="aceh-panel">
      <p className="text-sm font-semibold uppercase tracking-wider text-teal-800/70">Step 1</p>
      <h2 className="font-display mt-1 text-3xl text-teal-950 sm:text-4xl">Add a family member</h2>
      <p className="mt-2 text-lg text-slate-600">
        Fill in the name, then (if you know) connect them to someone already on the tree.
      </p>

      <form className="mt-8 space-y-7" noValidate onSubmit={onSubmit}>
        {apiError && <Alert tone="error">{apiError}</Alert>}
        {successName && (
          <Alert tone="success">
            Saved <strong>{successName}</strong>. You can add another person below.
          </Alert>
        )}

        <div>
          <Label htmlFor="first-name" className="text-base text-slate-800">
            Name <span className="text-teal-700">*</span>
          </Label>
          <Input
            id="first-name"
            className="mt-1 py-3.5 text-lg"
            autoComplete="name"
            placeholder="e.g. Cut Nyak Dhien"
            error={form.formState.errors.first_name?.message}
            {...form.register("first_name")}
          />
        </div>

        <fieldset>
          <legend className="mb-3 text-base font-medium text-slate-800">Gender</legend>
          <div className="grid grid-cols-2 gap-3">
            {genderOptions.map((opt) => {
              const selected = form.watch("gender") === opt.value;
              return (
                <button
                  key={opt.value}
                  type="button"
                  className={cn(
                    "rounded-xl px-3 py-3.5 text-base font-semibold ring-1 transition",
                    selected
                      ? "bg-teal-800 text-white ring-teal-900"
                      : "bg-white text-slate-800 ring-slate-300 hover:ring-teal-500",
                  )}
                  onClick={() => form.setValue("gender", opt.value, { shouldValidate: true })}
                >
                  {opt.label}
                </button>
              );
            })}
          </div>
        </fieldset>

        <div className="grid gap-6 sm:grid-cols-2">
          <div>
            <Label htmlFor="birth-date" className="text-base text-slate-800">
              Birth date (optional)
            </Label>
            <Input
              id="birth-date"
              type="date"
              className="mt-1 py-3.5 text-lg"
              error={form.formState.errors.birth_date?.message}
              {...form.register("birth_date")}
            />
          </div>
          <div>
            <Label htmlFor="death-date" className="text-base text-slate-800">
              Death date (optional)
            </Label>
            <Input
              id="death-date"
              type="date"
              className="mt-1 py-3.5 text-lg"
              error={form.formState.errors.death_date?.message}
              {...form.register("death_date")}
            />
          </div>
        </div>

        <div>
          <Label htmlFor="bio" className="text-base text-slate-800">
            Short note (optional)
          </Label>
          <Textarea
            id="bio"
            rows={3}
            className="mt-1 text-lg"
            placeholder="Anything helpful to remember this person"
            error={form.formState.errors.bio?.message}
            {...form.register("bio")}
          />
        </div>

        <div className="rounded-2xl bg-teal-950/5 p-5 ring-1 ring-teal-900/10">
          <p className="text-sm font-semibold uppercase tracking-wider text-teal-800/70">Step 2</p>
          <Label htmlFor="link-kind" className="mt-2 text-lg text-slate-900">
            How are they related? (optional)
          </Label>
          <select
            id="link-kind"
            className="mt-2 block w-full rounded-xl border-0 bg-white px-4 py-3.5 text-lg text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 focus:ring-2 focus:ring-teal-700"
            {...form.register("link_kind")}
          >
            {linkOptions.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>

          {linkKind !== "none" && (
            <div className="mt-4">
              <Label htmlFor="person-search" className="text-base text-slate-800">
                Find the person already on the tree
              </Label>
              {people.length === 0 ? (
                <p className="mt-2 text-base text-slate-600">
                  No one is on the tree yet — save this person first, then add their relatives.
                </p>
              ) : (
                <PersonSearchPicker
                  id="person-search"
                  people={people}
                  value={linkPersonId}
                  onChange={(id) =>
                    form.setValue("link_person_id", id, { shouldValidate: true, shouldDirty: true })
                  }
                  error={form.formState.errors.link_person_id?.message}
                />
              )}
            </div>
          )}
        </div>

        <Button type="submit" size="lg" className="w-full py-4 text-lg sm:w-auto sm:min-w-56" loading={busy}>
          Save this person
        </Button>
      </form>
    </div>
  );
}
