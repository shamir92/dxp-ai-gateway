"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { SiteHeader } from "@/components/public-shell";
import { ButtonLink } from "@/components/button-link";

export default function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    let redirectTimer: ReturnType<typeof setTimeout> | undefined;
    apiFetch("/api/v1/auth/logout", { method: "POST" })
      .catch(() => {})
      .finally(() => {
        redirectTimer = setTimeout(() => router.replace("/login"), 1800);
      });
    return () => {
      if (redirectTimer) clearTimeout(redirectTimer);
    };
  }, [router]);

  return (
    <div className="public-shell flex min-h-screen flex-col">
      <SiteHeader />
      <main className="flex flex-1">
        <section className="flex min-h-[374px] flex-1 items-center justify-center bg-cream py-12">
          <div className="w-full max-w-[420px] rounded-2xl border border-line-warm bg-white px-8 py-7 text-center">
            <h1 className="font-heading text-[22px] font-semibold text-ink">Signing out…</h1>
            <p className="mt-2 text-[12px] text-muted-warm">
              You are being redirected to the login page.
            </p>
            <ButtonLink
              href="/login"
              className="mt-5 h-10 rounded-full bg-ink px-5 text-[12px] text-white hover:bg-ink/90"
            >
              Back to login
            </ButtonLink>
          </div>
        </section>
      </main>
    </div>
  );
}
