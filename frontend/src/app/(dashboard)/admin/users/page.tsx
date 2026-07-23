"use client";

import { useState } from "react";
import Link from "next/link";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { useMe } from "@/hooks/use-auth";
import {
  useAdminUsers,
  useDeleteAdminUser,
  useIssuePasswordReset,
  useUpdateAdminUser,
} from "@/hooks/use-admin-users";
import { ApiError } from "@/types/api";
import type { Role, User } from "@/types/user";

export default function AdminUsersPage() {
  const { data: me } = useMe();
  const [search, setSearch] = useState("");
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(1);
  const [issuedToken, setIssuedToken] = useState<{ email: string; token: string } | null>(
    null,
  );

  const users = useAdminUsers(query, page, me?.role === "admin");
  const updateUser = useUpdateAdminUser();
  const deleteUser = useDeleteAdminUser();
  const issueReset = useIssuePasswordReset();

  if (me && me.role !== "admin") {
    return <Alert tone="error">Only administrators can manage accounts.</Alert>;
  }

  const mutationError =
    (updateUser.error instanceof ApiError && !updateUser.isPending
      ? updateUser.error.message
      : null) ||
    (deleteUser.error instanceof ApiError && !deleteUser.isPending
      ? deleteUser.error.message
      : null) ||
    (issueReset.error instanceof ApiError && !issueReset.isPending
      ? issueReset.error.message
      : null);

  return (
    <div className="mx-auto max-w-3xl space-y-8">
      <div>
        <p className="text-sm font-semibold uppercase tracking-wider text-teal-800/70">Admin</p>
        <h1 className="font-display mt-1 text-4xl text-teal-950">Accounts</h1>
        <p className="mt-3 text-lg text-slate-600">
          Search members, change roles, mark email verified, or issue a password-reset token.{" "}
          <Link
            href="/admin/requests"
            className="font-semibold text-teal-800 underline-offset-2 hover:underline"
          >
            Edit requests
          </Link>
        </p>
      </div>

      <form
        className="flex flex-col gap-3 sm:flex-row"
        onSubmit={(e) => {
          e.preventDefault();
          setPage(1);
          setQuery(search.trim());
        }}
      >
        <Input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search by name or email"
          className="py-3 text-base"
          aria-label="Search accounts"
        />
        <Button type="submit" size="lg" className="py-3 text-base sm:px-8">
          Find
        </Button>
      </form>

      {mutationError && <Alert tone="error">{mutationError}</Alert>}

      {issuedToken && (
        <Alert tone="success">
          Reset token for {issuedToken.email} (share out-of-band):
          <span className="mt-2 block break-all font-mono text-sm">{issuedToken.token}</span>
          <Link
            href={`/reset-password?token=${encodeURIComponent(issuedToken.token)}`}
            className="mt-2 inline-block font-semibold underline"
          >
            Open reset page
          </Link>
        </Alert>
      )}

      {users.isPending ? (
        <div className="grid place-items-center py-16">
          <Spinner className="size-10 text-teal-700" />
        </div>
      ) : users.data && users.data.items.length > 0 ? (
        <ul className="space-y-5">
          {users.data.items.map((user) => (
            <AdminUserRow
              key={user.id}
              user={user}
              isSelf={user.id === me?.id}
              busy={updateUser.isPending || deleteUser.isPending || issueReset.isPending}
              onRole={(role) => updateUser.mutate({ id: user.id, input: { role } })}
              onVerified={(email_verified) =>
                updateUser.mutate({ id: user.id, input: { email_verified } })
              }
              onReset={() =>
                issueReset.mutate(user.id, {
                  onSuccess: (data) =>
                    setIssuedToken({ email: user.email, token: data.reset_token }),
                })
              }
              onDelete={() => {
                if (
                  window.confirm(
                    `Soft-delete ${user.email}? Their sessions will be revoked.`,
                  )
                ) {
                  deleteUser.mutate(user.id);
                }
              }}
            />
          ))}
        </ul>
      ) : (
        <div className="aceh-panel">
          <h2 className="font-display text-2xl text-teal-950">No accounts found</h2>
          <p className="mt-2 text-lg text-slate-600">Try a different search.</p>
        </div>
      )}

      {users.data && users.data.meta.total_pages > 1 && (
        <div className="flex items-center justify-between gap-4">
          <Button
            variant="secondary"
            size="lg"
            disabled={page <= 1}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
          >
            Previous
          </Button>
          <p className="text-base text-slate-600">
            Page {users.data.meta.page} of {users.data.meta.total_pages}
          </p>
          <Button
            variant="secondary"
            size="lg"
            disabled={page >= users.data.meta.total_pages}
            onClick={() => setPage((p) => p + 1)}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  );
}

function AdminUserRow({
  user,
  isSelf,
  busy,
  onRole,
  onVerified,
  onReset,
  onDelete,
}: {
  user: User;
  isSelf: boolean;
  busy: boolean;
  onRole: (role: Role) => void;
  onVerified: (verified: boolean) => void;
  onReset: () => void;
  onDelete: () => void;
}) {
  return (
    <li className="aceh-panel space-y-4">
      <div>
        <h2 className="font-display text-2xl text-teal-950">{user.name}</h2>
        <p className="mt-1 text-base text-slate-600">{user.email}</p>
        <p className="mt-2 text-sm text-slate-500">
          Role: <span className="font-semibold text-teal-900">{user.role}</span>
          {" · "}
          Email:{" "}
          <span className="font-semibold text-teal-900">
            {user.email_verified ? "verified" : "unverified"}
          </span>
        </p>
      </div>

      <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap">
        {user.role !== "admin" ? (
          <Button size="lg" className="py-3 text-base" disabled={busy} onClick={() => onRole("admin")}>
            Make admin
          </Button>
        ) : (
          <Button
            size="lg"
            variant="secondary"
            className="py-3 text-base"
            disabled={busy || isSelf}
            onClick={() => onRole("user")}
          >
            Make user
          </Button>
        )}

        <Button
          size="lg"
          variant="secondary"
          className="py-3 text-base"
          disabled={busy}
          onClick={() => onVerified(!user.email_verified)}
        >
          {user.email_verified ? "Mark unverified" : "Mark verified"}
        </Button>

        <Button size="lg" variant="secondary" className="py-3 text-base" disabled={busy} onClick={onReset}>
          Issue password reset
        </Button>

        <Button
          size="lg"
          variant="danger"
          className="py-3 text-base"
          disabled={busy || isSelf}
          onClick={onDelete}
        >
          Soft-delete
        </Button>
      </div>
    </li>
  );
}
