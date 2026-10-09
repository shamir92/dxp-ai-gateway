import { PublicShell } from "@/components/public-shell";

const models = [
  {
    id: "deepseek-v4.1-flash",
    vendor: "DeepSeek",
    tag: "Reasoning",
    desc: "Fast, efficient reasoning for production workloads.",
  },
  {
    id: "deepseek-v4-flash",
    vendor: "DeepSeek",
    tag: "Speed",
    desc: "Low-latency chat, extraction, and routing.",
  },
  {
    id: "mimo-v2.6-pro",
    vendor: "MiMo",
    tag: "Flagship",
    desc: "Deep analysis, complex coding, and agent planning.",
  },
  {
    id: "mimo-v2.6-flash",
    vendor: "MiMo",
    tag: "Speed",
    desc: "High-speed general purpose inference.",
  },
  {
    id: "glm-5.3-flash",
    vendor: "GLM",
    tag: "Balanced",
    desc: "Strong balance of speed and quality.",
  },
  {
    id: "glm-5.3",
    vendor: "GLM",
    tag: "Reasoning",
    desc: "Strong reasoning and coding performance.",
  },
  {
    id: "kimi-k3",
    vendor: "Kimi",
    tag: "Long context",
    desc: "Built for long-context workflows.",
  },
  {
    id: "minimax-m3",
    vendor: "MiniMax",
    tag: "Multimodal",
    desc: "Multimodal understanding across inputs.",
  },
];

export default function ModelsPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
            Available models
          </p>
          <h1 className="mt-3 font-heading text-[42px] font-semibold text-ink">
            Models on the gateway.
          </h1>
          <p className="mt-3 max-w-[560px] text-[14px] text-muted-warm">
            One API key, eight production models — and more added regularly.
          </p>
          <div className="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {models.map((m) => (
              <div key={m.id} className="rounded-2xl border border-line-warm bg-white p-6">
                <div className="flex items-start justify-between gap-2">
                  <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
                    {m.vendor}
                  </p>
                  <span className="rounded-full bg-cream px-2 py-0.5 text-[9px] font-medium text-muted-warm">
                    {m.tag}
                  </span>
                </div>
                <h2 className="mt-3 font-mono text-[14px] font-medium text-ink">{m.id}</h2>
                <p className="mt-2 text-[12px] leading-relaxed text-muted-warm">{m.desc}</p>
              </div>
            ))}
          </div>
          <p className="mt-8 text-[13px] text-muted-warm">
            Model list continues to expand. Use the same OpenAI-compatible API for every model.
          </p>
        </div>
      </section>
    </PublicShell>
  );
}
