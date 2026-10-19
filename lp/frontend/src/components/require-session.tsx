"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";

/**
 * Client-side guard for pages that require a session.
 *
 * The gateway authenticates with an httpOnly cookie, so the backend is the only
 * thing that can answer "am I signed in?". This asks it once and redirects to
 * /login otherwise. While the check is in flight nothing is rendered, so a
 * protected page never flashes its content to a signed-out visitor.
 *
 * Children render as-is, so the wrapped pages stay server components.
 */
export function RequireSession({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [allowed, setAllowed] = useState(false);

  useEffect(() => {
    let cancelled = false;
    void (async () => {
      try {
        const res = await apiFetch("/api/v1/auth/me");
        if (cancelled) return;
        if (res.ok) {
          setAllowed(true);
          return;
        }
        router.replace("/login");
      } catch {
        if (!cancelled) router.replace("/login");
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [router]);

  if (!allowed) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[#F7F4EE]">
        <p className="text-[12px] text-muted-warm">Checking your session…</p>
      </div>
    );
  }

  return <>{children}</>;
}