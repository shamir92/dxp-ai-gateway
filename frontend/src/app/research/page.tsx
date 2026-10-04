import { PublicShell } from "@/components/public-shell";

export default function ResearchPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-24">
        <div className="mx-auto max-w-[1302px]">
          <h1 className="font-heading text-[42px] font-semibold text-ink">Research</h1>
          <p className="mt-3 max-w-[560px] text-[14px] text-muted-warm">
            Explorations in efficient reasoning, multimodal systems, and agentic workflows.
          </p>
        </div>
      </section>
    </PublicShell>
  );
}
