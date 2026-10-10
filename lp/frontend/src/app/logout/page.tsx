"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";

export default function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    apiFetch("/api/v1/auth/logout", { method: "POST" })
      .catch(() => {})
      .finally(() => {
        router.replace("/login");
      });
  }, [router]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-cream text-[13px] text-muted-warm">
      Signing you out…
    </div>
  );
}
