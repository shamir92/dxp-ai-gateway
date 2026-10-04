import Link from "next/link";
import { Check } from "lucide-react";
import { PublicShell } from "@/components/public-shell";
import { Button } from "@/components/ui/button";
import { ButtonLink } from "@/components/button-link";
import { Badge } from "@/components/ui/badge";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

const plans = [
  {
    name: "Starter",
    blurb: "For prototypes and evaluation.",
    price: "$0",
    period: "",
    features: [
      "1M tokens included",
      "2 requests / second",
      "Community support",
      "Core models",
    ],
    cta: "Start free",
    popular: false,
  },
  {
    name: "Builder",
    blurb: "For individual developers shipping products.",
    price: "$29",
    period: "/ month",
    features: [
      "38B plan credits",
      "20 requests / second",
      "Batch processing",
      "Email support",
    ],
    cta: "Choose Builder",
    popular: false,
  },
  {
    name: "Scale",
    blurb: "For production teams with growing usage.",
    price: "$149",
    period: "/ month",
    features: [
      "150B plan credits",
      "100 requests / second",
      "Priority batch queue",
      "Developer support",
    ],
    cta: "Choose Scale",
    popular: true,
  },
  {
    name: "Enterprise",
    blurb: "For critical workloads and tailored controls.",
    price: "Custom",
    period: "",
    features: [
      "Custom token volume",
      "Dedicated capacity",
      "SLA and SSO",
      "Solutions engineering",
    ],
    cta: "Contact sales",
    popular: false,
  },
];

const modelGroups = [
  {
    title: "Chat & reasoning",
    items: ["mimo-v2.6-pro", "mimo-v2.6-flash", "mimo-v2.5-pro", "mimo-v2.5"],
  },
  {
    title: "Speech & voice",
    items: ["mimo-v2.5-tts", "voiceclone", "voicedesign", "ASR"],
  },
  {
    title: "Developer tools",
    items: ["OpenClaw", "Codex", "Claude Code", "MiMo Code"],
  },
];

const benefits = [
  {
    title: "Transparent usage",
    desc: "Track tokens and requests by model across every day.",
  },
  {
    title: "Secure by default",
    desc: "Dedicated keys, regional endpoints, and safe rotation.",
  },
  {
    title: "Batch efficiency",
    desc: "Process large workloads at lower off-peak rates.",
  },
  {
    title: "No lock-in",
    desc: "OpenAI and Anthropic-compatible API interfaces.",
  },
];

const faqs = [
  {
    q: "How are token plan credits calculated?",
    a: "Credits are consumed from model token usage. Input, output, and cache-hit tokens follow each model’s published conversion rate.",
  },
  {
    q: "Can I change plans at any time?",
    a: "Yes. Upgrade or downgrade whenever you need — changes apply to the next billing cycle.",
  },
  {
    q: "What happens when I reach the plan limit?",
    a: "Requests are queued or rejected based on your settings. You can enable overage or upgrade instantly.",
  },
  {
    q: "Do unused credits roll over?",
    a: "Monthly credits reset each cycle. Annual plans include a modest rollover window.",
  },
  {
    q: "Is batch inference included?",
    a: "Batch is included on Builder and above, with priority queues on Scale and Enterprise.",
  },
];

