import Link from "next/link";

import { env } from "@/lib/env";
import { cn } from "@/utils/cn";

export function Logo({ className }: { className?: string }) {
  return (
    <Link
      href="/"
      className={cn(
        "font-display inline-flex items-center gap-3 text-2xl tracking-wide text-teal-950",
        className,
      )}
    >
      <span
        className="grid size-10 place-items-center rounded-xl bg-amber-400 text-teal-950 shadow-md shadow-amber-900/20"
        aria-hidden
      >
        <RencongIcon />
      </span>
      {env.appName}
    </Link>
  );
}

/** Stylized Acehnese rencong — the classic curved dagger emblem of Aceh. */
function RencongIcon() {
  return (
    <svg viewBox="0 0 24 24" className="size-6" fill="currentColor" aria-hidden>
      {/* Blade curve */}
      <path d="M18.6 3.4c-1.8 0.2-4.6 1.4-7.4 4.1C8.2 10.4 5.8 14.2 4.8 17.2c-0.2 0.6 0.5 1.1 1 0.7 2.4-1.8 5.2-5.4 7.4-8.2 1.9-2.4 3.6-4.2 4.9-5.1 0.5-0.3 0.7-1.1 0.5-1.2z" />
      {/* Hilt / grip */}
      <path d="M4.2 17.6c-0.9 0.8-1.5 1.9-1.4 2.6 0.1 0.7 0.9 1.1 1.8 0.9 0.7-0.2 1.5-0.8 2.1-1.5-0.9-0.4-1.8-1.1-2.5-2z" />
      {/* Cross guard accent */}
      <path d="M6.6 15.8c0.5 0.6 1.1 1.1 1.7 1.4l0.9-1.2c-0.5-0.3-1-0.7-1.4-1.2l-1.2 1z" />
    </svg>
  );
}
