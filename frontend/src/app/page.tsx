import Link from "next/link";
import {
  ArrowRight,
  Check,
  Globe,
  Key,
  Play,
  Route,
  SlidersHorizontal,
  Workflow,
  Wrench,
} from "lucide-react";
import { PublicShell } from "@/components/public-shell";
import { Button } from "@/components/ui/button";
import { ButtonLink } from "@/components/button-link";
import { Badge } from "@/components/ui/badge";

const metrics = [
  { value: "38B+", label: "tokens served daily" },
  { value: "99.9%", label: "gateway uptime" },
  { value: "284 ms", label: "median latency" },
  { value: "30+", label: "models and tools" },
];

const models = [
  {
    name: "MiMo V2.6 Pro",
    tag: "FLAGSHIP REASONING",
    stat: "256K context",
    desc: "Deep analysis, complex coding, and agent planning.",
    dark: true,
  },
  {
    name: "MiMo V2.6 Flash",
    tag: "HIGH-SPEED INFERENCE",
    stat: "3x faster",
    desc: "Low-latency chat, extraction, and routing.",
    dark: false,
  },
  {
    name: "MiMo V2.5 Omni",
    tag: "MULTIMODAL",
    stat: "4 modalities",
    desc: "Understand text, image, audio, and video together.",
    dark: false,
  },
];

const benefits = [
  "Unified authentication and usage tracking",
  "Streaming, tools, and structured output",
  "Regional endpoints and batch processing",
];

const capabilities = [
  {
    icon: Route,
    title: "Smart routing",
    desc: "Match each request to the best model for quality, latency, and cost.",
  },
  {
    icon: SlidersHorizontal,
    title: "Usage intelligence",
    desc: "Understand tokens, requests, cache efficiency, and spend by model.",
  },
  {
    icon: Workflow,
    title: "Batch operations",
    desc: "Run high-volume asynchronous workloads with status and retry controls.",
  },
  {
    icon: Wrench,
    title: "Trusted tools",
    desc: "Add web search, code execution, and private knowledge with permissions.",
  },
];

const useCaseTabs = [
  "Developer copilots",
  "Customer operations",
  "Content intelligence",
  "Voice experiences",
];

const codeSample = `from openai import OpenAI

client = OpenAI(
  base_url="https://api.dxp.ai/v1",
  api_key="$DXP_API_KEY"
)

response = client.chat.completions.create(
  model="mimo-v2.6-pro",
  messages=[{"role": "user", "content": "Plan our launch"}],
  stream=True
)

for chunk in response:
  print(chunk.choices[0].delta.content)`;

