import { redirect } from "next/navigation";

/** Multi-family detail UI retired; editors use /studio. */
export default function LegacyFamilyPage() {
  redirect("/studio");
}
