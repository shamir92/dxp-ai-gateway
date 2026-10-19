import { ArrowRight, Clock, Mail, MessageCircle } from "lucide-react";
import { ButtonLink } from "@/components/button-link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AppShell } from "@/components/app-shell";
import { TopSearchBar } from "@/components/token-plan-ui";
import { RequireSession } from "@/components/require-session";

const WA_NUMBER = "6285187963963";
const WA_DISPLAY = "+62 851-8796-3963";
const EMAIL = "hello@fromdxp.com";

const channels = [
  {
    icon: MessageCircle,
    title: "WhatsApp",
    desc: "Fastest for integration, billing, and account questions.",
    value: WA_DISPLAY,
    href: `https://wa.me/${WA_NUMBER}`,
    cta: "Chat on WhatsApp",
  },
  {
    icon: Mail,
    title: "Email",
    desc: "Send a detailed request with logs or screenshots attached.",
    value: EMAIL,
    href: `mailto:${EMAIL}`,
    cta: "Send an email",
  },
];



export default function SupportPage() {
  return (
    <RequireSession>
      <AppShell
        title="Support"
        subtitle="Reach the DXP team about integration, billing, models, and account access."
        actions={<TopSearchBar />}
      >
        <div className="mx-auto max-w-[1024px] space-y-8">
          <Card className="border-line-warm bg-white">
            <CardHeader>
              <CardTitle className="text-[15px]">Contact us</CardTitle>
              <p className="mt-1 text-[12px] text-muted-warm">
                We answer on WhatsApp and email — usually within one business day.
              </p>
            </CardHeader>
            <CardContent>
              <div className="grid gap-5 md:grid-cols-2">
                {channels.map((c) => (
                  <div
                    key={c.title}
                    className="flex flex-col rounded-2xl border border-line-warm bg-cream-card p-6"
                  >
                    <div className="flex size-11 items-center justify-center rounded-xl bg-white">
                      <c.icon className="size-5 text-flame" />
                    </div>
                    <h3 className="mt-5 text-[16px] font-semibold text-ink">{c.title}</h3>
                    <p className="mt-2 text-[12px] leading-relaxed text-muted-warm">
                      {c.desc}
                    </p>
                    <p className="mt-4 font-mono text-[13px] text-ink">{c.value}</p>
                    <ButtonLink
                      href={c.href}
                      className="mt-5 inline-flex h-10 w-fit items-center rounded-full bg-ink px-5 text-[12px] text-white hover:bg-ink/90"
                    >
                      {c.cta}
                      <ArrowRight className="ml-1.5 size-3.5" />
                    </ButtonLink>
                  </div>
                ))}
              </div>
              <p className="mt-5 flex items-center gap-2 text-[11px] text-muted-warm">
                <Clock className="size-3.5" />
                Typical first response: under 4 hours on business days.
              </p>
            </CardContent>
          </Card>

          <Card className="border-line-warm bg-white">
            <CardContent className="py-5">
              <p className="text-[12px] leading-relaxed text-muted-warm">
                For the quickest answer, include your account email, the model or endpoint
                involved, and what you expected to happen. Please do not send API keys —
                we never need them to investigate.
              </p>
            </CardContent>
          </Card>
        </div>
      </AppShell>
    </RequireSession>
  );
}