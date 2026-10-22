import type { NextConfig } from "next";

const backendUrl = process.env.BACKEND_URL ?? "http://localhost:8080";

const nextConfig: NextConfig = {
  rewrites() {
    // Production has vercel.json routing /api/* and /healthz to the Go service
    // at the edge, before Next.js sees the request. This only stands in for it
    // locally, so it must not apply to a prod build.
    if (process.env.NODE_ENV !== "development") return [];
    return [
      { source: "/api/:path*", destination: `${backendUrl}/api/:path*` },
      { source: "/healthz", destination: `${backendUrl}/healthz` },
    ];
  },
};

export default nextConfig;
