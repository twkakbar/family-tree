"use client";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import type { Person, Relationship } from "@/types/person";
import { displayName } from "@/utils/tree-layout";
import {
  allowedLinkChoices,
  type RelationshipLinkChoice,
} from "@/utils/relationship-rules";

export type { RelationshipLinkChoice };

interface RelationshipLinkPickerProps {
  source: Person;
  target: Person;
  relationships?: Relationship[];
  loading?: boolean;
  error?: string | null;
  onChoose: (choice: RelationshipLinkChoice) => void;
  onCancel: () => void;
}

/** Modal for choosing how two members are related after a drag-and-drop link. */
export function RelationshipLinkPicker({
  source,
  target,
  relationships = [],
  loading = false,
  error = null,
  onChoose,
  onCancel,
}: RelationshipLinkPickerProps) {
  const sourceName = displayName(source);
  const targetName = displayName(target);
  const choices = allowedLinkChoices(source, target, relationships);

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center p-4 sm:items-center" role="dialog">
      <button
        type="button"
        className="absolute inset-0 bg-slate-900/40 backdrop-blur-[1px]"
        aria-label="Close"
        onClick={onCancel}
      />

      <div
        data-link-picker
        className="relative z-10 w-full max-w-md rounded-2xl bg-white p-5 shadow-xl ring-1 ring-slate-200 sm:p-6"
      >
        <div className="flex items-start gap-3">
          <div className="flex min-w-0 flex-1 flex-col gap-1">
            <p className="text-base font-semibold text-slate-900">Connect members</p>
            <p className="text-sm text-slate-500">Choose how these two people are related.</p>
          </div>
        </div>

        <div className="mt-4 flex items-center gap-2 rounded-xl bg-slate-50 px-3 py-2.5 text-sm">
          <span className="truncate font-medium text-slate-900">{sourceName}</span>
          <span className="shrink-0 text-slate-400" aria-hidden>
            ↔
          </span>
          <span className="truncate font-medium text-slate-900">{targetName}</span>
        </div>

        {error && (
          <Alert tone="error" className="mt-4">
            {error}
          </Alert>
        )}

        <div className="mt-4 grid gap-2">
          {choices.length === 0 && (
            <p className="rounded-xl bg-slate-50 px-4 py-3 text-sm text-slate-600">
              These two are already connected in every way that fits a family tree. A spouse
              link stays hidden when they are already spouses, when a wife already has a
              husband, or when they are parent and child, siblings, or in the same bloodline.
              A husband can still take another wife.
            </p>
          )}
          {choices.includes("parent-of-target") && (
            <ChoiceButton
              title={`${sourceName} is the parent`}
              subtitle={`of ${targetName}`}
              loading={loading}
              onClick={() => onChoose("parent-of-target")}
            />
          )}
          {choices.includes("child-of-target") && (
            <ChoiceButton
              title={`${targetName} is the parent`}
              subtitle={`of ${sourceName}`}
              loading={loading}
              onClick={() => onChoose("child-of-target")}
            />
          )}
          {choices.includes("spouse") && (
            <ChoiceButton
              title="Spouse or partner"
              subtitle={`Link ${sourceName} and ${targetName}`}
              loading={loading}
              onClick={() => onChoose("spouse")}
            />
          )}
        </div>

        <Button type="button" variant="ghost" className="mt-4 w-full" onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </div>
  );
}

function ChoiceButton({
  title,
  subtitle,
  loading,
  onClick,
}: {
  title: string;
  subtitle: string;
  loading: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      disabled={loading}
      onClick={onClick}
      className="rounded-xl border border-slate-200 bg-white px-4 py-3 text-left transition hover:border-emerald-300 hover:bg-emerald-50/50 disabled:opacity-60"
    >
      <span className="block text-sm font-semibold text-slate-900">{title}</span>
      <span className="mt-0.5 block text-xs text-slate-500">{subtitle}</span>
    </button>
  );
}

export function choiceToRelationshipInput(
  choice: RelationshipLinkChoice,
  source: Person,
  target: Person,
) {
  switch (choice) {
    case "parent-of-target":
      return {
        from_person_id: source.id,
        to_person_id: target.id,
        type: "parent" as const,
      };
    case "child-of-target":
      return {
        from_person_id: target.id,
        to_person_id: source.id,
        type: "parent" as const,
      };
    case "spouse":
      return {
        from_person_id: source.id,
        to_person_id: target.id,
        type: "spouse" as const,
      };
  }
}
