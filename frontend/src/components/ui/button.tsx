import { forwardRef, type ButtonHTMLAttributes } from "react";

import { cn } from "@/utils/cn";
import { Spinner } from "./spinner";

type Variant = "primary" | "secondary" | "ghost" | "danger";
type Size = "sm" | "md" | "lg";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  size?: Size;
  /** Shows a spinner and disables the button. */
  loading?: boolean;
}

const variantClasses: Record<Variant, string> = {
  primary:
    "bg-emerald-600 text-white hover:bg-emerald-700 focus-visible:outline-emerald-600 shadow-sm",
  secondary:
    "bg-white text-slate-900 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus-visible:outline-slate-400",
  ghost: "text-slate-700 hover:bg-slate-100 focus-visible:outline-slate-400",
  danger: "bg-red-600 text-white hover:bg-red-700 focus-visible:outline-red-600 shadow-sm",
};

const sizeClasses: Record<Size, string> = {
  sm: "px-2.5 py-1.5 text-sm",
  md: "px-3.5 py-2 text-sm",
  lg: "px-4 py-2.5 text-base",
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "primary", size = "md", loading = false, disabled, children, ...props }, ref) => (
    <button
      ref={ref}
      disabled={disabled || loading}
      className={cn(
        "inline-flex items-center justify-center gap-2 rounded-lg font-semibold transition-colors",
        "focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2",
        "disabled:cursor-not-allowed disabled:opacity-60",
        variantClasses[variant],
        sizeClasses[size],
        className,
      )}
      {...props}
    >
      {loading && <Spinner className="size-4" />}
      {children}
    </button>
  ),
);
Button.displayName = "Button";
