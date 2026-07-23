"use client";

import { useRef, useState } from "react";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { useExportTree, useImportTree } from "@/hooks/use-transfer";
import { ApiError } from "@/types/api";
import type { TreeTransfer } from "@/services/transfer.service";

interface ExportImportPanelProps {
  familyId: string;
}

/** Large-button backup / bulk-load for the royal tree (JSON only). */
export function ExportImportPanel({ familyId }: ExportImportPanelProps) {
  const exportTree = useExportTree(familyId);
  const importTree = useImportTree(familyId);
  const fileRef = useRef<HTMLInputElement>(null);
  const [message, setMessage] = useState<string | null>(null);

  const apiError =
    (exportTree.error instanceof ApiError && !exportTree.isPending
      ? exportTree.error.message
      : null) ??
    (importTree.error instanceof ApiError && !importTree.isPending
      ? importTree.error.message
      : null);

  const download = async () => {
    setMessage(null);
    try {
      const doc = await exportTree.mutateAsync();
      const blob = new Blob([JSON.stringify(doc, null, 2)], { type: "application/json" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "kesultanan-aceh-export.json";
      a.click();
      URL.revokeObjectURL(url);
      setMessage(`Exported ${doc.persons.length} members.`);
    } catch {
      // surfaced via apiError
    }
  };

  const onFile = async (file: File) => {
    setMessage(null);
    try {
      const text = await file.text();
      const doc = JSON.parse(text) as TreeTransfer;
      if (!window.confirm(
        `Import this file into the tree?\n\nIt will ADD new members (not replace). Continue?`,
      )) {
        return;
      }
      const result = await importTree.mutateAsync(doc);
      setMessage(
        `Imported ${result.persons_created} members and ${result.relationships_created} links.`,
      );
    } catch (err) {
      if (!(err instanceof ApiError)) {
        setMessage("Could not read that file. Use a JSON export from this site.");
      }
    }
  };

  return (
    <div className="aceh-panel space-y-4">
      <div>
        <h2 className="font-display text-2xl text-teal-950">Backup &amp; load</h2>
        <p className="mt-1 text-base text-slate-600">
          Download a JSON copy of the tree, or load a previous export. Import always{" "}
          <strong>adds</strong> people — it does not delete existing ones.
        </p>
      </div>

      {apiError && <Alert tone="error">{apiError}</Alert>}
      {message && !apiError && <Alert tone="success">{message}</Alert>}

      <div className="flex flex-col gap-3 sm:flex-row">
        <Button
          type="button"
          size="lg"
          className="py-4 text-lg sm:flex-1"
          loading={exportTree.isPending}
          onClick={() => void download()}
        >
          Download JSON backup
        </Button>
        <Button
          type="button"
          size="lg"
          variant="secondary"
          className="py-4 text-lg sm:flex-1"
          loading={importTree.isPending}
          onClick={() => fileRef.current?.click()}
        >
          Load JSON file
        </Button>
        <input
          ref={fileRef}
          type="file"
          accept="application/json,.json"
          className="sr-only"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) void onFile(file);
            e.target.value = "";
          }}
        />
      </div>
    </div>
  );
}
