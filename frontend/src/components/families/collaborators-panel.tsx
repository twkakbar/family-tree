"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  useFamilyCollaborators,
  useInviteCollaborator,
  useRemoveCollaborator,
  useUpdateCollaborator,
} from "@/hooks/use-families";
import { inviteCollaboratorSchema, type InviteCollaboratorFormValues } from "@/lib/validations/family";
import { ApiError } from "@/types/api";
import type { Collaborator, Family, FamilyRole } from "@/types/family";

interface CollaboratorsPanelProps {
  family: Family;
}

const editableRoles: FamilyRole[] = ["owner", "editor", "viewer"];

/** Lists collaborators and lets owners invite or manage access. */
export function CollaboratorsPanel({ family }: CollaboratorsPanelProps) {
  const collaborators = useFamilyCollaborators(family.id);
  const invite = useInviteCollaborator(family.id);
  const updateRole = useUpdateCollaborator(family.id);
  const remove = useRemoveCollaborator(family.id);

  const form = useForm<InviteCollaboratorFormValues>({
    resolver: zodResolver(inviteCollaboratorSchema),
    defaultValues: { email: "", role: "editor" },
  });

  const apiError =
    (invite.error instanceof ApiError && !invite.isPending ? invite.error.message : null) ??
    (updateRole.error instanceof ApiError && !updateRole.isPending ? updateRole.error.message : null) ??
    (remove.error instanceof ApiError && !remove.isPending ? remove.error.message : null);

  return (
    <Card>
      <CardTitle>Collaborators</CardTitle>
      <CardDescription>
        {family.role === "owner"
          ? "Invite existing users by email and manage their roles."
          : "People who can access this family."}
      </CardDescription>

      {apiError && (
        <Alert tone="error" className="mt-4">
          {apiError}
        </Alert>
      )}

      {collaborators.isPending ? (
        <p className="mt-4 text-sm text-slate-500">Loading collaborators…</p>
      ) : (
        <ul className="mt-4 divide-y divide-slate-200 rounded-lg ring-1 ring-slate-200">
          {(collaborators.data ?? []).map((collaborator) => (
            <CollaboratorRow
              key={collaborator.user_id}
              collaborator={collaborator}
              canManage={family.role === "owner"}
              onRoleChange={(role) =>
                updateRole.mutate({ userId: collaborator.user_id, input: { role } })
              }
              onRemove={() => {
                if (
                  window.confirm(`Remove ${collaborator.name} from this family?`) &&
                  !remove.isPending
                ) {
                  remove.mutate(collaborator.user_id);
                }
              }}
            />
          ))}
        </ul>
      )}

      {family.role === "owner" && (
        <form
          className="mt-6 space-y-4 border-t border-slate-200 pt-6"
          noValidate
          onSubmit={form.handleSubmit((values) => {
            invite.mutate(values, {
              onSuccess: () => form.reset({ email: "", role: "editor" }),
            });
          })}
        >
          <div>
            <Label htmlFor="invite-email">Invite by email</Label>
            <Input
              id="invite-email"
              type="email"
              placeholder="collaborator@example.com"
              error={form.formState.errors.email?.message}
              {...form.register("email")}
            />
          </div>

          <div>
            <Label htmlFor="invite-role">Role</Label>
            <select
              id="invite-role"
              className="block w-full rounded-lg border-0 px-3 py-2 text-sm text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 focus:ring-2 focus:ring-inset focus:ring-emerald-600"
              {...form.register("role")}
            >
              <option value="editor">Editor — can edit members and relationships</option>
              <option value="viewer">Viewer — read-only access</option>
            </select>
          </div>

          <Button type="submit" loading={invite.isPending}>
            Send invite
          </Button>
        </form>
      )}
    </Card>
  );
}

interface CollaboratorRowProps {
  collaborator: Collaborator;
  canManage: boolean;
  onRoleChange: (role: FamilyRole) => void;
  onRemove: () => void;
}

function CollaboratorRow({ collaborator, canManage, onRoleChange, onRemove }: CollaboratorRowProps) {
  return (
    <li className="flex flex-col gap-3 px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p className="text-sm font-medium text-slate-900">{collaborator.name}</p>
        <p className="text-xs text-slate-500">{collaborator.email}</p>
      </div>

      {canManage ? (
        <div className="flex flex-wrap items-center gap-2">
          <select
            value={collaborator.role}
            className="rounded-lg border-0 px-2 py-1.5 text-sm text-slate-900 ring-1 ring-inset ring-slate-300 focus:ring-2 focus:ring-emerald-600"
            onChange={(event) => onRoleChange(event.target.value as FamilyRole)}
          >
            {editableRoles.map((role) => (
              <option key={role} value={role}>
                {role}
              </option>
            ))}
          </select>
          <Button type="button" variant="ghost" size="sm" onClick={onRemove}>
            Remove
          </Button>
        </div>
      ) : (
        <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-600">
          {collaborator.role}
        </span>
      )}
    </li>
  );
}
