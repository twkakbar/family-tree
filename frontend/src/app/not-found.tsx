import Link from "next/link";

import { Button } from "@/components/ui/button";

export default function NotFound() {
  return (
    <div className="grid min-h-dvh place-items-center px-4">
      <div className="text-center">
        <p className="text-sm font-semibold text-emerald-600">404</p>
        <h1 className="mt-2 text-3xl font-bold text-slate-900">Page not found</h1>
        <p className="mt-2 text-slate-600">The page you are looking for does not exist.</p>
        <Link href="/" className="mt-6 inline-block">
          <Button variant="secondary">Back to home</Button>
        </Link>
      </div>
    </div>
  );
}
