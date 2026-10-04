import { PublicShell } from "@/components/public-shell";

const models = [
  { name: "MiMo V2.6 Pro", tag: "Reasoning", desc: "Deep analysis, complex coding, and agent planning." },
  { name: "MiMo V2.6 Flash", tag: "Speed", desc: "Low-latency chat, extraction, and routing." },
  { name: "MiMo V2.5 Omni", tag: "Multimodal", desc: "Text, image, audio, and video in one model." },
];

export default function ModelsPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">Model family</p>
          <h1 className="mt-3 font-heading text-[42px] font-semibold text-ink">One family. Every workload.</h1>
          <p className="mt-3 max-w-[560px] text-[14px] text-muted-warm">
            Choose the right balance of intelligence, speed, and modality.
          </p>
          <div className="mt-10 grid gap-5 md:grid-cols-3">
            {models.map((m) => (
              <div key={m.name} className="rounded-2xl border border-line-warm bg-white p-6">
                <p className="text-[10px] font-semibold text-flame">{m.tag}</p>
                <h2 className="mt-2 text-[22px] font-semibold text-ink">{m.name}</h2>
                <p className="mt-2 text-[13px] text-muted-warm">{m.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>
    </PublicShell>
  );
}
