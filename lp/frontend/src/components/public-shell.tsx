import Link from "next/link";
import { ChevronDown } from "lucide-react";
import { Brand } from "@/components/brand";
import { ButtonLink } from "@/components/button-link";

type NavItem = {
  label: string;
  href?: string;
  children?: { label: string; href: string }[];
};

const navItems: NavItem[] = [
  { label: "Models", href: "/models" },
  {
    label: "Pricing",
    children: [
      { label: "Token plans", href: "/pricing" },
      { label: "API pricing", href: "/api-pricing" },
    ],
  },
  { label: "Contact", href: "/contact" },
  { label: "Privacy", href: "/privacy" },
  { label: "Terms", href: "/terms" },
];

export function SiteHeader() {
  return (
    <>
      <div className="flex h-[34px] items-center justify-center bg-night px-4">
        <p className="text-[10px] font-medium text-white">
          MiMo V2.6 series is now live — faster reasoning, full modality, built for
          production →
        </p>
      </div>
      <header className="border-b border-line-warm bg-white">
        <div className="mx-auto flex h-[72px] max-w-[1442px] items-center justify-between px-14">
          <Brand />
          <nav className="hidden items-center gap-8 md:flex">
            {navItems.map((item) =>
              item.children ? (
                <div key={item.label} className="group relative">
                  <button
                    type="button"
                    className="flex items-center gap-1 text-[13px] text-muted-warm transition-colors hover:text-ink"
                  >
                    {item.label}
                    <ChevronDown className="size-3" />
                  </button>
                  <div className="invisible absolute left-1/2 top-full z-50 w-40 -translate-x-1/2 pt-2 opacity-0 transition-all group-hover:visible group-hover:opacity-100 group-focus-within:visible group-focus-within:opacity-100">
                    <div className="rounded-xl border border-line-warm bg-white py-2 shadow-lg">
                      {item.children.map((child) => (
                        <Link
                          key={child.href}
                          href={child.href}
                          className="block px-4 py-2 text-[12px] text-muted-warm transition-colors hover:bg-cream-card hover:text-ink"
                        >
                          {child.label}
                        </Link>
                      ))}
                    </div>
                  </div>
                </div>
              ) : (
                <Link
                  key={item.label}
                  href={item.href!}
                  className="text-[13px] text-muted-warm transition-colors hover:text-ink"
                >
                  {item.label}
                </Link>
              ),
            )}
          </nav>
          <div className="flex items-center gap-2">
            <ButtonLink href="/login" variant="outline"
              className="h-9 rounded-full border-line-warm bg-white px-4 text-[12px] text-ink">Log in</ButtonLink>
            <ButtonLink href="/register" className="h-9 rounded-full bg-ink px-4 text-[12px] text-white hover:bg-ink/90">Start building</ButtonLink>
          </div>
        </div>
      </header>
    </>
  );
}

const footerColumns = [
  {
    title: "Product",
    links: [{ label: "Models", href: "/models" }],
  },
  {
    title: "Company",
    links: [{ label: "Contact", href: "/contact" }],
  },
  {
    title: "Legal",
    links: [
      { label: "Privacy", href: "/privacy" },
      { label: "Terms", href: "/terms" },
    ],
  },
];

export function SiteFooter() {
  return (
    <footer className="bg-night px-14 py-16 text-white">
      <div className="mx-auto max-w-[1330px]">
        <div className="grid gap-12 lg:grid-cols-[360px_1fr]">
          <div>
            <Brand inverted />
            <p className="mt-5 max-w-[360px] text-[12px] leading-relaxed text-[#697689]">
              Production-ready AI infrastructure for developers and teams building
              the next generation of intelligent products.
            </p>
          </div>
          <div className="grid grid-cols-2 gap-8 sm:grid-cols-3">
            {footerColumns.map((col) => (
              <div key={col.title}>
                <h4 className="text-[11px] font-semibold text-white">{col.title}</h4>
                <ul className="mt-4 space-y-2.5">
                  {col.links.map((link) => (
                    <li key={link.label}>
                      <Link
                        href={link.href}
                        className="text-[11px] text-[#697689] transition-colors hover:text-white"
                      >
                        {link.label}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>
        <div className="mt-12 flex items-center justify-between border-t border-white/10 pt-8">
          <p className="text-[10px] text-[#697689]">
            © 2026 DXP AI Gateway. All rights reserved.
          </p>
          <div className="flex gap-3">
            {["X", "GH", "LI"].map((s) => (
              <span
                key={s}
                className="flex size-6 items-center justify-center rounded-full bg-white/10 text-[8px] text-white/70"
              >
                {s}
              </span>
            ))}
          </div>
        </div>
      </div>
    </footer>
  );
}

export function PublicShell({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader />
      <main className="flex-1">{children}</main>
      <SiteFooter />
    </div>
  );
}
