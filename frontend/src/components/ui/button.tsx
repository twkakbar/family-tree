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
    "bg-teal-800 text-white hover:bg-teal-900 focus-visible:outline-teal-800 shadow-sm shadow-teal-950/10",
  secondary:
    "bg-white text-teal-950 ring-1 ring-inset ring-teal-900/15 hover:bg-teal-50 focus-visible:outline-teal-700",
  ghost: "text-teal-900 hover:bg-teal-900/5 focus-visible:outline-teal-700",
  danger: "bg-red-700 text-white hover:bg-red-800 focus-visible:outline-red-700 shadow-sm",
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
