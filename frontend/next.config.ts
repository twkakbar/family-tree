import type { NextConfig } from "next";

/**
 * In development the Next server proxies API and upload paths to the Go
 * backend, so the browser always talks to a single origin and the refresh
 * cookie stays first-party (ADR-0005). In production nginx does this routing
 * and these rewrites never match.
 */
const apiProxyTarget = process.env.API_PROXY_TARGET ?? "http://localhost:8080";

const nextConfig: NextConfig = {
  output: "standalone",
  async rewrites() {
    return [
      { source: "/api/:path*", destination: `${apiProxyTarget}/api/:path*` },
      { source: "/uploads/:path*", destination: `${apiProxyTarget}/uploads/:path*` },
    ];
  },
};

export default nextConfig;
