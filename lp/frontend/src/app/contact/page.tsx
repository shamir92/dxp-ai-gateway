import { Mail, MessageCircle, MapPin, ArrowRight } from "lucide-react";
import { PublicShell } from "@/components/public-shell";
import { ButtonLink } from "@/components/button-link";

const WA_NUMBER = "6285187963963";
const WA_DISPLAY = "+62 851-8796-3963";
const EMAIL = "hello@fromdxp.com";

const channels = [
  {
    icon: MessageCircle,
    title: "WhatsApp",
    desc: "Fastest for product, pricing, and API questions.",
    value: WA_DISPLAY,
    href: `https://wa.me/${WA_NUMBER}`,
    cta: "Chat on WhatsApp",
    external: true,
  },
  {
    icon: Mail,
    title: "Email",
    desc: "For partnerships, support tickets, and formal requests.",
    value: EMAIL,
    href: `mailto:${EMAIL}`,
    cta: "Send an email",
    external: false,
  },
];

export default function ContactPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-16">
        <div className="mx-auto max-w-[1302px]">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
            Contact
          </p>
          <h1 className="mt-3 max-w-[640px] font-heading text-[42px] font-semibold leading-tight text-ink">
            Talk to the DXP team.
          </h1>
          <p className="mt-4 max-w-[520px] text-[14px] leading-relaxed text-muted-warm">
            Questions about the AI gateway, API access, pricing, or partnerships —
            reach us on WhatsApp or email. We usually reply within one business day.
          </p>
        </div>
      </section>

      <section className="bg-white px-14 py-16">
        <div className="mx-auto max-w-[1302px]">
          <div className="grid gap-5 md:grid-cols-2">
            {channels.map((c) => (
              <div
                key={c.title}
                className="flex flex-col rounded-2xl border border-line-warm bg-cream-card p-8"
              >
                <div className="flex size-11 items-center justify-center rounded-xl bg-white">
                  <c.icon className="size-5 text-flame" />
                </div>
                <h2 className="mt-5 text-[20px] font-semibold text-ink">{c.title}</h2>
                <p className="mt-2 text-[13px] leading-relaxed text-muted-warm">{c.desc}</p>
                <p className="mt-5 font-mono text-[15px] text-ink">{c.value}</p>
                <ButtonLink
                  href={c.href}
                  className="mt-6 inline-flex h-11 w-fit items-center rounded-full bg-ink px-6 text-[13px] text-white hover:bg-ink/90"
                >
                  {c.cta}
                  <ArrowRight className="ml-1.5 size-3.5" />
                </ButtonLink>
              </div>
            ))}
          </div>

          <div className="mt-10 grid gap-5 rounded-2xl border border-line-warm bg-cream-card p-8 sm:grid-cols-[auto_1fr] sm:items-start">
            <div className="flex size-11 items-center justify-center rounded-xl bg-white">
              <MapPin className="size-5 text-flame" />
            </div>
            <div>
              <h2 className="text-[16px] font-semibold text-ink">Office</h2>
              <p className="mt-2 text-[13px] leading-relaxed text-muted-warm">
                PT DXP Global Inovasi
                <br />
                Jl. Tarumajaya Raya, Sagara Makmur, Kec. Tarumajaya,
                <br />
                Kabupaten Bekasi, Jawa Barat 17211, Indonesia
              </p>
            </div>
          </div>

          <p className="mt-10 text-[12px] text-muted-warm">
            Looking for account and API help while signed in? Visit{" "}
            <span className="text-ink">Support</span> in the console after you log in.
          </p>
        </div>
      </section>
    </PublicShell>
  );
}
