"use client";

import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import {
  useDeletePerson,
  useDeleteRelationship,
  usePerson,
  usePersonAncestors,
  usePersonDescendants,
  useUpdatePerson,
  useUploadPersonPhoto,
} from "@/hooks/use-persons";
import { ApiError } from "@/types/api";
import type { Family } from "@/types/family";
import type { Person, PersonWithDepth } from "@/types/person";
import { displayName } from "@/utils/tree-layout";

const editSchema = z.object({
  first_name: z.string().trim().min(1, "Name is required").max(100),
  gender: z.enum(["male", "female"], { message: "Choose male or female" }),
  birth_date: z
    .string()
    .trim()
    .regex(/^\d{4}-\d{2}-\d{2}$/, "Use YYYY-MM-DD")
    .optional()
    .or(z.literal("")),
  death_date: z
    .string()
    .trim()
    .regex(/^\d{4}-\d{2}-\d{2}$/, "Use YYYY-MM-DD")
    .optional()
    .or(z.literal("")),
  bio: z.string().trim().max(5000).optional().or(z.literal("")),
});

type EditValues = z.infer<typeof editSchema>;

interface PersonDetailPanelProps {
  family: Family;
  person: Person;
  onClose: () => void;
  onUpdated?: (person: Person) => void;
  /** Switch selection from ancestor/descendant links inside the panel. */
  onSelectRelated?: (person: Person) => void;
}

