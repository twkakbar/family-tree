import Link from "next/link";

import { env } from "@/lib/env";
import { cn } from "@/utils/cn";

export function Logo({ className }: { className?: string }) {
  return (
    <Link
      href="/"
      className={cn("inline-flex items-center gap-2 font-semibold text-slate-900", className)}
    >
      <span className="grid size-8 place-items-center rounded-lg bg-emerald-600 text-white">
        <TreeIcon />
      </span>
      {env.appName}
    </Link>
  );
}

function TreeIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="size-5">
      <circle cx="12" cy="5" r="2.2" />
      <circle cx="5.5" cy="18" r="2.2" />
      <circle cx="18.5" cy="18" r="2.2" />
      <path d="M12 7.5v4m0 0-5 4.5m5-4.5 5 4.5" strokeLinecap="round" />
    </svg>
  );
}
