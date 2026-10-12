export function TreeLegend({ tone }: { tone: "light" | "dark" }) {
  const marriage = tone === "dark" ? "#d4af37" : "#a16207";
  const descent = tone === "dark" ? "#5eead4" : "#0f766e";
  const label = tone === "dark" ? "text-teal-100/80" : "text-slate-600";

  return (
    <div className={`flex flex-wrap items-center gap-x-4 gap-y-1 text-xs ${label}`}>
      <span className="inline-flex items-center gap-2">
        <svg width="28" height="8" aria-hidden>
          <line x1="1" y1="4" x2="27" y2="4" stroke={marriage} strokeWidth="2" strokeLinecap="round" />
        </svg>
        Husband and wife
      </span>
      <span className="inline-flex items-center gap-2">
        <svg width="28" height="14" aria-hidden>
          <path
            d={`M 14 1 L 14 6 M 4 6 L 24 6 M 4 6 L 4 13 M 14 6 L 14 13 M 24 6 L 24 13`}
            fill="none"
            stroke={descent}
            strokeWidth="1.6"
            strokeLinecap="round"
          />
        </svg>
        Parents and children
      </span>
    </div>
  );
}
