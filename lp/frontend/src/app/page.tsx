import Link from "next/link";
import {
  ArrowRight,
  Check,
  Gauge,
  Key,
  Monitor,
  Play,
  Route,
  ShieldCheck,
  SlidersHorizontal,
  Terminal,
  Wallet,
  Workflow,
  Wrench,
  Zap,
} from "lucide-react";
import { PublicShell } from "@/components/public-shell";
import { Button } from "@/components/ui/button";
import { ButtonLink } from "@/components/button-link";
import { Badge } from "@/components/ui/badge";
import { ApiPlayground } from "@/components/api-playground";

const metrics = [
  { value: "OpenAI", label: "compatible API" },
  { value: "8+", label: "models and growing" },
  { value: "Streaming", label: "tools and batch jobs" },
  { value: "Analytics", label: "usage and spend tracking" },
];

const models = [
  { id: "deepseek-v4.1-flash", vendor: "DeepSeek", note: "Fast, efficient reasoning" },
  { id: "deepseek-v4-flash", vendor: "DeepSeek", note: "Low-latency chat and extraction" },
  { id: "mimo-v2.6-pro", vendor: "MiMo", note: "Deep analysis and agent planning" },
  { id: "mimo-v2.6-flash", vendor: "MiMo", note: "High-speed general purpose" },
  { id: "glm-5.3-flash", vendor: "GLM", note: "Balanced speed and quality" },
  { id: "glm-5.3", vendor: "GLM", note: "Strong reasoning and coding" },
  { id: "kimi-k3", vendor: "Kimi", note: "Long-context workflows" },
  { id: "minimax-m3", vendor: "MiniMax", note: "Multimodal understanding" },
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
    desc: "Match each request to the right model for quality, latency, and cost.",
  },
  {
    icon: Zap,
    title: "Streaming & tools",
    desc: "SSE streaming, function calling, and structured JSON output built in.",
  },
  {
    icon: Gauge,
    title: "Rate limits & quotas",
    desc: "Per-key RPS, budgets, and concurrency controls for every environment.",
  },
  {
    icon: SlidersHorizontal,
    title: "Usage intelligence",
    desc: "Tokens, requests, cache efficiency, and spend by model and key.",
  },
  {
    icon: Wallet,
    title: "Credits & billing",
    desc: "Metered API billing and credit packs for tools — clear cost per call.",
  },
  {
    icon: ShieldCheck,
    title: "Secure by default",
    desc: "Scoped API keys, rotation, audit trails, and regional endpoints.",
  },
  {
    icon: Workflow,
    title: "Batch operations",
    desc: "High-volume async workloads with status, retries, and priority queues.",
  },
  {
    icon: Wrench,
    title: "Fallback & reliability",
    desc: "Retries, timeouts, and healthy fallback paths to keep products online.",
  },
];

const officialTools = [
  {
    name: "Claude Code",
    desc: "CLI and extension — point Anthropic-compatible settings at DXP.",
  },
  {
    name: "Codex",
    desc: "Coding agent — set Base URL and your key, then keep working.",
  },
  {
    name: "OpenClaw",
    desc: "Coding agent and gateway with full tool-calling support.",
  },
  {
    name: "Hermes",
    desc: "Add a custom provider in a few steps.",
  },
];

const supportedEditors = ["VS Code", "Cursor", "Antigravity", "Claude Desktop"];

