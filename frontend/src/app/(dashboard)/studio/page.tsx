"use client";

import { useState } from "react";
import Link from "next/link";

import { ExportImportPanel } from "@/components/persons/export-import-panel";
import { MemberSearch } from "@/components/persons/member-search";
import { PersonDetailPanel } from "@/components/persons/person-detail-panel";
import { SimpleMemberForm } from "@/components/persons/simple-member-form";
import { FamilyTreeView } from "@/components/tree/family-tree-view";
import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Dialog } from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";
import { useRoyalMeta, useStudioCanEdit } from "@/hooks/use-public";
import type { Family } from "@/types/family";
import type { Person } from "@/types/person";

export default function StudioPage() {
  const canEdit = useStudioCanEdit();
  const meta = useRoyalMeta();
  const [selectedPerson, setSelectedPerson] = useState<Person | null>(null);
  const [showTree, setShowTree] = useState(false);

  if (canEdit.isPending) {
    return (
      <div className="grid place-items-center py-24">
        <Spinner className="size-10 text-teal-700" />
      </div>
    );
  }

  if (!canEdit.data?.can_edit || !canEdit.data.family_id) {
    return (
      <div className="aceh-panel mx-auto max-w-lg space-y-4">
        <Alert tone="info">
          Your account cannot edit the tree yet. An administrator must approve edit access first.
        </Alert>
        <Link href="/" className="inline-block text-lg font-medium text-teal-800 hover:text-teal-700">
          ← Back to the tree
        </Link>
      </div>
    );
  }

  const family: Family = {
    id: canEdit.data.family_id,
    name: meta.data?.name ?? "Kesultanan Aceh",
    description: meta.data?.description ?? "",
    role: "editor",
    persons_count: 0,
    created_at: new Date(0).toISOString(),
    updated_at: new Date(0).toISOString(),
  };

  return (
    <div className="mx-auto max-w-3xl space-y-8 xl:max-w-5xl">
      <div>
        <p className="text-sm font-semibold uppercase tracking-wider text-teal-800/70">Editor studio</p>
        <h1 className="font-display mt-1 text-4xl text-teal-950 sm:text-5xl">
          {family.name}
        </h1>
        <p className="mt-3 max-w-2xl text-lg text-slate-600">
          Add one person at a time. Start with a living relative, then connect them toward older
          generations.{" "}
          <Link href="/" className="font-semibold text-teal-800 underline-offset-2 hover:underline">
            Open the public tree
          </Link>
        </p>
      </div>

      <SimpleMemberForm
        family={family}
        onCreated={(person) => {
          setSelectedPerson(person);
          setShowTree(true);
        }}
      />

      <MemberSearch
        familyId={family.id}
        selectedId={selectedPerson?.id}
        onSelect={(person) => {
          setSelectedPerson(person);
          setShowTree(true);
        }}
      />

      <ExportImportPanel familyId={family.id} />

      <div className="aceh-panel space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 className="font-display text-2xl text-teal-950">Tree diagram</h2>
            <p className="mt-1 text-base text-slate-600">
              Tap a name to open details in a popup — edit, photo, or remove links without scrolling.
            </p>
          </div>
          <Button
            type="button"
            variant="secondary"
            size="lg"
            className="text-base"
            onClick={() => setShowTree((v) => !v)}
          >
            {showTree ? "Hide diagram" : "Show diagram"}
          </Button>
        </div>

        {showTree && (
          <FamilyTreeView
            familyId={family.id}
            canEdit
            selectedPersonId={selectedPerson?.id}
            onSelectPerson={setSelectedPerson}
          />
        )}
      </div>

      <Dialog
        open={selectedPerson !== null}
        onClose={() => setSelectedPerson(null)}
        title={selectedPerson ? `Member: ${selectedPerson.first_name}` : "Member details"}
      >
        {selectedPerson && (
          <PersonDetailPanel
            family={family}
            person={selectedPerson}
            onClose={() => setSelectedPerson(null)}
            onUpdated={setSelectedPerson}
            onSelectRelated={setSelectedPerson}
          />
        )}
      </Dialog>
    </div>
  );
}
