import { ArrowRight, Terminal } from "lucide-react";
import { PublicShell } from "@/components/public-shell";
import { ButtonLink } from "@/components/button-link";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { PlanGrid, ModelsSection, type Plan } from "@/components/pricing-shared";

const codingTools = [
  "MiMo Desktop",
  "MiMo Code",
  "OpenCode",
  "OpenClaw",
  "Claude Code",
  "Codex",
  "Hermes Agent",
  "Kilo Code",
  "Chatbox AI",
  "Cline",
  "Cherry Studio",
  "Qwen Code",
  "CodeBuddy",
];

const quickStart = [
  {
    step: "1",
    title: "Subscription Plans",
    desc: "Choose and complete your subscription.",
    cta: "Subscribe Now",
    href: "#plans",
  },
  {
    step: "2",
    title: "Install Code Tools",
    desc: "Install your preferred coding tools.",
    cta: "Quick Integration",
    href: "#tools",
  },
  {
    step: "3",
    title: "Configure Your Plan",
    desc: "Set up the model, your dedicated API Key, and Base URL.",
    cta: "Configuration Guide",
    href: "/contact",
  },
  {
    step: "4",
    title: "Get Started",
    desc: "Unleash your productivity.",
    cta: "",
    href: "",
  },
];

const flagship: Plan["features"] = [
  { title: "V2.6 Flagship Model Access", sub: "Text / Multimodal / Speech" },
  { title: "Coding Framework Support", sub: "OpenClaw / Codex / MiMo Claw etc." },
  { title: "Free TTS Models", sub: "Limited time offer" },
  { title: "Night 0.8x Usage", sub: "00:00–08:00 UTC+8" },
  { title: "Unlimited Usage", sub: "No weekly cap · No 5-hour cap" },
];

const tokenPlans: Plan[] = [
  {
    name: "Lite",
    blurb: "Light use.",
    price: "$5.28",
    period: "/ month",
    oldPrice: "$6",
    features: [
      { title: "4.1 Billion Credits", sub: "Monthly package total" },
      ...flagship,
    ],
    cta: "Subscribe",
    href: "/register",
  },
  {
    name: "Standard",
    blurb: "Daily work.",
    price: "$14.08",
    period: "/ month",
    oldPrice: "$16",
    features: [
      { title: "11 Billion Credits", sub: "≈ 2.7 × Lite" },
      ...flagship,
    ],
    cta: "Subscribe",
    href: "/register",
    popular: true,
  },
  {
    name: "Pro",
    blurb: "Pro development.",
    price: "$44",
    period: "/ month",
    oldPrice: "$50",
    features: [
      { title: "38 Billion Credits", sub: "≈ 9.3 × Lite" },
      ...flagship,
    ],
    cta: "Subscribe",
    href: "/register",
  },
  {
    name: "Max",
    blurb: "Power development.",
    price: "$88",
    period: "/ month",
    oldPrice: "$100",
    features: [
      { title: "82 Billion Credits", sub: "≈ 20 × Lite" },
      ...flagship,
    ],
    cta: "Subscribe",
    href: "/register",
  },
];

const benefits = [
  {
    title: "Same models, monthly credits",
    desc: "Flagship V2.6 access on every plan — pick the credit pack that matches your pace.",
  },
  {
    title: "Coding tools included",
    desc: "Works with OpenClaw, Codex, MiMo Claw, and other frameworks.",
  },
  {
    title: "Night discount",
    desc: "0.8× credit usage from 00:00–08:00 UTC+8.",
  },
  {
    title: "Unlimited usage",
    desc: "No weekly cap and no 5-hour cap on paid monthly plans.",
  },
];

const faqs = [
  {
    q: "What is included in each monthly plan?",
    a: "Every plan includes V2.6 flagship model access (text, multimodal, speech), coding framework support, free TTS models (limited time), night 0.8× usage, and unlimited usage. Plans differ by monthly credit volume: Lite 4.1B, Standard 11B, Pro 38B, Max 82B.",
  },
  {
    q: "How does billing work?",
    a: "Plans are billed monthly per individual account. You can pay monthly; optional first-purchase or annual discounts may be shown at checkout when available.",
  },
  {
    q: "When is night 0.8× usage active?",
    a: "Off-peak hours are 00:00–08:00 UTC+8 (Beijing). During those hours, usage is charged at a 0.8× credits consumption rate.",
  },
  {
    q: "Does the plan renew automatically?",
    a: "Yes. Monthly plans support automatic renewal. Cancellations and refunds are not supported except where required by law.",
  },
  {
    q: "Which models can I use?",
    a: "The MiMo V2.6 series is available on all plans. Personal plans may also include additional models as listed in the console. The catalog can grow over time.",
  },
  {
    q: "How do I use Token Plan in coding tools?",
    a: "Configure your dedicated API key and Base URL in tools such as OpenClaw, Codex, MiMo Claw, Claude Code, or Cursor. Usage across tools counts toward the same monthly allowance.",
  },
  {
    q: "Is this different from API plans?",
    a: "Yes. Token plans are individual monthly credit packs for tools and agents. API plans are for products that call the gateway over HTTP — see API pricing.",
  },
];