export default function PricingPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-20">
        <div className="mx-auto max-w-[1302px] text-center">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
            Simple, predictable pricing
          </p>
          <h1 className="mx-auto mt-4 max-w-[700px] font-heading text-[42px] font-semibold leading-tight text-ink">
            Flexible plans, built for real value.
          </h1>
          <p className="mx-auto mt-4 max-w-[560px] text-[14px] leading-relaxed text-muted-warm">
            Choose a plan that fits how you build. Every tier includes secure API
            access, transparent usage, and the latest production-ready models.
          </p>
          <div className="mt-8 inline-flex rounded-full border border-line-warm bg-white p-1">
            <button className="rounded-full bg-ink px-5 py-2 text-[12px] text-white">
              Monthly
            </button>
            <button className="rounded-full px-5 py-2 text-[12px] text-muted-warm">
              Annual · save 20%
            </button>
          </div>
        </div>
      </section>

      <section className="bg-cream px-14 pb-20">
        <div className="mx-auto max-w-[1302px]">
          <h2 className="text-center font-heading text-[22px] font-semibold text-ink">
            Plans for every stage
          </h2>
          <div className="mt-10 grid gap-5 md:grid-cols-2 lg:grid-cols-4">
            {plans.map((plan) => (
              <div
                key={plan.name}
                className={`relative flex flex-col rounded-2xl p-6 ${
                  plan.popular
                    ? "bg-night-card text-white shadow-xl"
                    : "border border-line-warm bg-white"
                }`}
              >
                {plan.popular ? (
                  <Badge className="absolute -top-3 left-6 rounded-full bg-flame px-3 text-[9px] text-white">
                    MOST POPULAR
                  </Badge>
                ) : null}
                <h3 className={`text-[18px] font-semibold ${plan.popular ? "text-white" : "text-ink"}`}>
                  {plan.name}
                </h3>
                <p className={`mt-1 text-[12px] ${plan.popular ? "text-[#91A0B6]" : "text-muted-warm"}`}>
                  {plan.blurb}
                </p>
                <div className="mt-6 flex items-baseline gap-1">
                  <span className={`font-heading text-[36px] font-semibold ${plan.popular ? "text-white" : "text-ink"}`}>
                    {plan.price}
                  </span>
                  {plan.period ? (
                    <span className={`text-[12px] ${plan.popular ? "text-[#91A0B6]" : "text-muted-warm"}`}>
                      {plan.period}
                    </span>
                  ) : null}
                </div>
                <ul className="mt-6 space-y-3">
                  {plan.features.map((f) => (
                    <li key={f} className="flex items-start gap-2 text-[12px]">
                      <Check
                        className={`mt-0.5 size-3.5 shrink-0 ${plan.popular ? "text-success" : "text-flame"}`}
                      />
                      <span className={plan.popular ? "text-white" : "text-ink"}>{f}</span>
                    </li>
                  ))}
                </ul>
                <ButtonLink
                  href={plan.name === "Enterprise" ? "/support" : "/token-plan"}
                  className={`mt-8 h-10 rounded-full text-[12px] ${
                    plan.popular
                      ? "bg-white text-ink hover:bg-white/90"
                      : "bg-ink text-white hover:bg-ink/90"
                  }`}
                >
                  {plan.cta}
                </ButtonLink>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="bg-white px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <h2 className="font-heading text-[28px] font-semibold text-ink">
            Use the models you need
          </h2>
          <p className="mt-2 text-[13px] text-muted-warm">
            All paid plans include access to production-ready chat, reasoning, voice,
            and multimodal models.
          </p>
          <div className="mt-10 grid gap-5 md:grid-cols-3">
            {modelGroups.map((group) => (
              <div key={group.title} className="rounded-2xl border border-line-warm p-6">
                <h3 className="text-[14px] font-semibold text-ink">{group.title}</h3>
                <ul className="mt-4 space-y-2.5">
                  {group.items.map((item) => (
                    <li
                      key={item}
                      className="flex items-center gap-2 rounded-lg bg-cream-card px-3 py-2 text-[12px] text-ink"
                    >
                      <Check className="size-3.5 text-flame" />
                      <span className="font-mono">{item}</span>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="bg-cream-deep px-14 py-20">
        <div className="mx-auto max-w-[1302px]">
          <h2 className="font-heading text-[28px] font-semibold text-ink">
            Built to scale with you
          </h2>
          <p className="mt-2 text-[13px] text-muted-warm">
            Predictable plans with the controls and visibility production teams expect.
          </p>
          <div className="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {benefits.map((b) => (
              <div key={b.title} className="rounded-2xl bg-white p-6">
                <div className="flex size-9 items-center justify-center rounded-xl bg-cream">
                  <Check className="size-4 text-flame" />
                </div>
                <h3 className="mt-4 text-[14px] font-semibold text-ink">{b.title}</h3>
                <p className="mt-2 text-[12px] leading-relaxed text-muted-warm">{b.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="bg-white px-14 py-20">
        <div className="mx-auto max-w-[760px]">
          <div className="flex flex-wrap items-end justify-between gap-4">
            <div>
              <h2 className="font-heading text-[28px] font-semibold text-ink">FAQ</h2>
              <p className="mt-2 text-[22px] font-semibold text-ink">Questions, answered.</p>
              <p className="mt-2 text-[13px] text-muted-warm">
                Need more detail? Explore the documentation or talk with our team about
                your workload.
              </p>
            </div>
            <ButtonLink href="/support" variant="outline" className="rounded-full border-line-warm">Contact us</ButtonLink>
          </div>
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
              Ready to build with DXP AI Gateway?
            </h2>
            <p className="mt-2 text-[13px] text-[#FFE4DD]">
              Create an account, generate a key, and send your first request today.
            </p>
          </div>
          <div className="flex gap-3">
            <ButtonLink href="/token-plan" className="h-11 rounded-full bg-ink px-6 text-white hover:bg-ink/90">Start building</ButtonLink>
            <ButtonLink href="/support" variant="outline"
              className="h-11 rounded-full border-white bg-white px-6 text-ink">Talk to sales</ButtonLink>
          </div>
        </div>
      </section>
    </PublicShell>
  );
}
