import type { NextConfig } from "next";

const API_URL = process.env.API_URL ?? "http://localhost:8080";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/auth/signin/google",
        destination: `${API_URL}/api/v1/auth/oauth/google`,
      },
      {
        source: "/api/auth/signin/github",
        destination: `${API_URL}/api/v1/auth/oauth/github`,
      },
      {
        source: "/api/v1/:path*",
        destination: `${API_URL}/api/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
