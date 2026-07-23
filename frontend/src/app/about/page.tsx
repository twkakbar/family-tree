import Link from "next/link";

export const metadata = {
  title: "About",
  description: "About the public Aceh royal family tree.",
};

/** Brand / story page — no login or edit-request CTAs. */
export default function AboutPage() {
  return (
    <div className="aceh-sea-bg relative min-h-dvh overflow-hidden text-teal-50">
      <div
        className="pointer-events-none absolute inset-0 opacity-[0.12]"
        style={{
          backgroundImage:
            "url(\"data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%232dd4bf' fill-opacity='0.35'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E\")",
        }}
      />

      <main className="relative mx-auto flex min-h-dvh max-w-3xl flex-col justify-center px-6 py-16 sm:px-10">
        <p className="text-sm font-medium uppercase tracking-[0.28em] text-amber-200/80">
          Aceh · Indonesia
        </p>
        <h1 className="font-display mt-4 text-5xl leading-[1.05] text-amber-50 sm:text-6xl">
          Kesultanan Aceh
        </h1>
        <p className="mt-6 text-lg leading-relaxed text-teal-100/85 sm:text-xl">
          This site shares the living royal lineage of Aceh — an open tree anyone can explore,
          from today&apos;s descendants back toward the sultan line.
        </p>
        <p className="mt-4 text-base leading-relaxed text-teal-100/70 sm:text-lg">
          Drag, zoom, and tap names in the 3D view to learn more about each person. The tree grows
          as trusted family editors add members over time.
        </p>

        <div className="mt-10">
          <Link
            href="/"
            className="inline-flex items-center justify-center rounded-full bg-amber-400 px-8 py-3.5 text-base font-semibold text-slate-950 shadow-lg shadow-amber-900/30 transition hover:bg-amber-300"
          >
            Back to the tree
          </Link>
        </div>
      </main>
    </div>
  );
}
