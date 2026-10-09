"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Gift } from "lucide-react";
import { cn } from "@/lib/utils";
import { BrandMark } from "@/components/brand";

const sections = [
  {
    label: "WORKSPACE",
    items: [
      { label: "Token plan", href: "/token-plan" },
      { label: "API lifecycle", href: "/api-lifecycle" },
      { label: "Plugins", href: "/plugins" },
    ],
  },
  {
    label: "ACCOUNT",
    items: [
      { label: "Profile", href: "/profile" },
      { label: "Settings", href: "/settings" },
    ],
  },
  {
    label: "RESOURCES",
    items: [
      { label: "Documentation", href: "/documentation" },
      { label: "Support", href: "/support" },
    ],
  },
];

export function AppSidebar() {
  const pathname = usePathname();

  return (
    <aside className="flex w-[236px] shrink-0 flex-col bg-night-card p-5 text-white">
      <div className="flex items-center">
        <BrandMark inverted className="h-8 w-auto" />
      </div>

      <nav className="mt-8 flex flex-1 flex-col gap-7">
        {sections.map((section) => (
          <div key={section.label}>
            <p className="px-3 text-[10px] font-medium tracking-wider text-[#7E8A9E]">
              {section.label}
            </p>
            <div className="mt-2 flex flex-col gap-0.5">
              {section.items.map((item) => {
                const active =
                  pathname === item.href || pathname.startsWith(item.href + "/");
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={cn(
                      "rounded-lg px-3 py-2.5 text-[13px] transition-colors",
                      active
                        ? "bg-white/[0.07] text-white"
                        : "text-[#AEB8C8] hover:bg-white/5 hover:text-white",
                    )}
                  >
                    {item.label}
                  </Link>
                );
              })}
            </div>
          </div>
        ))}
      </nav>

      <div className="mt-auto rounded-xl bg-prompt p-3.5">
        <Gift className="size-5 text-[#FFD29A]" />
        <p className="mt-2.5 text-[12px] font-medium text-white">Build with your team</p>
        <p className="mt-1 text-[10px] leading-relaxed text-[#91A0B6]">
          Invite collaborators and share one workspace.
        </p>
        <button className="mt-3 w-full rounded-lg bg-white py-2 text-[11px] font-medium text-night-card">
          Invite members
        </button>
      </div>
    </aside>
  );
}

export function AppShell({
  children,
  title,
  subtitle,
  actions,
}: {
  children: React.ReactNode;
  title: string;
  subtitle?: string;
  actions?: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen bg-[#F7F4EE]">
      <AppSidebar />
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex min-h-[72px] items-center justify-between gap-4 border-b border-line-warm bg-white px-8 py-4">
          <div>
            <h1 className="font-heading text-[22px] font-semibold text-ink">{title}</h1>
            {subtitle ? (
              <p className="mt-0.5 text-[12px] text-muted-warm">{subtitle}</p>
            ) : null}
          </div>
          {actions ? <div className="flex flex-wrap items-center gap-2">{actions}</div> : null}
        </header>
        <main className="flex-1 overflow-auto p-8">{children}</main>
      </div>
    </div>
  );
}
