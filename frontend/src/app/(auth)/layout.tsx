import type { ReactNode } from "react";
import Link from "next/link";

import { Logo } from "@/components/layout/logo";

/** Centered auth layout with Aceh sea backdrop. Login/register stay unlinked from public pages. */
export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="aceh-sea-bg flex min-h-dvh flex-col items-center justify-center px-4 py-10">
      <Logo className="mb-8 text-amber-50" />
      <div className="w-full max-w-md">{children}</div>
      <Link
        href="/"
        className="mt-8 text-base text-teal-100/70 underline-offset-2 hover:text-amber-100 hover:underline"
      >
        ← Back to the tree
      </Link>
    </div>
  );
}
