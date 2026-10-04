import Link from "next/link";
import {
  BookOpen,
  Boxes,
  CreditCard,
  FileText,
  Key,
  Layers,
  Search,
  Terminal,
  Zap,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { AppShell } from "@/components/app-shell";
import { TopSearchBar } from "@/components/token-plan-ui";

const shortcuts = [
  {
    icon: Zap,
    title: "Quickstart",
    desc: "Send your first request in under five minutes.",
  },
  {
    icon: FileText,
    title: "API reference",
    desc: "Explore endpoints, schemas, and errors.",
  },
  {
    icon: Boxes,
    title: "SDK libraries",
    desc: "Use official Python, Node, and Go clients.",
  },
  {
    icon: Layers,
    title: "Batch processing",
    desc: "Submit large asynchronous workloads.",
  },
];

const browse = [
  "Getting started",
  "Authentication",
  "Models",
  "Token plans",
  "Batch jobs",
  "Billing",
  "Errors & limits",
];

const guides = [
  ["Build your first chat completion", "Create a key, select a model, and stream a response.", "6 min"],
  ["Understand token consumption", "Measure input, output, cache, and model-level usage.", "8 min"],
  ["Process requests in batches", "Prepare JSONL input and monitor asynchronous jobs.", "12 min"],
  ["Set balance alerts", "Avoid service interruptions with threshold notifications.", "4 min"],
  ["Secure production credentials", "Rotate keys and apply environment-level restrictions.", "7 min"],
];

const updates = [
  ["OCT 04", "MiMo v2.6 models", "New model family and context limits."],
  ["SEP 28", "Batch retry policy", "Automatic retries and error exports."],
  ["SEP 20", "Usage API", "Daily model-level metrics endpoint."],
];

export default function DocumentationPage() {
  return (
    <AppShell
      title="Documentation"
      subtitle="Everything you need to integrate, operate, and scale the AI gateway."
      actions={
        <div className="flex items-center gap-3">
          <TopSearchBar />
          <Badge className="rounded-full bg-cream text-[11px] text-muted-warm">API v1.8</Badge>
        </div>
      }
    >
      <div className="mx-auto max-w-[1024px] space-y-10">
        <div>
          <h2 className="font-heading text-[28px] font-semibold text-ink">
            How can we help you build?
          </h2>
          <div className="mt-4 flex items-center gap-2 rounded-full border border-line-warm bg-white px-4">
            <Search className="size-4 text-muted-warm" />
            <Input
              placeholder="Search guides, APIs, models, and examples…"
              className="h-12 border-0 bg-transparent p-0 text-[13px] shadow-none focus-visible:ring-0"
            />
            <kbd className="rounded border border-line-warm px-1.5 py-0.5 text-[10px] text-muted-warm">
              ⌘ K
            </kbd>
          </div>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {shortcuts.map((s) => (
            <Card key={s.title} className="border-line-warm bg-white">
              <CardContent className="p-5">
                <s.icon className="size-5 text-flame" />
                <h3 className="mt-4 text-[14px] font-semibold text-ink">{s.title}</h3>
                <p className="mt-1.5 text-[12px] leading-relaxed text-muted-warm">{s.desc}</p>
              </CardContent>
            </Card>
          ))}
        </div>

        <div>
          <p className="text-[11px] font-semibold uppercase tracking-wider text-muted-warm">
            Browse docs
          </p>
          <div className="mt-3 flex flex-wrap gap-2">
            {browse.map((item) => (
              <Button
                key={item}
                variant="outline"
                className="rounded-full border-line-warm bg-white text-[12px]"
              >
                {item}
              </Button>
            ))}
          </div>
        </div>

        <div>
          <div className="flex items-center justify-between">
            <h3 className="text-[16px] font-semibold text-ink">Popular guides</h3>
            <Link href="#" className="text-[12px] text-flame">
              View all guides →
            </Link>
          </div>
          <div className="mt-4 space-y-3">
            {guides.map(([title, desc, time]) => (
              <div
                key={title}
                className="flex flex-wrap items-center gap-4 rounded-2xl border border-line-warm bg-white p-5"
              >
                <div className="flex size-10 items-center justify-center rounded-xl bg-cream">
                  <BookOpen className="size-4 text-flame" />
                </div>
                <div className="min-w-[200px] flex-1">
                  <p className="text-[13px] font-semibold text-ink">{title}</p>
                  <p className="mt-0.5 text-[12px] text-muted-warm">{desc}</p>
                </div>
                <Badge className="rounded-full bg-cream text-[10px] text-muted-warm">{time}</Badge>
              </div>
            ))}
          </div>
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          <Card className="border-line-warm bg-white">
            <CardContent className="p-6">
              <h3 className="text-[15px] font-semibold text-ink">Latest updates</h3>
              <div className="mt-4 space-y-4">
                {updates.map(([date, title, desc]) => (
                  <div key={title} className="flex gap-4">
                    <span className="w-14 shrink-0 text-[10px] font-semibold text-flame">
                      {date}
                    </span>
                    <div>
                      <p className="text-[13px] font-medium text-ink">{title}</p>
                      <p className="text-[11px] text-muted-warm">{desc}</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
          <Card className="border-line-warm bg-white">
            <CardContent className="p-6">
              <h3 className="text-[15px] font-semibold text-ink">Developer resources</h3>
              <div className="mt-4 space-y-3">
                {[
                  [Key, "Authentication guide", "API keys, rotation, and scopes"],
                  [Terminal, "SDK quickstart", "Python, Node, and Go examples"],
                  [CreditCard, "Billing & credits", "Plans, invoices, and usage"],
                ].map(([Icon, title, desc]) => {
                  const I = Icon as typeof Key;
                  return (
                    <div
                      key={title as string}
                      className="flex items-center gap-3 rounded-xl border border-line-warm p-3"
                    >
                      <I className="size-4 text-flame" />
                      <div>
                        <p className="text-[12px] font-medium text-ink">{title as string}</p>
                        <p className="text-[11px] text-muted-warm">{desc as string}</p>
                      </div>
                    </div>
                  );
                })}
              </div>
              <div className="mt-5 flex items-center gap-2 rounded-xl bg-success/10 px-4 py-3">
                <span className="size-2 rounded-full bg-success" />
                <p className="text-[12px] text-success">API operational · All regions healthy</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </AppShell>
  );
}
