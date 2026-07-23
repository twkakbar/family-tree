import { forwardRef, type InputHTMLAttributes } from "react";

import { cn } from "@/utils/cn";

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  /** Validation message; also switches the border to the error state. */
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className, error, id, ...props }, ref) => (
    <div>
      <input
        ref={ref}
        id={id}
        aria-invalid={error ? true : undefined}
        aria-describedby={error && id ? `${id}-error` : undefined}
        className={cn(
          "block w-full rounded-lg border-0 px-3 py-2 text-sm text-slate-900 shadow-sm",
          "ring-1 ring-inset ring-slate-300 placeholder:text-slate-400",
          "focus:ring-2 focus:ring-inset focus:ring-teal-700",
          "disabled:cursor-not-allowed disabled:bg-slate-50",
          error && "ring-red-400 focus:ring-red-500",
          className,
        )}
        {...props}
      />
      {error && (
        <p id={id ? `${id}-error` : undefined} className="mt-1.5 text-sm text-red-600">
          {error}
        </p>
      )}
    </div>
  ),
);
Input.displayName = "Input";
