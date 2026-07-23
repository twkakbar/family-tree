"use client";

import Link from "next/link";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { useMe } from "@/hooks/use-auth";
import { usePendingAccessRequests, useReviewAccessRequest } from "@/hooks/use-public";
import { ApiError } from "@/types/api";

export default function AdminRequestsPage() {
  const { data: user } = useMe();
  const pending = usePendingAccessRequests(user?.role === "admin");
  const review = useReviewAccessRequest();

  if (user && user.role !== "admin") {
    return <Alert tone="error">Only administrators can review edit-access requests.</Alert>;
  }

  const apiError =
    review.error instanceof ApiError && !review.isPending ? review.error.message : null;

  return (
    <div className="mx-auto max-w-2xl space-y-8">
      <div>
        <p className="text-sm font-semibold uppercase tracking-wider text-teal-800/70">Admin</p>
        <h1 className="font-display mt-1 text-4xl text-teal-950">Who may edit?</h1>
        <p className="mt-3 text-lg text-slate-600">
          Approve people who asked to help grow the royal tree.{" "}
          <Link href="/admin/users" className="font-semibold text-teal-800 underline-offset-2 hover:underline">
            Manage accounts
          </Link>
          {" · "}
          <Link href="/studio" className="font-semibold text-teal-800 underline-offset-2 hover:underline">
            Open studio
          </Link>
        </p>
      </div>

      {apiError && <Alert tone="error">{apiError}</Alert>}

      {pending.isPending ? (
        <div className="grid place-items-center py-16">
          <Spinner className="size-10 text-teal-700" />
        </div>
      ) : pending.data && pending.data.length > 0 ? (
        <ul className="space-y-5">
          {pending.data.map((item) => (
            <li key={item.id} className="aceh-panel">
              <h2 className="font-display text-2xl text-teal-950">{item.name}</h2>
              <p className="mt-1 text-base text-slate-600">{item.email}</p>
              <p className="mt-5 whitespace-pre-wrap text-lg leading-relaxed text-slate-800">
                {item.message}
              </p>
              <div className="mt-6 flex flex-col gap-3 sm:flex-row">
                <Button
                  size="lg"
                  className="flex-1 py-4 text-lg"
                  loading={review.isPending}
                  onClick={() => review.mutate({ id: item.id, action: "approve" })}
                >
                  Approve — allow editing
                </Button>
                <Button
                  size="lg"
                  variant="secondary"
                  className="flex-1 py-4 text-lg"
                  loading={review.isPending}
                  onClick={() => review.mutate({ id: item.id, action: "reject" })}
                >
                  Decline
                </Button>
              </div>
            </li>
          ))}
        </ul>
      ) : (
        <div className="aceh-panel">
          <h2 className="font-display text-2xl text-teal-950">No waiting requests</h2>
          <p className="mt-2 text-lg text-slate-600">
            When someone asks for edit access, their message will appear here.
          </p>
        </div>
      )}
    </div>
  );
}
