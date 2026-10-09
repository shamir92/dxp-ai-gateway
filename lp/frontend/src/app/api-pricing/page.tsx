import { ArrowRight, Info } from "lucide-react";
import { PublicShell } from "@/components/public-shell";
import { ButtonLink } from "@/components/button-link";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

/** Overseas rates, USD per 1M tokens (from pay-as-you-go reference). */
const languageModels = [
  {
    inference: "Real-time API",
    id: "mimo-v2.6-pro",
    cacheHit: "$0.0036",
    cacheMiss: "$0.435",
    output: "$0.87",
  },
  {
    inference: "",
    id: "mimo-v2.6-flash",
    cacheHit: "$0.0028",
    cacheMiss: "$0.14",
    output: "$0.28",
  },
  {
    inference: "Batch API",
    id: "mimo-v2.6-pro",
    cacheHit: "$0.0018",
    cacheMiss: "$0.2175",
    output: "$0.435",
  },
  {
    inference: "",
    id: "mimo-v2.6-flash",
    cacheHit: "$0.0014",
    cacheMiss: "$0.07",
    output: "$0.14",
  },
];

const otherModels = [
  { id: "deepseek-v4.1-flash", vendor: "DeepSeek" },
  { id: "deepseek-v4-flash", vendor: "DeepSeek" },
  { id: "glm-5.3-flash", vendor: "GLM" },
  { id: "glm-5.3", vendor: "GLM" },
  { id: "kimi-k3", vendor: "Kimi" },
  { id: "minimax-m3", vendor: "MiniMax" },
];

const billingNotes = [
  "Billing unit: USD / 1M tokens (overseas)",
  "Cache hit: billed at the cache-hit rate when the request prefix hits the prompt cache",
  "Cache write: limited-time free",
  "ASR models: billed on input audio duration (accurate to the second, converted to hourly)",
  "Internet search / web tools: billed independently per call, not included in token price",
];

const faqs = [
  {
    q: "How is API usage billed?",
    a: "Pay-as-you-go by token usage. Input is priced separately for cache hit vs cache miss; output is billed per token. You can see spend by model and API key in the console.",
  },
  {
    q: "Do you store my prompts for billing?",
    a: "No. We do not store request/response content. We log usage metadata only (model, token counts, cost, status) for billing and analytics.",
  },
  {
    q: "What is batch API pricing?",
    a: "Batch jobs use discounted rates versus real-time (about half in the table above). Batch is not available for ultra-speed or deprecated models.",
  },
  {
    q: "Are there rate limits?",
    a: "Yes — by plan and key. Start free with shared limits, or talk to us about higher RPS and dedicated capacity.",
  },
  {
    q: "Do you support OpenAI-compatible clients?",
    a: "Yes. Point your Base URL and API key at DXP. Streaming, tools, and structured output are supported on real-time models.",
  },
];

