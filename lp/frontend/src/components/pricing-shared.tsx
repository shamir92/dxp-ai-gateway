import { Check } from "lucide-react";
import { ButtonLink } from "@/components/button-link";
import { Badge } from "@/components/ui/badge";

export type PlanFeature = {
  title: string;
  sub?: string;
};

export type Plan = {
  name: string;
  blurb: string;
  price: string;
  period: string;
  oldPrice?: string;
  features: PlanFeature[];
  cta: string;
  href: string;
  popular?: boolean;
};

export const modelGroups = [
  {
    title: "DeepSeek",
    items: ["deepseek-v4.1-flash", "deepseek-v4-flash"],
  },
  {
    title: "MiMo",
    items: ["mimo-v2.6-pro", "mimo-v2.6-flash"],
  },
  {
    title: "GLM · Kimi · MiniMax",
    items: ["glm-5.3-flash", "glm-5.3", "kimi-k3", "minimax-m3"],
  },
];

export function PlanGrid({ plans }: { plans: Plan[] }) {
  return (
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
          <div className="mt-5 flex items-baseline gap-2">
            <span className={`font-heading text-[34px] font-semibold ${plan.popular ? "text-white" : "text-ink"}`}>
              {plan.price}
            </span>
            {plan.period ? (
              <span className={`text-[12px] ${plan.popular ? "text-[#91A0B6]" : "text-muted-warm"}`}>
                {plan.period}
              </span>
            ) : null}
            {plan.oldPrice ? (
              <span className={`text-[12px] line-through ${plan.popular ? "text-[#697689]" : "text-muted-warm"}`}>
                {plan.oldPrice}
              </span>
            ) : null}
          </div>
          <ButtonLink
            href={plan.href}
            className={`mt-5 h-10 w-full rounded-full text-[12px] ${
              plan.popular
                ? "bg-white text-ink hover:bg-white/90"
                : "bg-ink text-white hover:bg-ink/90"
            }`}
          >
            {plan.cta}
          </ButtonLink>
          <ul className="mt-6 space-y-3.5">
            {plan.features.map((f) => (
              <li key={f.title} className="flex items-start gap-2">
                <Check
                  className={`mt-0.5 size-3.5 shrink-0 ${plan.popular ? "text-success" : "text-flame"}`}
                />
                <div>
                  <p className={`text-[12px] font-medium ${plan.popular ? "text-white" : "text-ink"}`}>
                    {f.title}
                  </p>
                  {f.sub ? (
                    <p className={`mt-0.5 text-[11px] ${plan.popular ? "text-[#91A0B6]" : "text-muted-warm"}`}>
                      {f.sub}
                    </p>
                  ) : null}
                </div>
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
}

export function ModelsSection({
  title = "Use the models you need",
  blurb = "Every plan includes the same production-ready model catalog.",
}: {
  title?: string;
  blurb?: string;
}) {
  return (
    <section className="bg-white px-14 py-20">
      <div className="mx-auto max-w-[1302px]">
        <h2 className="font-heading text-[28px] font-semibold text-ink">{title}</h2>
        <p className="mt-2 text-[13px] text-muted-warm">{blurb}</p>
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
  );
}
