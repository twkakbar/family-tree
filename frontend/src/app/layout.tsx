import type { Metadata } from "next";
import type { ReactNode } from "react";
import { Cormorant_Garamond, Source_Sans_3 } from "next/font/google";

import { Providers } from "./providers";
import "./globals.css";

const sourceSans = Source_Sans_3({
  subsets: ["latin"],
  variable: "--font-source-sans",
  display: "swap",
});

const cormorant = Cormorant_Garamond({
  subsets: ["latin"],
  weight: ["500", "600", "700"],
  variable: "--font-cormorant",
  display: "swap",
});

export const metadata: Metadata = {
  title: {
    default: "Kesultanan Aceh",
    template: `%s · Kesultanan Aceh`,
  },
  description: "Public family tree of the Aceh royal lineage.",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className={`${sourceSans.variable} ${cormorant.variable}`}>
      <body className="min-h-dvh font-sans">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
