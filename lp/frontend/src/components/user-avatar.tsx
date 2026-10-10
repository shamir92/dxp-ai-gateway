"use client";

import { useEffect, useState } from "react";
import { apiJSON } from "@/lib/api";

type Me = {
  user?: {
    id: string;
    email: string;
    name: string;
    avatar_url: string;
  };
};

function initials(name: string, email: string) {
  const src = (name || email || "?").trim();
  const parts = src.split(/[\s@._-]+/).filter(Boolean);
  if (parts.length === 0) return "?";
  if (parts.length === 1) return parts[0]!.slice(0, 2).toUpperCase();
  return (parts[0]![0]! + parts[1]![0]!).toUpperCase();
}

export function UserAvatar({ className = "" }: { className?: string }) {
  const [label, setLabel] = useState("··");
  const [avatarUrl, setAvatarUrl] = useState("");

  useEffect(() => {
    let cancelled = false;
    apiJSON<Me>("/api/v1/auth/me")
      .then((data) => {
        if (cancelled) return;
        const u = data.user ?? data;
        const anyU = u as { name?: string; email?: string; avatar_url?: string };
        if (anyU?.avatar_url) setAvatarUrl(anyU.avatar_url);
        setLabel(initials(anyU?.name ?? "", anyU?.email ?? ""));
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div
      className={`flex size-9 items-center justify-center overflow-hidden rounded-full bg-prompt text-[11px] font-semibold text-white ${className}`}
      title="Account"
    >
      {avatarUrl ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={avatarUrl} alt="" className="size-full object-cover" />
      ) : (
        label
      )}
    </div>
  );
}
