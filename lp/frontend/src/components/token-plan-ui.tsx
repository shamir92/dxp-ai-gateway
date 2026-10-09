import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Copy, ExternalLink, Key } from "lucide-react";

export function TopSearchBar() {
  return (
    <div className="flex items-center gap-3">
      <div className="flex h-9 min-w-[260px] items-center gap-2 rounded-full border border-line-warm bg-white px-4 text-[12px] text-muted-warm">
        <span>⌕</span>
        Search models, keys, requests…
      </div>
      <div className="flex size-9 items-center justify-center rounded-full bg-prompt text-[11px] font-semibold text-white">
        SK
      </div>
    </div>
  );
}

export function PlanUsageCard() {
  return (
    <Card className="border-line-warm bg-white">
      <CardHeader className="flex flex-row items-start justify-between gap-4 pb-2">
        <div>
          <CardTitle className="text-[13px] font-medium text-muted-warm">
            CURRENT PLAN USAGE
          </CardTitle>
          <p className="mt-2 font-heading text-[28px] font-semibold text-ink">
            26,008,233,659
          </p>
          <p className="mt-1 text-[12px] text-muted-warm">
            of 38,000,000,000 · 68.4%
          </p>
        </div>
        <Badge className="rounded-full bg-success/10 text-success">Active</Badge>
      </CardHeader>
      <CardContent>
        <Progress value={68.4} className="h-2" />
        <div className="mt-5 grid grid-cols-2 gap-4 sm:grid-cols-4">
          {[
            ["Remaining", "11.99B"],
            ["Daily average", "867M"],
            ["Forecast", "Oct 08"],
            ["Renews", "Oct 10"],
          ].map(([label, value]) => (
            <div key={label}>
              <p className="text-[10px] uppercase tracking-wide text-muted-warm">{label}</p>
              <p className="mt-1 text-[14px] font-semibold text-ink">{value}</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

export function PlanMetaStrip() {
  return (
    <Card className="border-line-warm bg-white">
      <CardContent className="flex flex-wrap items-center gap-x-10 gap-y-4 py-5">
        <div>
          <p className="text-[11px] font-medium text-muted-warm">Individual plan</p>
          <p className="mt-1 text-[14px] font-semibold text-ink">Pro monthly plan</p>
          <p className="text-[11px] text-muted-warm">
            Valid until 2026-10-10 23:59:59 UTC
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-6 text-[11px]">
          <div>
            <p className="text-muted-warm">Started</p>
            <p className="font-medium text-ink">Sep 10</p>
          </div>
          <div>
            <p className="text-muted-warm">Current</p>
            <p className="font-medium text-ink">Oct 04</p>
          </div>
          <div>
            <p className="text-muted-warm">Renews</p>
            <p className="font-medium text-ink">Oct 10</p>
          </div>
        </div>
        <div className="ml-auto flex gap-2">
          <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
            Auto-renewal management
          </Button>
          <Button className="rounded-full bg-ink text-[12px] text-white hover:bg-ink/90">
            Subscribe plan
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}

export function ApiKeyCard() {
  return (
    <Card className="border-line-warm bg-white">
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle className="flex items-center gap-2 text-[15px]">
            <Key className="size-4 text-flame" /> Dedicated API key
          </CardTitle>
          <p className="mt-1 text-[12px] text-muted-warm">
            Keep your API key secure. Do not share it or expose it in browser or
            client-side code.
          </p>
        </div>
        <Button variant="ghost" className="gap-1 text-[11px] text-flame">
          Quick integration <ExternalLink className="size-3" />
        </Button>
      </CardHeader>
      <CardContent className="space-y-5">
        <div>
          <p className="text-[11px] font-medium text-muted-warm">API key</p>
          <div className="mt-2 flex items-center gap-2 rounded-xl border border-line-warm bg-cream-card px-4 py-3">
            <code className="flex-1 font-mono text-[12px] text-ink">
              tp-s44p2o••••••••••••••••••••zHy
            </code>
            <Button size="sm" variant="ghost" className="size-8 p-0">
              <Copy className="size-3.5" />
            </Button>
          </div>
          <p className="mt-2 text-[11px] text-muted-warm">
            This plan is intended for interactive AI coding and agent tools only.
          </p>
        </div>

        <div>
          <p className="text-[11px] font-medium text-muted-warm">Dedicated base URL</p>
          <div className="mt-2 space-y-2">
            {[
              ["OpenAI compatible", "https://token-plan-sgp.api.dxp.ai/v1"],
              ["Anthropic compatible", "https://token-plan-sgp.api.dxp.ai/anthropic"],
            ].map(([label, url]) => (
              <div
                key={label}
                className="flex flex-wrap items-center gap-2 rounded-xl border border-line-warm px-4 py-2.5"
              >
                <Badge variant="secondary" className="rounded-full bg-cream text-[10px]">
                  {label}
                </Badge>
                <code className="flex-1 font-mono text-[11px] text-ink">{url}</code>
                <Button size="sm" variant="ghost" className="size-7 p-0">
                  <Copy className="size-3.5" />
                </Button>
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export function PlanBenefitsCard() {
  return (
    <Card className="border-line-warm bg-white">
      <CardHeader>
        <CardTitle className="text-[15px]">Plan benefits</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4 text-[12px]">
        {[
          ["Models", "deepseek, mimo, glm, kimi, minimax"],
          ["Credits", "38 billion tokens"],
          ["Coding tools", "OpenClaw, Codex, Claude Code, MiMo Code and others"],
          ["Other benefits", "20% off during off-peak hours · selected TTS models included"],
        ].map(([label, value]) => (
          <div key={label}>
            <p className="font-medium text-muted-warm">{label}</p>
            <p className="mt-1 text-ink">{value}</p>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}

export function UsageBarChart({
  series,
  labels,
  max = 500,
}: {
  series: { name: string; color: string; values: number[] }[];
  labels: string[];
  max?: number;
}) {
  return (
    <div className="mt-4">
      <div className="flex h-[180px] items-end gap-1.5">
        {labels.map((_, i) => (
          <div key={i} className="flex flex-1 flex-col justify-end gap-0.5">
            {series.map((s) => {
              const v = s.values[i] ?? 0;
              const h = Math.max(2, (v / max) * 160);
              return (
                <div
                  key={s.name}
                  style={{ height: h, backgroundColor: s.color }}
                  className="rounded-sm"
                  title={`${s.name}: ${v}`}
                />
              );
            })}
          </div>
        ))}
      </div>
      <div className="mt-2 flex justify-between text-[9px] text-muted-warm">
        {labels.map((l) => (
          <span key={l}>{l}</span>
        ))}
      </div>
    </div>
  );
}

export function UsageLegend({
  items,
}: {
  items: { name: string; value: string; share: string; color: string }[];
}) {
  return (
    <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
      {items.map((item) => (
        <div key={item.name} className="rounded-xl border border-line-warm p-3">
          <div className="flex items-center gap-2">
            <span
              className="size-2 rounded-full"
              style={{ backgroundColor: item.color }}
            />
            <p className="text-[11px] text-muted-warm">{item.name}</p>
          </div>
          <p className="mt-1.5 text-[16px] font-semibold text-ink">{item.value}</p>
          <p className="text-[10px] text-muted-warm">{item.share}</p>
        </div>
      ))}
    </div>
  );
}

export function StatTile({
  label,
  value,
  hint,
}: {
  label: string;
  value: string;
  hint?: string;
}) {
  return (
    <Card className="border-line-warm bg-white">
      <CardContent className="py-5">
        <p className="text-[11px] uppercase tracking-wide text-muted-warm">{label}</p>
        <p className="mt-1.5 font-heading text-[22px] font-semibold text-ink">{value}</p>
        {hint ? <p className="mt-1 text-[11px] text-muted-warm">{hint}</p> : null}
      </CardContent>
    </Card>
  );
}
