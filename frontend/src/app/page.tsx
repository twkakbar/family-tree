import Link from "next/link";

import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Logo } from "@/components/layout/logo";
import { env } from "@/lib/env";

const features = [
  {
    title: "Build your tree",
    description: "Add family members, connect generations and watch your history take shape.",
  },
  {
    title: "Collaborate",
    description: "Share a family with relatives as editors or viewers - everyone contributes.",
  },
  {
    title: "Own your data",
    description: "Self-hosted, exportable and importable. Your family history stays yours.",
  },
];

export default function LandingPage() {
  return (
    <div className="flex min-h-dvh flex-col">
      <header className="mx-auto flex w-full max-w-6xl items-center justify-between px-4 py-5 sm:px-6">
        <Logo />
        <nav className="flex items-center gap-2">
          <Link href="/login">
            <Button variant="ghost">Sign in</Button>
          </Link>
          <Link href="/register">
            <Button>Get started</Button>
          </Link>
        </nav>
      </header>

      <main className="flex-1">
        <section className="mx-auto max-w-3xl px-4 py-20 text-center sm:px-6 sm:py-28">
          <h1 className="text-4xl font-bold tracking-tight text-slate-900 sm:text-5xl">
            Preserve your family history,
            <span className="text-emerald-600"> beautifully</span>
          </h1>
          <p className="mx-auto mt-5 max-w-xl text-lg text-slate-600">
            {env.appName} helps you build interactive family trees, connect generations and share
            your heritage with the people who matter.
          </p>
          <div className="mt-8 flex justify-center gap-3">
            <Link href="/register">
              <Button size="lg">Create your family tree</Button>
            </Link>
            <Link href="/login">
              <Button size="lg" variant="secondary">
                Sign in
              </Button>
            </Link>
          </div>
        </section>

        <section className="mx-auto grid max-w-6xl gap-4 px-4 pb-24 sm:grid-cols-3 sm:px-6">
          {features.map((feature) => (
            <Card key={feature.title}>
              <CardTitle>{feature.title}</CardTitle>
              <CardDescription>{feature.description}</CardDescription>
            </Card>
          ))}
        </section>
      </main>

      <footer className="border-t border-slate-200 py-6 text-center text-sm text-slate-500">
        {env.appName} — self-hosted family history
      </footer>
    </div>
  );
}
