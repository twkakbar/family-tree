import { redirect } from "next/navigation";

/** Kept for old links; the tree now lives at `/`. */
export default function ExploreRedirectPage() {
  redirect("/");
}
