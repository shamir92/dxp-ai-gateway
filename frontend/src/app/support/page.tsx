import {
  BookOpen,
  CreditCard,
  Key,
  LifeBuoy,
  MessageSquare,
  Mail,
  Users,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { AppShell } from "@/components/app-shell";
import { StatTile, TopSearchBar } from "@/components/token-plan-ui";

const channels = [
  {
    icon: MessageSquare,
    title: "Chat with support",
    desc: "Talk to a specialist for technical and account questions.",
    cta: "Start chat",
  },
  {
    icon: Mail,
    title: "Email support",
    desc: "Send a detailed request with logs and attachments.",
    cta: "Create ticket",
  },
  {
    icon: Users,
    title: "Developer community",
    desc: "Ask questions and share solutions with other builders.",
    cta: "Join community",
  },
];

const topics = [
  { icon: Key, title: "API & authentication", desc: "Keys, endpoints, SDK setup" },
  { icon: BookOpen, title: "Models & tokens", desc: "Selection, limits, and consumption" },
  { icon: LifeBuoy, title: "Batch jobs", desc: "Input files, status, retries" },
  { icon: CreditCard, title: "Balance & billing", desc: "Funding, invoices, refunds" },
  { icon: Users, title: "Account & security", desc: "Profile, login, and 2FA" },
];

const tickets = [
  ["Batch output delay", "In progress", "#48219 · Updated 2 days ago"],
  ["Invoice correction", "Resolved", "#47802 · Updated 2 days ago"],
  ["Plugin permission", "Resolved", "#47194 · Updated 2 days ago"],
];

export default function SupportPage() {
  return (
    <AppShell
      title="Support center"
      subtitle="Get help with integration, billing, models, and account access."
      actions={
        <div className="flex items-center gap-3">
          <TopSearchBar />
          <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
            My tickets
          </Button>
        </div>
      }
    >
      <div className="mx-auto max-w-[1024px] space-y-8">
        <Card className="border-line-warm bg-white">
          <CardContent className="flex flex-wrap items-center gap-4 py-5">
            <span className="flex size-10 items-center justify-center rounded-full bg-success/10">
              <span className="size-2.5 rounded-full bg-success" />
            </span>
            <div className="flex-1">
              <p className="text-[14px] font-semibold text-ink">All services operational</p>
              <p className="text-[12px] text-muted-warm">
                Gateway, batch processing, billing, and plugins are healthy.
              </p>
            </div>
            <Button variant="ghost" className="text-[12px] text-flame">
              View system status ↗
            </Button>
          </CardContent>
        </Card>

        <div className="grid gap-4 md:grid-cols-3">
          {channels.map((c) => (
            <Card key={c.title} className="border-line-warm bg-white">
              <CardContent className="flex h-full flex-col p-6">
                <c.icon className="size-5 text-flame" />
                <h3 className="mt-4 text-[14px] font-semibold text-ink">{c.title}</h3>
                <p className="mt-1.5 text-[12px] leading-relaxed text-muted-warm">{c.desc}</p>
                <Button className="mt-5 rounded-full bg-ink text-[12px] text-white hover:bg-ink/90">
                  {c.cta}
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>

        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">Browse help topics</CardTitle>
            <div className="mt-2 flex items-center gap-2 rounded-full border border-line-warm px-4">
              <Input
                placeholder="Search help"
                className="h-10 border-0 bg-transparent p-0 text-[12px] shadow-none focus-visible:ring-0"
              />
            </div>
          </CardHeader>
          <CardContent className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {topics.map((t) => (
              <div
                key={t.title}
                className="flex items-start gap-3 rounded-xl border border-line-warm p-4"
              >
                <t.icon className="mt-0.5 size-4 text-flame" />
                <div>
                  <p className="text-[12px] font-semibold text-ink">{t.title}</p>
                  <p className="text-[11px] text-muted-warm">{t.desc}</p>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <div className="grid gap-6 lg:grid-cols-2">
          <Card className="border-line-warm bg-white">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle className="text-[15px]">Recent tickets</CardTitle>
              <Button variant="ghost" className="text-[12px] text-flame">
                View all
              </Button>
            </CardHeader>
            <CardContent className="space-y-3">
              {tickets.map(([title, status, meta]) => (
                <div
                  key={title as string}
                  className="flex items-center gap-3 rounded-xl border border-line-warm p-4"
                >
                  <div className="flex-1">
                    <p className="text-[12px] font-semibold text-ink">{title}</p>
                    <p className="text-[11px] text-muted-warm">{meta}</p>
                  </div>
                  <Badge
                    className={`rounded-full text-[10px] ${
                      status === "Resolved"
                        ? "bg-success/10 text-success"
                        : "bg-flame/10 text-flame"
                    }`}
                  >
                    {status}
                  </Badge>
                </div>
              ))}
            </CardContent>
          </Card>

          <Card className="border-line-warm bg-white">
            <CardHeader>
              <CardTitle className="text-[15px]">Your support plan</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-[18px] font-semibold text-ink">Developer support</p>
              <p className="mt-1 text-[12px] text-muted-warm">
                Typical first response · under 4 hours
              </p>
              <ul className="mt-5 space-y-2.5 text-[12px] text-ink">
                <li>Technical integration support</li>
                <li>Billing and account assistance</li>
                <li>Ticket history and attachments</li>
              </ul>
              <Button variant="outline" className="mt-6 rounded-full border-line-warm text-[12px]">
                Compare support plans
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </AppShell>
  );
}
