import type { Metadata } from "next";
import type { ReactNode } from "react";

import { env } from "@/lib/env";
import { Providers } from "./providers";
import "./globals.css";

export const metadata: Metadata = {
  title: {
    default: env.appName,
    template: `%s · ${env.appName}`,
  },
  description: "Create, manage and explore your family tree.",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-dvh font-sans">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
