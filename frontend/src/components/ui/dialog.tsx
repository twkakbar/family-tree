"use client";

import { useEffect, useId, useRef, type ReactNode } from "react";
import { createPortal } from "react-dom";

import { cn } from "@/utils/cn";

interface DialogProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  className?: string;
}

/**
 * Lightweight modal overlay (no new dependency). Focuses the panel on open
 * and closes on Escape / backdrop click.
 */
export function Dialog({ open, onClose, title, children, className }: DialogProps) {
  const titleId = useId();
  const panelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;

    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKeyDown);
    panelRef.current?.focus();

    return () => {
      document.body.style.overflow = previousOverflow;
      window.removeEventListener("keydown", onKeyDown);
    };
  }, [open, onClose]);

  if (!open || typeof document === "undefined") {
    return null;
  }

  return createPortal(
    <div className="fixed inset-0 z-50 flex items-end justify-center p-0 sm:items-center sm:p-4">
      <button
        type="button"
        aria-label="Close dialog"
        className="absolute inset-0 bg-teal-950/55 backdrop-blur-[2px]"
        onClick={onClose}
      />
      <div
        ref={panelRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby={title ? titleId : undefined}
        tabIndex={-1}
        className={cn(
          "relative z-10 flex max-h-[min(92dvh,44rem)] w-full max-w-lg flex-col overflow-hidden",
          "rounded-t-3xl bg-white shadow-2xl ring-1 ring-teal-900/10 sm:rounded-3xl",
          "outline-none",
          className,
        )}
      >
        {title ? (
          <h2 id={titleId} className="sr-only">
            {title}
          </h2>
        ) : null}
        <div className="overflow-y-auto overscroll-contain p-5 sm:p-6">{children}</div>
      </div>
    </div>,
    document.body,
  );
}
