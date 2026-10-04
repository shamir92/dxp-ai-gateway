import { PublicShell } from "@/components/public-shell";

export default function UpdatesPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-24">
        <div className="mx-auto max-w-[1302px]">
          <h1 className="font-heading text-[42px] font-semibold text-ink">Updates</h1>
          <p className="mt-3 max-w-[560px] text-[14px] text-muted-warm">
            Product announcements and release notes for DXP AI Gateway.
          </p>
        </div>
      </section>
    </PublicShell>
  );
}