function PriceTable() {
  return (
    <div className="overflow-x-auto rounded-2xl border border-line-warm bg-white">
      <table className="w-full min-w-[680px] text-left">
        <thead>
          <tr className="border-b border-line-warm bg-cream-card">
            <th className="px-5 py-4 text-[11px] font-semibold text-muted-warm">Inference type</th>
            <th className="px-5 py-4 text-[11px] font-semibold text-muted-warm">Model name</th>
            <th className="px-5 py-4 text-[11px] font-semibold text-muted-warm">Input (cache hit)</th>
            <th className="px-5 py-4 text-[11px] font-semibold text-muted-warm">Input (cache miss)</th>
            <th className="px-5 py-4 text-[11px] font-semibold text-muted-warm">Output</th>
          </tr>
        </thead>
        <tbody>
          {languageModels.map((m, i) => (
            <tr key={`${m.inference}-${m.id}-${i}`} className="border-b border-line-warm last:border-0">
              <td className="px-5 py-4 text-[12px] text-ink">{m.inference || ""}</td>
              <td className="px-5 py-4">
                <div className="flex flex-wrap items-center gap-2">
                  <code className="rounded-md bg-cream-card px-2 py-1 text-[12px] text-ink">
                    {m.id}
                  </code>
                </div>
              </td>
              <td className="px-5 py-4 font-mono text-[12px] text-ink">{m.cacheHit}</td>
              <td className="px-5 py-4 font-mono text-[12px] text-ink">{m.cacheMiss}</td>
              <td className="px-5 py-4 font-mono text-[12px] text-ink">{m.output}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <p className="border-t border-line-warm px-5 py-3 text-[11px] text-muted-warm">
        Units: USD per 1M tokens · Real-time and batch rates · Last updated for the published catalog
      </p>
    </div>
  );
}

export default function ApiPricingPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-16">
        <div className="mx-auto max-w-[1302px]">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
            API pricing
          </p>
          <h1 className="mt-3 max-w-[720px] font-heading text-[42px] font-semibold leading-tight text-ink">
            Pay-as-you-go for the DXP API.
          </h1>
          <p className="mt-4 max-w-[640px] text-[14px] leading-relaxed text-muted-warm">
            Use a standard API key and consume balance based on actual token
            usage. Pay-as-you-go is not interoperable with Token Plan package
            quota — see{" "}
            <a href="/pricing" className="text-flame">
              Token plans
            </a>{" "}
            for monthly credit packs for tools.
          </p>
        </div>
      </section>

      <section className="bg-white px-14 py-12">
        <div className="mx-auto max-w-[1302px]">
          <div className="rounded-2xl border border-line-warm bg-cream-card p-6">
            <div className="flex items-center gap-2">
              <Info className="size-4 text-flame" />
              <h2 className="text-[15px] font-semibold text-ink">Billing instructions</h2>
            </div>
            <ul className="mt-4 list-disc space-y-2 pl-5 text-[13px] leading-relaxed text-muted-warm">
              {billingNotes.map((n) => (
                <li key={n}>{n}</li>
              ))}
            </ul>
          </div>
        </div>
      </section>

      <section className="bg-cream px-14 pb-16">
        <div className="mx-auto max-w-[1302px]">
          <h2 className="font-heading text-[28px] font-semibold text-ink">
            Language models
          </h2>
          <p className="mt-2 text-[13px] text-muted-warm">
            Overseas pricing · USD / 1M tokens · cache hit and cache miss input rates.
          </p>
          <div className="mt-8">
            <PriceTable />
          </div>

          <div className="mt-12">
            <h2 className="font-heading text-[22px] font-semibold text-ink">
              Batch API
            </h2>
            <p className="mt-2 text-[13px] text-muted-warm">
              Discounted async rates (see table). Not supported on ultra-speed
              models. Status and retries available in the console.
            </p>
          </div>

          <div className="mt-12">
            <h2 className="font-heading text-[22px] font-semibold text-ink">
              Additional models
            </h2>
            <p className="mt-2 text-[13px] text-muted-warm">
              Also available on the gateway — rates in the console or from our
              team:
            </p>
            <div className="mt-5 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
              {otherModels.map((m) => (
                <div
                  key={m.id}
                  className="flex items-center justify-between rounded-2xl border border-line-warm bg-white px-5 py-4"
                >
                  <div>
                    <p className="text-[11px] font-medium text-flame">{m.vendor}</p>
                    <code className="mt-1 block text-[13px] text-ink">{m.id}</code>
                  </div>
                  <ButtonLink
                    href="/contact"
                    variant="outline"
                    className="h-9 rounded-full border-line-warm text-[11px]"
                  >
                    Get rates
                  </ButtonLink>
                </div>
              ))}
            </div>
          </div>

          <div className="mt-12 rounded-2xl border border-line-warm bg-white p-6">
            <h2 className="text-[15px] font-semibold text-ink">ASR & TTS</h2>
            <ul className="mt-3 list-disc space-y-2 pl-5 text-[13px] text-muted-warm">
              <li>
                <strong className="text-ink">ASR</strong> — billed on input audio
                duration (accurate to the second, converted to hourly). Rate in
                console (e.g. ~$0.074 / h reference).
              </li>
              <li>
                <strong className="text-ink">TTS</strong> — free for a limited
                time on supported speech-synthesis models.
              </li>
            </ul>
          </div>
        </div>
      </section>

      <section className="bg-cream-deep px-14 py-16">
        <div className="mx-auto max-w-[1302px]">
          <h2 className="font-heading text-[28px] font-semibold text-ink">
            Specs of the API
          </h2>
          <div className="mt-8 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {[
              {
                title: "OpenAI-compatible",
                desc: "Chat completions, streaming (SSE), tools, and structured output.",
              },
              {
                title: "Anthropic-compatible",
                desc: "Messages API for Claude Code and other Anthropic clients.",
              },
              {
                title: "Prompt cache",
                desc: "Cache hit billed at a lower input rate. Cache write free for a limited time.",
              },
              {
                title: "Batch",
                desc: "Lower rates for async workloads with status and retries.",
              },
              {
                title: "Usage & keys",
                desc: "Scoped API keys, spend by model, and token usage logs for billing only.",
              },
              {
                title: "No content storage",
                desc: "Prompts and completions are not stored on the gateway after the request.",
              },
            ].map((b) => (
              <div key={b.title} className="rounded-2xl border border-line-warm bg-white p-6">
                <h3 className="text-[14px] font-semibold text-ink">{b.title}</h3>
                <p className="mt-2 text-[12px] leading-relaxed text-muted-warm">{b.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="bg-white px-14 py-16">
        <div className="mx-auto max-w-[760px]">
          <h2 className="font-heading text-[28px] font-semibold text-ink">FAQ</h2>
          <p className="mt-2 text-[13px] text-muted-warm">Common questions about API billing.</p>
          <Accordion className="mt-8">
            {faqs.map((faq, i) => (
              <AccordionItem key={faq.q} value={`item-${i}`}>
                <AccordionTrigger className="text-left text-[14px] text-ink">
                  {faq.q}
                </AccordionTrigger>
                <AccordionContent className="text-[13px] text-muted-warm">
                  {faq.a}
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </div>
      </section>

      <section className="bg-flame px-14 py-16">
        <div className="mx-auto flex max-w-[1302px] flex-wrap items-center justify-between gap-8">
          <div>
            <h2 className="font-heading text-[28px] font-semibold text-white">
              Start with pay-as-you-go.
            </h2>
            <p className="mt-2 text-[13px] text-[#FFE4DD]">
              Create a key and send your first request — or talk to us for volume
              pricing.
            </p>
          </div>
          <div className="flex gap-3">
            <ButtonLink
              href="/register"
              className="h-11 rounded-full bg-ink px-6 text-white hover:bg-ink/90"
            >
              Get API key <ArrowRight className="ml-1.5 size-3.5" />
            </ButtonLink>
            <ButtonLink
              href="/contact"
              variant="outline"
              className="h-11 rounded-full border-white bg-white px-6 text-ink"
            >
              Contact sales
            </ButtonLink>
          </div>
        </div>
      </section>
    </PublicShell>
  );
}
