"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { BrandMark } from "@/components/brand";

const sections = [
  {
    label: "WORKSPACE",
    items: [
      { label: "Token plan", href: "/token-plan" },
      // Not released yet: shown so the surface is discoverable, but inert — the
      // page behind it is still under construction.
      { label: "API plan", href: "/api-plan", disabled: true },
      // Hidden until implemented: the Plugins page and /plugins route still exist.
    ],
  },
  {
    label: "RESOURCES",
    items: [
      // Hidden until implemented: the Documentation page and /documentation route
      // still exist.
      { label: "Support", href: "/support" },
    ],
  },
  {
    label: "ACCOUNT",
    items: [
      { label: "Profile", href: "/profile" },
      { label: "Log out", href: "/logout" },
    ],
  },
];

export function AppSidebar() {
  const pathname = usePathname();

  return (
    <aside className="hidden w-[196px] shrink-0 flex-col bg-night-card p-5 text-white md:flex">
      <div className="flex items-center">
        <BrandMark inverted />
      </div>

      <nav className="mt-8 flex flex-1 flex-col gap-7">
        {sections.map((section) => (
          <div
            key={section.label}
            className={section.label === "ACCOUNT" ? "border-t border-white/10 pt-5" : undefined}
          >
            <p className="px-3 text-[10px] font-medium tracking-wider text-[#7E8A9E]">
              {section.label}
            </p>
            <div className="mt-2 flex flex-col gap-0.5">
              {section.items.map((item) => {
                const active =
                  pathname === item.href ||
                  pathname.startsWith(item.href + "/") ||
                  (item.href === "/token-plan" && pathname === "/consumption");
                if (item.disabled) {
                  return (
                    <span
                      key={item.href}
                      aria-disabled="true"
                      className="flex cursor-not-allowed items-center justify-between rounded-lg px-3 py-2.5 text-[13px] text-[#AEB8C8]/50 select-none"
                    >
                      {item.label}
                      <span className="rounded-full bg-white/10 px-2 py-0.5 text-[9px] font-medium tracking-wide text-[#AEB8C8] uppercase">
                        Coming soon
                      </span>
                    </span>
                  );
                }
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
    </aside>
  );
}

function AppMobileNav() {
  const pathname = usePathname();
  const items = sections.flatMap((section) => section.items);

  return (
    <nav aria-label="Workspace navigation" className="flex gap-1 overflow-x-auto border-b border-line-warm bg-night-card px-3 py-2 md:hidden">
      {items.map((item) => {
        const active = pathname === item.href || (item.href === "/token-plan" && pathname === "/consumption");
        if (item.disabled) {
          return (
            <span
              key={item.href}
              aria-disabled="true"
              className="shrink-0 cursor-not-allowed rounded-full px-3 py-2 text-[11px] text-[#AEB8C8]/50 select-none"
            >
              {item.label} · coming soon
            </span>
          );
        }
        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "shrink-0 rounded-full px-3 py-2 text-[11px] transition-colors",
              active ? "bg-white/[0.1] text-white" : "text-[#AEB8C8] hover:text-white",
            )}
          >
            {item.label}
          </Link>
        );
      })}
    </nav>
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
    <div className="app-shell flex min-h-screen flex-col bg-[#F7F4EE] md:flex-row">
      <AppSidebar />
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex min-h-[72px] flex-wrap items-center justify-between gap-3 border-b border-line-warm bg-white px-4 py-3 md:flex-nowrap md:gap-4 md:px-8 md:py-4">
          <div>
            <h1 className="font-heading text-[22px] font-semibold text-ink">{title}</h1>
            {subtitle ? (
              <p className="mt-0.5 text-[12px] text-muted-warm">{subtitle}</p>
            ) : null}
          </div>
          {actions ? <div className="flex w-full flex-wrap items-center gap-2 md:w-auto">{actions}</div> : null}
        </header>
        <AppMobileNav />
        <main className="flex-1 overflow-auto p-4 md:p-8">{children}</main>
      </div>
    </div>
  );
}
