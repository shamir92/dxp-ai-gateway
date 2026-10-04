import Link from "next/link";
import { ChevronDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ButtonLink } from "@/components/button-link";
import { Brand } from "@/components/brand";

const navLinks = [
  { label: "Research", href: "/research" },
  { label: "Models", href: "/models", caret: true },
  { label: "Updates", href: "/updates", caret: true },
  { label: "Documentation", href: "/documentation" },
  { label: "Pricing", href: "/pricing", caret: true },
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
            {navLinks.map((link) => (
              <Link
                key={link.label}
                href={link.href}
                className="flex items-center gap-1 text-[13px] text-muted-warm transition-colors hover:text-ink"
              >
                {link.label}
                {link.caret ? <ChevronDown className="size-3" /> : null}
              </Link>
            ))}
          </nav>
          <div className="flex items-center gap-2">
            <ButtonLink href="/login" variant="outline"
              className="h-9 rounded-full border-line-warm bg-white px-4 text-[12px] text-ink">Log in</ButtonLink>
            <ButtonLink href="/token-plan" className="h-9 rounded-full bg-ink px-4 text-[12px] text-white hover:bg-ink/90">Start building</ButtonLink>
          </div>
        </div>
      </header>
    </>
  );
}

const footerColumns = [
  {
    title: "Product",
    links: ["Models", "Token plans", "Batch jobs", "Plugins"],
  },
  {
    title: "Developers",
    links: ["Documentation", "API reference", "SDKs", "Status"],
  },
  {
    title: "Company",
    links: ["Research", "Updates", "Contact", "Careers"],
  },
  {
    title: "Legal",
    links: ["Privacy", "Terms", "Security", "Compliance"],
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
          <div className="grid grid-cols-2 gap-8 sm:grid-cols-4">
            {footerColumns.map((col) => (
              <div key={col.title}>
                <h4 className="text-[11px] font-semibold text-white">{col.title}</h4>
                <ul className="mt-4 space-y-2.5">
                  {col.links.map((link) => (
                    <li key={link}>
                      <Link
                        href="#"
                        className="text-[11px] text-[#697689] transition-colors hover:text-white"
                      >
                        {link}
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