const codeSample = `from openai import OpenAI

client = OpenAI(
  base_url="...",
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
              <ButtonLink href="/register" className="h-11 rounded-full bg-ink px-6 text-[13px] text-white hover:bg-ink/90">
                  Start building <ArrowRight className="ml-1.5 size-3.5" />
                </ButtonLink>
              <ButtonLink href="/models" variant="outline"
                className="h-11 rounded-full border-line-warm bg-white px-6 text-[13px] text-ink">
                  <Play className="mr-1.5 size-3.5 text-flame" /> Explore models
                </ButtonLink>
            </div>
            <p className="mt-5 text-[11px] text-muted-warm">
              OpenAI-compatible · Built-in usage analytics · 8+ models
            </p>
          </div>

          {/* API playground visual */}
          <ApiPlayground />
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
                Available models
              </p>
              <h2 className="mt-3 font-heading text-[36px] font-semibold text-ink">
                Models on the gateway.
              </h2>
              <p className="mt-2 text-[13px] text-muted-warm">
                One API key, eight production models — and more added regularly.
              </p>
            </div>
            <Link href="/models" className="text-[11px] text-flame">
              Explore all models →
            </Link>
          </div>
          <div className="mt-10 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            {models.map((m, i) => (
              <div
                key={m.id}
                className={`flex min-h-[140px] flex-col rounded-2xl p-5 ${
                  i === 0
                    ? "bg-night-card text-white"
                    : "border border-line-warm bg-white"
                }`}
              >
                <div className="flex items-start justify-between gap-2">
                  <span
                    className={`text-[10px] font-semibold uppercase tracking-wider ${
                      i === 0 ? "text-flame-soft" : "text-flame"
                    }`}
                  >
                    {m.vendor}
                  </span>
                </div>
                <p
                  className={`mt-3 font-mono text-[13px] ${
                    i === 0 ? "text-white" : "text-ink"
                  }`}
                >
                  {m.id}
                </p>
                <p
                  className={`mt-1.5 text-[11px] leading-relaxed ${
                    i === 0 ? "text-[#91A0B6]" : "text-muted-warm"
                  }`}
                >
                  {m.note}
                </p>
              </div>
            ))}
          </div>
          <p className="mt-5 text-[12px] text-muted-warm">
            Model list continues to expand — check the models page for the latest.
          </p>
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
            <ButtonLink href="/login" className="mt-8 h-10 rounded-full bg-white px-5 text-[12px] text-ink hover:bg-white/90">
                Sign in to dashboard <ArrowRight className="ml-1.5 size-3.5" />
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
              <div key={c.title} className="flex min-h-[190px] flex-col rounded-2xl bg-cream-card p-6">
                <div className="flex size-10 items-center justify-center rounded-xl bg-white">
                  <c.icon className="size-4 text-flame" />
                </div>
                <h3 className="mt-5 text-[15px] font-semibold text-ink">{c.title}</h3>
                <p className="mt-2 text-[12px] leading-relaxed text-muted-warm">{c.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Compatibility */}
      <section className="bg-cream-deep px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
            Compatibility
          </p>
          <h2 className="mt-3 font-heading text-[36px] font-semibold leading-tight text-ink">
            Works with the tools you already use.
          </h2>
          <p className="mt-3 max-w-[560px] text-[13px] leading-relaxed text-muted-warm">
            Use DXP credits and models from coding agents, editors, and apps —
            OpenAI and Anthropic-compatible, with one API key.
          </p>

          <div className="mt-12 grid gap-10 lg:grid-cols-[1.15fr_0.85fr]">
            <div>
              <p className="text-[10px] font-semibold tracking-wider text-muted-warm">
                OFFICIAL TOOLS
              </p>
              <div className="mt-5 grid gap-4 sm:grid-cols-2">
                {officialTools.map((tool) => (
                  <div
                    key={tool.name}
                    className="rounded-2xl border border-line-warm bg-white p-6"
                  >
                    <div className="flex size-10 items-center justify-center rounded-xl bg-cream">
                      <Terminal className="size-4 text-flame" />
                    </div>
                    <h3 className="mt-4 text-[16px] font-semibold text-ink">{tool.name}</h3>
                    <p className="mt-2 text-[12px] leading-relaxed text-muted-warm">{tool.desc}</p>
                  </div>
                ))}
              </div>
            </div>

            <div>
              <p className="text-[10px] font-semibold tracking-wider text-muted-warm">
                SUPPORTED EDITORS
              </p>
              <div className="mt-5 grid gap-3 sm:grid-cols-2 lg:grid-cols-1">
                {supportedEditors.map((editor) => (
                  <div
                    key={editor}
                    className="flex items-center gap-3 rounded-2xl border border-line-warm bg-white px-5 py-4"
                  >
                    <span className="flex size-9 items-center justify-center rounded-full bg-cream">
                      <Monitor className="size-4 text-flame" />
                    </span>
                    <span className="text-[14px] font-medium text-ink">{editor}</span>
                  </div>
                ))}
              </div>
              <p className="mt-6 text-[12px] leading-relaxed text-muted-warm">
                Compatible with OpenAI and Anthropic — any tool that can set a Base
                URL and API key works: agents, automation, and your own SDK.
              </p>
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
              Explore the models and scale when you are ready.
            </p>
          </div>
          <div className="flex gap-3">
            <ButtonLink href="/register" className="h-11 rounded-full bg-ink px-6 text-[13px] text-white hover:bg-ink/90">
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