export default function LandingPage() {
  return (
    <PublicShell>
      {/* Hero */}
      <section className="relative overflow-hidden bg-cream px-14 py-20">
        <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,rgba(255,90,42,0.12),transparent_55%)]" />
        <div className="relative mx-auto grid max-w-[1302px] items-center gap-14 lg:grid-cols-[1fr_1.1fr]">
          <div>
            <Badge className="rounded-full border-line-warm bg-white px-3 py-1.5 text-[10px] font-medium text-ink">
              <span className="mr-1.5 size-1.5 rounded-full bg-flame" />
              MIMO V2.6 IS NOW AVAILABLE
            </Badge>
            <h1 className="mt-6 font-heading text-[58px] font-semibold leading-[1.05] tracking-tight text-ink">
              Build intelligence that moves at your speed.
            </h1>
            <p className="mt-5 max-w-[550px] text-[15px] leading-relaxed text-muted-warm">
              One gateway for high-performance reasoning, multimodal creation, voice,
              and agentic workflows — built for developers shipping real products.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <ButtonLink href="/token-plan" className="h-11 rounded-full bg-ink px-6 text-[13px] text-white hover:bg-ink/90">
                  Start building <ArrowRight className="ml-1.5 size-3.5" />
                </ButtonLink>
              <ButtonLink href="/models" variant="outline"
                className="h-11 rounded-full border-line-warm bg-white px-6 text-[13px] text-ink">
                  <Play className="mr-1.5 size-3.5 text-flame" /> Explore models
                </ButtonLink>
            </div>
            <p className="mt-5 text-[11px] text-muted-warm">
              OpenAI-compatible · Built-in usage analytics · Start free
            </p>
          </div>

          {/* Gateway visual */}
          <div className="relative rounded-[28px] bg-night-card p-10 shadow-2xl">
            <div className="absolute right-8 top-8 text-[10px] text-[#697689]">sgp-1</div>
            <div className="flex items-center gap-2 text-[11px] text-white">
              <span className="size-2 rounded-full bg-success" />
              DXP GATEWAY · LIVE
            </div>
            <div className="mt-6 rounded-2xl bg-prompt p-5">
              <p className="text-[9px] uppercase tracking-wider text-[#697689]">Request</p>
              <p className="mt-2 text-[13px] leading-relaxed text-white">
                Analyze this product brief, identify risks, and propose a launch plan.
              </p>
            </div>
            <div className="mt-4 grid grid-cols-3 gap-3">
              {models.map((m) => (
                <div key={m.name} className="rounded-xl bg-[#FFE8E1] p-3">
                  <p className="text-[10px] font-mono text-ink">{m.name.toLowerCase().replace(/ /g, "-")}</p>
                  <p className="mt-1 text-[9px] text-flame">{m.tag.split(" ")[0].toLowerCase()}</p>
                  <p className="mt-2 text-[9px] font-semibold text-flame">SELECTED</p>
                </div>
              ))}
            </div>
            <div className="mt-4 rounded-2xl bg-white p-5">
              <div className="flex items-center justify-between">
                <p className="text-[9px] uppercase tracking-wider text-muted-warm">Response</p>
                <p className="text-[10px] text-success">284 ms</p>
              </div>
              <div className="mt-3 space-y-2">
                <div className="h-2 w-full rounded bg-[#EDEAE4]" />
                <div className="h-2 w-4/5 rounded bg-[#EDEAE4]" />
                <div className="h-2 w-3/5 rounded bg-[#EDEAE4]" />
              </div>
              <div className="mt-4 flex gap-2 text-[9px] text-muted-warm">
                <span className="rounded-full bg-cream px-2 py-1">4,828 tokens</span>
                <span className="rounded-full bg-cream px-2 py-1">cache 56%</span>
                <span className="rounded-full bg-cream px-2 py-1">$0.012</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Metrics */}
      <section className="bg-white px-14 py-12">
        <div className="mx-auto grid max-w-[1302px] grid-cols-2 gap-8 md:grid-cols-4">
          {metrics.map((m) => (
            <div key={m.label} className="text-center">
              <p className="font-heading text-[30px] font-semibold text-ink">{m.value}</p>
              <p className="mt-1 text-[11px] text-muted-warm">{m.label}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Models */}
      <section className="bg-cream px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <div className="flex flex-wrap items-end justify-between gap-4">
            <div>
              <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
                Model family
              </p>
              <h2 className="mt-3 font-heading text-[36px] font-semibold text-ink">
                One family. Every workload.
              </h2>
              <p className="mt-2 text-[13px] text-muted-warm">
                Choose the right balance of intelligence, speed, and modality.
              </p>
            </div>
            <Link href="/models" className="text-[11px] text-flame">
              Explore all models →
            </Link>
          </div>
          <div className="mt-10 grid gap-5 md:grid-cols-3">
            {models.map((m) => (
              <div
                key={m.name}
                className={`flex min-h-[280px] flex-col rounded-2xl p-6 ${
                  m.dark ? "bg-night-card text-white" : "border border-line-warm bg-white"
                }`}
              >
                <div className="flex items-start justify-between">
                  <div
                    className={`flex size-10 items-center justify-center rounded-xl ${
                      m.dark ? "bg-flame/20" : "bg-cream"
                    }`}
                  >
                    <Globe className={`size-5 ${m.dark ? "text-flame-soft" : "text-flame"}`} />
                  </div>
                  <span className={`text-[10px] ${m.dark ? "text-[#697689]" : "text-muted-warm"}`}>
                    {m.stat}
                  </span>
                </div>
                <p className={`mt-6 text-[10px] font-semibold ${m.dark ? "text-flame-soft" : "text-flame"}`}>
                  {m.tag}
                </p>
                <h3 className={`mt-2 text-[22px] font-semibold ${m.dark ? "text-white" : "text-ink"}`}>
                  {m.name}
                </h3>
                <p className={`mt-2 text-[12px] leading-relaxed ${m.dark ? "text-[#91A0B6]" : "text-muted-warm"}`}>
                  {m.desc}
                </p>
                <Link
                  href="#"
                  className={`mt-auto pt-8 text-[11px] ${m.dark ? "text-flame-soft" : "text-flame"}`}
                >
                  View model details →
                </Link>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Developer */}
      <section className="bg-night px-14 py-20 text-white">
        <div className="mx-auto grid max-w-[1302px] items-center gap-12 lg:grid-cols-[1fr_1.3fr]">
          <div>
            <p className="text-[10px] font-semibold uppercase tracking-wider text-flame-soft">
              Built for developers
            </p>
            <h2 className="mt-4 font-heading text-[36px] font-semibold leading-tight">
              One endpoint. Every capability.
            </h2>
            <p className="mt-4 max-w-[440px] text-[13px] leading-relaxed text-[#AAB5C5]">
              Use familiar OpenAI-compatible APIs, then route across reasoning,
              multimodal, speech, and batch models without rewriting your stack.
            </p>
            <ul className="mt-8 space-y-4">
              {benefits.map((b) => (
                <li key={b} className="flex items-center gap-3">
                  <span className="flex size-8 items-center justify-center rounded-lg bg-prompt">
                    <Key className="size-3.5 text-success" />
                  </span>
                  <span className="text-[12px] text-white">{b}</span>
                </li>
              ))}
            </ul>
            <ButtonLink href="/documentation" className="mt-8 h-10 rounded-full bg-white px-5 text-[12px] text-ink hover:bg-white/90">
                Read the documentation <ArrowRight className="ml-1.5 size-3.5" />
              </ButtonLink>
          </div>
          <div className="overflow-hidden rounded-2xl bg-ink">
            <div className="flex items-center gap-2 bg-[#172033] px-4 py-3">
              <span className="size-2.5 rounded-full bg-[#FF5F57]" />
              <span className="size-2.5 rounded-full bg-[#FEBC2E]" />
              <span className="size-2.5 rounded-full bg-[#28C840]" />
              <span className="ml-auto text-[10px] text-[#697689]">quickstart</span>
            </div>
            <pre className="overflow-x-auto p-6 font-mono text-[11px] leading-relaxed text-[#AAB5C5]">
              <code>{codeSample}</code>
            </pre>
          </div>
        </div>
      </section>

      {/* Capabilities */}
      <section className="bg-white px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <div className="flex flex-wrap items-end justify-between gap-4">
            <div>
              <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
                Platform capabilities
              </p>
              <h2 className="mt-3 font-heading text-[36px] font-semibold text-ink">
                From first request to production scale.
              </h2>
            </div>
            <p className="text-[12px] text-muted-warm">Everything managed in one console.</p>
          </div>
          <div className="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {capabilities.map((c) => (
              <div key={c.title} className="flex min-h-[210px] flex-col rounded-2xl bg-cream-card p-6">
                <div className="flex size-10 items-center justify-center rounded-xl bg-white">
                  <c.icon className="size-4 text-flame" />
                </div>
                <h3 className="mt-5 text-[15px] font-semibold text-ink">{c.title}</h3>
                <p className="mt-2 text-[12px] leading-relaxed text-muted-warm">{c.desc}</p>
                <Link href="#" className="mt-auto pt-6 text-[11px] text-flame">
                  Learn more →
                </Link>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Use cases */}
      <section className="bg-cream-deep px-14 py-20">
        <div className="mx-auto grid max-w-[1302px] items-start gap-12 lg:grid-cols-[410px_1fr]">
          <div>
            <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
              Made for real work
            </p>
            <h2 className="mt-3 font-heading text-[34px] font-semibold leading-tight text-ink">
              A platform for every AI product.
            </h2>
            <p className="mt-3 text-[13px] leading-relaxed text-muted-warm">
              From customer experiences to internal automation, DXP AI Gateway gives
              teams the models and controls to move from idea to production.
            </p>
            <div className="mt-8 space-y-2">
              {useCaseTabs.map((tab, i) => (
                <div
                  key={tab}
                  className={`rounded-xl px-4 py-3 text-[12px] ${
                    i === 0
                      ? "border-l-2 border-flame bg-white font-medium text-ink"
                      : "text-muted-warm"
                  }`}
                >
                  {tab}
                </div>
              ))}
            </div>
          </div>
          <div className="rounded-2xl bg-white p-8">
            <div className="flex items-center justify-between">
              <h3 className="text-[15px] font-semibold text-ink">Developer copilot workflow</h3>
              <span className="flex items-center gap-1.5 rounded-full bg-success/10 px-2.5 py-1 text-[10px] text-success">
                <span className="size-1.5 rounded-full bg-success" /> RUNNING
              </span>
            </div>
            <div className="mt-6 grid grid-cols-3 gap-4">
              {["Issue context", "MiMo V2.6 Pro", "Code patch"].map((step) => (
                <div key={step} className="rounded-xl bg-cream-card p-5 text-center">
                  <Workflow className="mx-auto size-5 text-flame" />
                  <p className="mt-3 text-[12px] font-medium text-ink">{step}</p>
                </div>
              ))}
            </div>
            <div className="mt-4 rounded-xl bg-ink p-5 font-mono text-[11px] leading-relaxed text-[#AAB5C5]">
              <p>✓ Context indexed · 24 files</p>
              <p className="mt-1">✓ Plan generated · 6 steps</p>
              <p className="mt-1 text-flame-soft">→ Writing tests and implementation…</p>
            </div>
          </div>
        </div>
      </section>

      {/* Final CTA */}
      <section className="bg-flame px-14 py-16">
        <div className="mx-auto flex max-w-[1302px] flex-wrap items-center justify-between gap-8">
          <div>
            <h2 className="font-heading text-[28px] font-semibold text-white">
              Turn your next idea into production AI.
            </h2>
            <p className="mt-2 text-[13px] text-[#FFE4DD]">
              Start free, explore the models, and scale when you are ready.
            </p>
          </div>
          <div className="flex gap-3">
            <ButtonLink href="/token-plan" className="h-11 rounded-full bg-ink px-6 text-[13px] text-white hover:bg-ink/90">
                Start building <ArrowRight className="ml-1.5 size-3.5" />
              </ButtonLink>
            <ButtonLink href="/pricing" variant="outline"
              className="h-11 rounded-full border-white bg-white px-6 text-[13px] text-ink">View pricing</ButtonLink>
          </div>
        </div>
      </section>
    </PublicShell>
  );
}