export function PersonDetailPanel({
  family,
  person,
  onClose,
  onUpdated,
  onSelectRelated,
}: PersonDetailPanelProps) {
  const personId = person.id;
  const detail = usePerson(personId);
  const ancestors = usePersonAncestors(personId);
  const descendants = usePersonDescendants(personId);
  const updatePerson = useUpdatePerson(family.id, personId);
  const deletePerson = useDeletePerson(family.id);
  const deleteRelationship = useDeleteRelationship(family.id);
  const uploadPhoto = useUploadPersonPhoto(family.id);
  const fileRef = useRef<HTMLInputElement>(null);
  const [saveNotice, setSaveNotice] = useState<"saved" | "invalid" | null>(null);

  const canEdit = family.role === "owner" || family.role === "editor";
  const shown = detail.data ?? person;

  const form = useForm<EditValues>({
    resolver: zodResolver(editSchema),
    defaultValues: valuesFromPerson(shown),
  });

  useEffect(() => {
    setSaveNotice(null);
  }, [person.id]);

  useEffect(() => {
    form.reset(valuesFromPerson(detail.data ?? person));
    // Re-seed when switching members or when detail finishes loading.
    // eslint-disable-next-line react-hooks/exhaustive-deps -- form.reset is stable enough here
  }, [person.id, detail.data]);

  const apiError =
    (updatePerson.error instanceof ApiError && !updatePerson.isPending
      ? updatePerson.error.message
      : null) ??
    (deletePerson.error instanceof ApiError && !deletePerson.isPending
      ? deletePerson.error.message
      : null) ??
    (deleteRelationship.error instanceof ApiError && !deleteRelationship.isPending
      ? deleteRelationship.error.message
      : null) ??
    (uploadPhoto.error instanceof ApiError && !uploadPhoto.isPending
      ? uploadPhoto.error.message
      : null);

  return (
    <div>
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-start gap-4">
          {shown.photo_url ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={shown.photo_url}
              alt={displayName(shown)}
              className="size-20 rounded-2xl object-cover object-top ring-1 ring-teal-900/15"
            />
          ) : (
            <div className="grid size-20 place-items-center rounded-2xl bg-teal-100 font-display text-2xl text-teal-900">
              {displayName(shown).slice(0, 1).toUpperCase()}
            </div>
          )}
          <div>
            <h2 className="font-display text-2xl text-teal-950">{displayName(shown)}</h2>
            <p className="mt-1 capitalize text-slate-600">{shown.gender}</p>
          </div>
        </div>
        <Button variant="ghost" size="sm" onClick={onClose}>
          Close
        </Button>
      </div>

      {shown.photo_url && (
        <div className="mt-4 overflow-hidden rounded-2xl bg-teal-950/5 ring-1 ring-teal-900/10">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={shown.photo_url}
            alt={displayName(shown)}
            className="mx-auto max-h-80 w-full object-contain"
          />
        </div>
      )}

      {apiError && (
        <Alert tone="error" className="mt-4">
          {apiError}
        </Alert>
      )}

      {canEdit && (
        <div className="mt-5">
          <input
            ref={fileRef}
            type="file"
            accept="image/jpeg,image/png,image/webp"
            className="sr-only"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (!file) return;
              uploadPhoto.mutate(
                { personId: shown.id, file },
                {
                  onSuccess: (updated) => onUpdated?.(updated),
                },
              );
              e.target.value = "";
            }}
          />
          <Button
            type="button"
            variant="secondary"
            size="lg"
            className="text-base"
            loading={uploadPhoto.isPending}
            onClick={() => fileRef.current?.click()}
          >
            {shown.photo_url ? "Change photo" : "Add photo"}
          </Button>
          <p className="mt-2 text-sm text-slate-500">JPEG, PNG, or WebP.</p>
        </div>
      )}

      {canEdit && (
        <form
          className="mt-6 space-y-4 border-t border-teal-900/10 pt-5"
          noValidate
          onSubmit={form.handleSubmit(
            (values) => {
              setSaveNotice(null);
              updatePerson.mutate(
                {
                  first_name: values.first_name,
                  gender: values.gender,
                  birth_date: values.birth_date || undefined,
                  death_date: values.death_date || undefined,
                  bio: values.bio || undefined,
                },
                {
                  onSuccess: (updated) => {
                    onUpdated?.(updated);
                    setSaveNotice("saved");
                  },
                  onError: () => setSaveNotice(null),
                },
              );
            },
            () => setSaveNotice("invalid"),
          )}
        >
          <h3 className="text-base font-semibold text-slate-900">Edit details</h3>

          <div>
            <Label htmlFor="edit-name" className="text-base">
              Name
            </Label>
            <Input
              id="edit-name"
              className="mt-1 py-3 text-base"
              error={form.formState.errors.first_name?.message}
              {...form.register("first_name")}
            />
          </div>

          <fieldset>
            <legend className="text-base font-medium text-slate-800">Gender</legend>
            <div className="mt-2 grid grid-cols-2 gap-2">
              {(["male", "female"] as const).map((value) => {
                const selected = form.watch("gender") === value;
                return (
                  <button
                    key={value}
                    type="button"
                    className={
                      selected
                        ? "rounded-xl bg-teal-800 px-3 py-3 text-base font-semibold capitalize text-white"
                        : "rounded-xl bg-teal-950/5 px-3 py-3 text-base font-semibold capitalize text-teal-950 ring-1 ring-teal-900/10"
                    }
                    onClick={() => form.setValue("gender", value, { shouldValidate: true })}
                  >
                    {value}
                  </button>
                );
              })}
            </div>
            {form.formState.errors.gender?.message && (
              <p className="mt-1 text-sm text-red-700">{form.formState.errors.gender.message}</p>
            )}
          </fieldset>

          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <Label htmlFor="edit-birth" className="text-base">
                Born (optional)
              </Label>
              <Input
                id="edit-birth"
                className="mt-1 py-3 text-base"
                placeholder="YYYY-MM-DD"
                error={form.formState.errors.birth_date?.message}
                {...form.register("birth_date")}
              />
            </div>
            <div>
              <Label htmlFor="edit-death" className="text-base">
                Died (optional)
              </Label>
              <Input
                id="edit-death"
                className="mt-1 py-3 text-base"
                placeholder="YYYY-MM-DD"
                error={form.formState.errors.death_date?.message}
                {...form.register("death_date")}
              />
            </div>
          </div>

          <div>
            <Label htmlFor="edit-bio" className="text-base">
              Note (optional)
            </Label>
            <Textarea
              id="edit-bio"
              className="mt-1 min-h-24 text-base"
              {...form.register("bio")}
            />
          </div>

          {saveNotice === "saved" && !updatePerson.isPending && (
            <Alert tone="success">Changes saved.</Alert>
          )}
          {saveNotice === "invalid" && (
            <Alert tone="error">Could not save. Check the name and use dates as YYYY-MM-DD.</Alert>
          )}
          {updatePerson.error instanceof ApiError && !updatePerson.isPending && (
            <Alert tone="error">{updatePerson.error.message}</Alert>
          )}

          <Button type="submit" size="lg" className="w-full py-3.5 text-base" loading={updatePerson.isPending}>
            Save changes
          </Button>
        </form>
      )}

      {detail.isPending ? (
        <div className="mt-4 grid place-items-center py-6">
          <Spinner className="size-6 text-teal-700" />
        </div>
      ) : (
        <div className="mt-6">
          <h3 className="text-base font-semibold text-slate-900">Relationships</h3>
          {detail.data?.relationships.length ? (
            <ul className="mt-2 space-y-2">
              {detail.data.relationships.map((rel) => (
                <li
                  key={rel.id}
                  className="flex items-center justify-between rounded-lg bg-teal-950/5 px-3 py-2 text-sm"
                >
                  <span className="capitalize text-slate-700">{rel.type}</span>
                  {canEdit && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => deleteRelationship.mutate(rel.id)}
                    >
                      Remove
                    </Button>
                  )}
                </li>
              ))}
            </ul>
          ) : (
            <p className="mt-2 text-sm text-slate-500">No relationships yet.</p>
          )}
        </div>
      )}

      <div className="mt-6 grid gap-4 sm:grid-cols-2">
        <TraversalList title="Ancestors" query={ancestors} onSelect={onSelectRelated} />
        <TraversalList title="Descendants" query={descendants} onSelect={onSelectRelated} />
      </div>

      {canEdit && (
        <Button
          className="mt-6"
          variant="danger"
          loading={deletePerson.isPending}
          onClick={() => {
            if (window.confirm(`Delete ${displayName(shown)} and their relationships?`)) {
              deletePerson.mutate(shown.id, { onSuccess: onClose });
            }
          }}
        >
          Delete member
        </Button>
      )}
    </div>
  );
}

