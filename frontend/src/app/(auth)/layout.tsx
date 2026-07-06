import type { ReactNode } from "react";

import { Logo } from "@/components/layout/logo";

/** Centered card layout shared by the login and register pages. */
export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-dvh flex-col items-center justify-center px-4 py-10">
      <Logo className="mb-8" />
      <div className="w-full max-w-md">{children}</div>
    </div>
  );
}
