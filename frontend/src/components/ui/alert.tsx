import type { HTMLAttributes } from "react";

import { cn } from "@/utils/cn";

type Tone = "error" | "success" | "info";

const toneClasses: Record<Tone, string> = {
  error: "bg-red-50 text-red-800 ring-red-200",
  success: "bg-emerald-50 text-emerald-800 ring-emerald-200",
  info: "bg-sky-50 text-sky-800 ring-sky-200",
};

export interface AlertProps extends HTMLAttributes<HTMLDivElement> {
  tone?: Tone;
}

export function Alert({ className, tone = "info", ...props }: AlertProps) {
  return (
    <div
      role="alert"
      className={cn("rounded-lg p-3 text-sm ring-1 ring-inset", toneClasses[tone], className)}
      {...props}
    />
  );
}