function valuesFromPerson(person: Person): EditValues {
  const gender = person.gender === "female" ? "female" : "male";
  return {
    first_name: displayName(person),
    gender,
    birth_date: person.birth_date ?? "",
    death_date: person.death_date ?? "",
    bio: person.bio ?? "",
  };
}

function TraversalList({
  title,
  query,
  onSelect,
}: {
  title: string;
  query: { isPending: boolean; data?: PersonWithDepth[] };
  onSelect?: (person: Person) => void;
}) {
  return (
    <div>
      <h3 className="text-base font-semibold text-slate-900">{title}</h3>
      {query.isPending ? (
        <p className="mt-2 text-sm text-slate-500">Loading…</p>
      ) : query.data && query.data.length > 0 ? (
        <ul className="mt-2 space-y-1 text-sm text-slate-600">
          {query.data.map((item) => (
            <li key={item.id}>
              {onSelect ? (
                <button
                  type="button"
                  className="text-left font-medium text-teal-800 underline-offset-2 hover:underline"
                  onClick={() => onSelect(item)}
                >
                  {displayName(item)}
                </button>
              ) : (
                displayName(item)
              )}
              <span className="text-slate-400"> · gen {item.depth}</span>
            </li>
          ))}
        </ul>
      ) : (
        <p className="mt-2 text-sm text-slate-500">None within depth limit.</p>
      )}
    </div>
  );
}
