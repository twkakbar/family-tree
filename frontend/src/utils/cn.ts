import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Merge conditional class names with correct Tailwind conflict resolution
 * (clsx for conditionals, tailwind-merge so "p-2 p-4" resolves to "p-4").
 * Used by every UI primitive.
 */
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}