export default function PricingPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-20">
        <div className="mx-auto max-w-[1302px] text-center">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
            Token plans
          </p>
          <h1 className="mx-auto mt-4 max-w-[760px] font-heading text-[42px] font-semibold leading-tight text-ink">
            Monthly credits. One flagship model family.
          </h1>
          <p className="mx-auto mt-4 max-w-[620px] text-[14px] leading-relaxed text-muted-warm">
            Individual monthly plans with generous credits — built for coding
            agents, editors, and everyday AI work.
          </p>
          <div className="mt-6 inline-flex items-center rounded-full border border-line-warm bg-white px-4 py-2 text-[12px] text-muted-warm">
            Individual · Monthly
          </div>
          <div className="mt-8">
            <ButtonLink
              href="/api-pricing"
              variant="outline"
              className="h-11 rounded-full border-line-warm bg-white px-6 text-[13px] text-ink"
            >
              Looking for API plans? <ArrowRight className="ml-1.5 size-3.5" />
            </ButtonLink>
          </div>
        </div>
      </section>

      <section id="plans" className="bg-cream px-14 pb-20">
        <div className="mx-auto max-w-[1302px]">
          <div className="flex flex-wrap items-end justify-between gap-3">
            <h2 className="font-heading text-[22px] font-semibold text-ink">
              Choose your plan
            </h2>
            <p className="text-[12px] text-muted-warm">
              Monthly billing · credits refresh each cycle
            </p>
          </div>
          <PlanGrid plans={tokenPlans} />
        </div>
      </section>

      {/* Coding tools */}
      <section id="tools" className="bg-white px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <h2 className="text-center font-heading text-[32px] font-semibold text-ink">
            Coding Tools Support
          </h2>
          <p className="mx-auto mt-3 max-w-[560px] text-center text-[13px] text-muted-warm">
            Use the same monthly credits in the coding tools and agents you already
            use every day.
          </p>
          <div className="mx-auto mt-10 grid max-w-[1100px] grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
            {codingTools.map((tool) => (
              <div
                key={tool}
                className="flex items-center gap-3 rounded-2xl border border-line-warm bg-cream-card px-4 py-4"
              >
                <span className="flex size-9 shrink-0 items-center justify-center rounded-full bg-white">
                  <Terminal className="size-4 text-ink" />
                </span>
                <span className="text-[13px] font-medium text-ink">{tool}</span>
              </div>
            ))}
            <div className="flex items-center justify-center rounded-2xl border border-line-warm bg-cream-card px-4 py-4">
              <span className="text-[18px] font-semibold text-muted-warm">···</span>
            </div>
          </div>
          <p className="mx-auto mt-6 max-w-[560px] text-center text-[12px] text-muted-warm">
            Compatible with OpenAI &amp; Anthropic — any tool that can set a Base
            URL and API key works.
          </p>
        </div>
      </section>

      {/* Quick start */}
      <section className="bg-cream px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <h2 className="text-center font-heading text-[32px] font-semibold text-ink">
            Quick Start
          </h2>
          <div className="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {quickStart.map((s) => (
              <div
                key={s.step}
                className="flex min-h-[220px] flex-col rounded-2xl border border-line-warm bg-white p-6"
              >
                <div className="flex items-center gap-2">
                  <span className="flex size-6 items-center justify-center rounded-full bg-cream text-[11px] font-semibold text-ink">
                    {s.step}
                  </span>
                  <h3 className="text-[15px] font-semibold text-ink">{s.title}</h3>
                </div>
                <p className="mt-4 text-[12px] leading-relaxed text-muted-warm">{s.desc}</p>
                {s.cta && s.href ? (
                  <ButtonLink
                    href={s.href}
                    variant="outline"
                    className="mt-auto h-9 w-fit rounded-full border-line-warm text-[11px] text-ink"
                  >
                    {s.cta} <ArrowRight className="ml-1 size-3" />
                  </ButtonLink>
                ) : null}
              </div>
            ))}
          </div>
        </div>
      </section>

      <ModelsSection />

      <section className="bg-cream-deep px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <h2 className="font-heading text-[28px] font-semibold text-ink">
            Built for individual builders
          </h2>
          <p className="mt-2 text-[13px] text-muted-warm">
            Clear monthly credits, tool-friendly access, and no surprise caps.
          </p>
          <div className="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {benefits.map((b) => (
              <div key={b.title} className="rounded-2xl bg-white p-6">
                <h3 className="text-[14px] font-semibold text-ink">{b.title}</h3>
                <p className="mt-2 text-[12px] leading-relaxed text-muted-warm">{b.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="bg-white px-14 py-20">
        <div className="mx-auto max-w-[760px]">
          <h2 className="font-heading text-[28px] font-semibold text-ink">FAQ</h2>
          <p className="mt-2 text-[22px] font-semibold text-ink">Token plans, answered.</p>
          <p className="mt-2 text-[13px] text-muted-warm">
            Need more detail? Contact us — we are happy to help.
          </p>
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
              Ready for your monthly credit plan?
            </h2>
            <p className="mt-2 text-[13px] text-[#FFE4DD]">
              Subscribe, connect your tools, and start building today.
            </p>
          </div>
          <div className="flex gap-3">
            <ButtonLink
              href="/register"
              className="h-11 rounded-full bg-ink px-6 text-white hover:bg-ink/90"
            >
              Subscribe <ArrowRight className="ml-1.5 size-3.5" />
            </ButtonLink>
            <ButtonLink
              href="/contact"
              variant="outline"
              className="h-11 rounded-full border-white bg-white px-6 text-ink"
            >
              Contact us
            </ButtonLink>
          </div>
        </div>
      </section>
    </PublicShell>
  );
}
