"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Copy, ExternalLink, Key, Trash2 } from "lucide-react";
import { UserAvatar } from "@/components/user-avatar";
import { MODEL_PRICING } from "@/lib/model-pricing";
import { codeSamples, toolGuides, baseURLFor, KEY_PLACEHOLDER } from "@/lib/integration";
import {
  formatCredits,
  formatDateTime,
  formatKeyForDisplay,
  type ApiKey,
  type Endpoint,
  type SubscriptionView,
} from "@/lib/api-keys";
export function TopSearchBar() {
  return (
    <div className="flex items-center gap-3">
      <div className="flex h-9 min-w-[260px] items-center gap-2 rounded-full border border-line-warm bg-white px-4 text-[12px] text-muted-warm">
        <span>⌕</span>
        Search models, keys, requests…
      </div>
      <UserAvatar />
    </div>
  );
}

export function PlanUsageCard({
  view,
  loading,
}: {
  view: SubscriptionView | null;
  loading: boolean;
}) {
  if (loading) {
    return (
      <Card className="border-line-warm bg-white">
        <CardContent className="py-10 text-[13px] text-muted-warm">
          Loading plan usage…
        </CardContent>
      </Card>
    );
  }

  if (!view || view.status === "") {
    return (
      <Card className="border-line-warm bg-white">
        <CardHeader>
          <CardTitle className="text-[13px] font-medium text-muted-warm">
            CURRENT PLAN USAGE
          </CardTitle>
        </CardHeader>
        <CardContent className="text-[13px] text-muted-warm">
          No subscription is active on this platform yet.
        </CardContent>
      </Card>
    );
  }

  const unlimited = view.unlimited;
  const total = view.credits_total;
  const used = view.credits_used;
  const percent = !unlimited && total && total > 0 ? (used / total) * 100 : 0;
  const remaining = !unlimited && total ? Math.max(total - used, 0) : null;

  return (
    <Card className="border-line-warm bg-white">
      <CardHeader className="flex flex-row items-start justify-between gap-4 pb-2">
        <div>
          <CardTitle className="text-[13px] font-medium text-muted-warm">
            CURRENT PLAN USAGE
          </CardTitle>
          <p className="mt-2 font-heading text-[28px] font-semibold text-ink">
            {formatCredits(used)}
          </p>
          <p className="mt-1 text-[12px] text-muted-warm">
            {unlimited
              ? "of Unlimited"
              : `of ${formatCredits(total ?? 0)} · ${percent.toFixed(1)}%`}
          </p>
        </div>
        <Badge
          className={
            view.status === "active"
              ? "rounded-full bg-success/10 text-success"
              : "rounded-full bg-cream text-muted-warm"
          }
        >
          {view.status || "unknown"}
        </Badge>
      </CardHeader>
      <CardContent>
        {!unlimited ? <Progress value={percent} className="h-2" /> : null}
        <div className="mt-5 grid grid-cols-2 gap-4 sm:grid-cols-4">
          {[
            ["Remaining", unlimited ? "Unlimited" : formatCredits(remaining ?? 0)],
            ["Plan", view.plan_id ? `#${view.plan_id}` : "—"],
            ["Renews", formatDateTime(view.end_time)],
            ["Resets", formatDateTime(view.next_reset_time)],
          ].map(([label, value]) => (
            <div key={label}>
              <p className="text-[10px] uppercase tracking-wide text-muted-warm">
                {label}
              </p>
              <p className="mt-1 text-[14px] font-semibold text-ink">{value}</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

export function PlanMetaStrip({ view }: { view: SubscriptionView | null }) {
  return (
    <Card className="border-line-warm bg-white">
      <CardContent className="flex flex-wrap items-center gap-x-10 gap-y-4 py-5">
        <div>
          <p className="text-[11px] font-medium text-muted-warm">Individual plan</p>
          <p className="mt-1 text-[14px] font-semibold text-ink">
            {view && view.plan_id ? `Plan #${view.plan_id}` : "No active plan"}
          </p>
          <p className="text-[11px] text-muted-warm">
            {view && view.end_time
              ? `Valid until ${formatDateTime(view.end_time)} 23:59:59 UTC`
              : "Subscribe to activate a plan"}
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-6 text-[11px]">
          <div>
            <p className="text-muted-warm">Started</p>
            <p className="font-medium text-ink">
              {formatDateTime(view?.start_time ?? 0)}
            </p>
          </div>
          <div>
            <p className="text-muted-warm">Status</p>
            <p className="font-medium text-ink">{view?.status || "—"}</p>
          </div>
          <div>
            <p className="text-muted-warm">Renews</p>
            <p className="font-medium text-ink">
              {formatDateTime(view?.end_time ?? 0)}
            </p>
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

function CopyButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <Button
      size="sm"
      variant="ghost"
      className="size-8 p-0"
      aria-label="Copy"
      onClick={async () => {
        try {
          await navigator.clipboard.writeText(value);
          setCopied(true);
          setTimeout(() => setCopied(false), 1500);
        } catch {
          setCopied(false);
        }
      }}
    >
      <Copy className="size-3.5" />
      {copied ? <span className="sr-only">Copied</span> : null}
    </Button>
  );
}

export function ApiKeyCard({
  keys,
  loading,
  error,
  createdKey,
  onCreate,
  onDelete,
  onDismissCreated,
  endpoints,
}: {
  keys: ApiKey[];
  loading: boolean;
  error: string | null;
  createdKey: string | null;
  onCreate: (name: string) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
  onDismissCreated: () => void;
  endpoints: Endpoint[];
}) {
  const [showCreate, setShowCreate] = useState(false);
  const [name, setName] = useState("");
  const [busy, setBusy] = useState(false);

  const closeCreate = () => {
    setShowCreate(false);
    setName("");
  };

  return (
    <Card className="border-line-warm bg-white">
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle className="flex items-center gap-2 text-[15px]">
            <Key className="size-4 text-flame" /> Dedicated API keys
          </CardTitle>
          <p className="mt-1 text-[12px] text-muted-warm">
            Key values are shown once, when created. We only store a masked copy.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Dialog open={showCreate} onOpenChange={(open) => (open ? setShowCreate(true) : closeCreate())}>
            <DialogTrigger
              render={
                <Button className="h-8 rounded-full bg-ink px-4 text-[11px] text-white hover:bg-ink/90">
                  Create API key
                </Button>
              }
            />
            <DialogContent className="sm:max-w-[420px]">
              <DialogHeader>
                <DialogTitle className="text-[15px]">Create API key</DialogTitle>
                <DialogDescription className="text-[12px] text-muted-warm">
                  Give the key a name so you can recognise it later.
                </DialogDescription>
              </DialogHeader>
              <div className="py-2">
                <label htmlFor="api-key-name" className="text-[11px] font-medium text-muted-warm">
                  API key name
                </label>
                <Input
                  id="api-key-name"
                  value={name}
                  maxLength={50}
                  autoFocus
                  placeholder="e.g. Mimo Code"
                  className="mt-2 h-9 rounded-xl text-[12px]"
                  onChange={(e) => setName(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" && name.trim() !== "" && !busy) {
                      setBusy(true);
                      void onCreate(name.trim())
                        .then(closeCreate)
                        .finally(() => setBusy(false));
                    }
                  }}
                />
              </div>
              <DialogFooter>
                <Button
                  variant="outline"
                  className="h-9 rounded-full border-line-warm text-[12px]"
                  onClick={closeCreate}
                >
                  Cancel
                </Button>
                <Button
                  className="h-9 rounded-full bg-ink text-[12px] text-white hover:bg-ink/90"
                  disabled={busy || name.trim() === ""}
                  onClick={() => {
                    setBusy(true);
                    void onCreate(name.trim())
                      .then(closeCreate)
                      .finally(() => setBusy(false));
                  }}
                >
                  Create
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </CardHeader>
      <CardContent className="space-y-5">
        {createdKey ? (
          <div className="rounded-xl border border-success/40 bg-success/5 p-4">
            <p className="text-[11px] font-medium text-success">
              API key created — copy it now, it will not be shown again
            </p>
            <div className="mt-2 flex items-center gap-2 rounded-lg border border-line-warm bg-white px-3 py-2">
              <code className="flex-1 break-all font-mono text-[12px] text-ink">
                {formatKeyForDisplay(createdKey)}
              </code>
              <CopyButton value={formatKeyForDisplay(createdKey)} />
            </div>
            <Button
              variant="ghost"
              className="mt-2 h-7 text-[11px] text-muted-warm"
              onClick={onDismissCreated}
            >
              Dismiss
            </Button>
          </div>
        ) : null}

        {error ? <p className="text-[12px] text-flame">{error}</p> : null}

        <div>
          <p className="text-[11px] font-medium text-muted-warm">
            Your keys {loading ? "" : `(${keys.length})`}
          </p>
          <ul className="mt-2 space-y-2">
            {loading ? (
              <li className="rounded-xl border border-line-warm px-4 py-3 text-[12px] text-muted-warm">
                Loading keys…
              </li>
            ) : keys.length === 0 ? (
              <li className="rounded-xl border border-line-warm px-4 py-3 text-[12px] text-muted-warm">
                No API keys yet.
              </li>
            ) : (
              keys.map((k) => (
                <li
                  key={k.id}
                  className="flex flex-wrap items-center gap-3 rounded-xl border border-line-warm bg-cream-card px-4 py-3"
                >
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-[12px] font-medium text-ink">{k.name}</p>
                    <code className="font-mono text-[11px] text-muted-warm">
                      {formatKeyForDisplay(k.key)}
                    </code>
                  </div>
                  <span className="text-[11px] text-muted-warm">
                    {k.unlimited_quota ? "Unlimited" : formatCredits(k.remain_quota)}
                  </span>
                  <Button
                    size="sm"
                    variant="ghost"
                    className="size-8 p-0 text-flame"
                    aria-label={`Delete ${k.name}`}
                    onClick={() => void onDelete(k.id)}
                  >
                    <Trash2 className="size-3.5" />
                  </Button>
                </li>
              ))
            )}
          </ul>
        </div>

        <div>
          <p className="text-[11px] font-medium text-muted-warm">Base URL</p>
          <p className="mt-1 text-[11px] text-muted-warm">
            Compatible with OpenAI and Anthropic clients.
          </p>
          <div className="mt-2 space-y-2">
            {endpoints.length === 0 ? (
              <div className="rounded-xl border border-line-warm px-4 py-2.5 text-[11px] text-muted-warm">
                No base URL configured for this platform yet.
              </div>
            ) : (
              endpoints.map((endpoint) => (
                <div
                  key={endpoint.url}
                  className="flex flex-wrap items-center gap-2 rounded-xl border border-line-warm px-4 py-2.5"
                >
                  {endpoint.label ? (
                    <Badge variant="secondary" className="rounded-full bg-cream text-[10px]">
                      {endpoint.label}
                    </Badge>
                  ) : null}
                  <code className="flex-1 font-mono text-[11px] text-ink">{endpoint.url}</code>
                  <CopyButton value={endpoint.url} />
                </div>
              ))
            )}
          </div>
        </div>

        <QuickIntegrationDialog baseURL={endpoints[0]?.url ? baseURLFor(endpoints[0].url) : ""} apiKey={createdKey ? formatKeyForDisplay(createdKey) : KEY_PLACEHOLDER} />
      </CardContent>
    </Card>
  );
}

/** The plan benefit: every plan includes the same ten models. */
const planModels = [
  "deepseek-v4.1-flash-fast",
  "deepseek-v4-flash-fast",
  "minimax-m3",
  "mimo-v2.6-pro",
  "mimo-v2.6-flash",
  "kimi-k3",
  "glm-5.3",
  "glm-5.3-flash",
  "deepseek-v4.1-flash",
  "deepseek-v4-flash",
];

export function PlanBenefitsCard({ credits }: { credits: string }) {
  return (
    <Card className="border-line-warm bg-white">
      <CardHeader>
        <CardTitle className="text-[15px]">Plan benefits</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4 text-[12px]">
        <div>
          <p className="font-medium text-muted-warm">Models ({planModels.length})</p>
          <ul className="mt-2 flex flex-wrap gap-1.5">
            {planModels.map((model) => (
              <li
                key={model}
                className="rounded-full bg-cream-card px-2.5 py-1 font-mono text-[10px] text-ink"
              >
                {model}
              </li>
            ))}
          </ul>
        </div>
        <div>
          <p className="font-medium text-muted-warm">Credits</p>
          <p className="mt-1 text-ink">{credits}</p>
        </div>
        <div>
          <p className="font-medium text-muted-warm">Coding tools</p>
          <p className="mt-1 text-ink">
            OpenClaw, Codex, Claude Code, MiMo Code and others
          </p>
        </div>
      </CardContent>
    </Card>
  );
}

/** Token pricing per model, in credits per token. Transcribed from the platform's rate card. */
export function TokenPricingCard() {
  const columns = [
    "Model",
    "Cached input",
    "Input",
    "Output",
    "Cache write 5m",
    "Cache write 1h",
  ];
  const credits = (value: number | null) =>
    value === null ? <span className="text-muted-warm">n.a.</span> : value.toLocaleString();

  return (
    <Card className="border-line-warm bg-white">
      <CardHeader>
        <CardTitle className="text-[15px]">Token pricing</CardTitle>
        <p className="mt-1 text-[11px] text-muted-warm">
          Credits per token. Cache write applies only to models that support it.
        </p>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto rounded-xl border border-line-warm">
          <Table className="min-w-[720px]">
            <TableHeader>
              <TableRow className="bg-cream-card">
                {columns.map((column, i) => (
                  <TableHead
                    key={column}
                    className={`text-[10px] uppercase tracking-wide text-muted-warm ${
                      i === 0 ? "px-4 py-3" : "px-4 py-3 text-right"
                    }`}
                  >
                    {column}
                  </TableHead>
                ))}
              </TableRow>
            </TableHeader>
            <TableBody>
              {MODEL_PRICING.map((row) => (
                <TableRow key={row.model} className="border-t border-line-warm">
                  <TableCell className="px-4 py-3 text-[12px] font-medium text-ink">
                    {row.model}
                  </TableCell>
                  <TableCell className="px-4 py-3 text-right font-mono text-[12px] text-ink">
                    {credits(row.cachedInput)}
                  </TableCell>
                  <TableCell className="px-4 py-3 text-right font-mono text-[12px] text-ink">
                    {credits(row.input)}
                  </TableCell>
                  <TableCell className="px-4 py-3 text-right font-mono text-[12px] text-ink">
                    {credits(row.output)}
                  </TableCell>
                  <TableCell className="px-4 py-3 text-right font-mono text-[12px] text-ink">
                    {credits(row.cacheWrite5m)}
                  </TableCell>
                  <TableCell className="px-4 py-3 text-right font-mono text-[12px] text-ink">
                    {credits(row.cacheWrite1h)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}

/** Quick integration guide: copy-paste snippets plus per-tool setup. */
export function QuickIntegrationDialog({ baseURL, apiKey }: { baseURL: string; apiKey: string }) {
  const samples = codeSamples(baseURL);
  const tools = toolGuides(baseURL, apiKey);

  return (
    <Dialog>
      <DialogTrigger
        render={
          <Button variant="ghost" className="gap-1 px-0 text-[11px] text-flame">
            Quick integration <ExternalLink className="size-3" />
          </Button>
        }
      />
      <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-[720px]">
        <DialogHeader>
          <DialogTitle className="text-[15px]">Quick integration</DialogTitle>
          <DialogDescription className="text-[12px] text-muted-warm">
            One base URL serves both the OpenAI and Anthropic styles. Send the key as a
            bearer token.
          </DialogDescription>
        </DialogHeader>

        {baseURL ? (
          <div className="rounded-xl border border-line-warm bg-cream-card px-4 py-3">
            <p className="text-[10px] uppercase tracking-wide text-muted-warm">Base URL</p>
            <div className="mt-1 flex items-center gap-2">
              <code className="flex-1 break-all font-mono text-[12px] text-ink">{baseURL}</code>
              <CopyButton value={baseURL} />
            </div>
          </div>
        ) : (
          <div className="rounded-xl border border-line-warm px-4 py-3 text-[12px] text-muted-warm">
            No base URL is configured for this platform yet.
          </div>
        )}

        <div className="space-y-4">
          {samples.map((sample) => (
            <div key={sample.id}>
              <div className="flex items-center justify-between">
                <p className="text-[11px] font-medium text-ink">{sample.label}</p>
                <CopyButton value={sample.code} />
              </div>
              <pre className="mt-2 overflow-x-auto rounded-xl border border-line-warm bg-cream-card p-3">
                <code className="font-mono text-[11px] leading-relaxed text-ink">
                  {sample.code}
                </code>
              </pre>
            </div>
          ))}
        </div>

        <div className="border-t border-line-warm pt-4">
          <p className="text-[13px] font-semibold text-ink">Tool setup</p>
          <p className="mt-1 text-[11px] text-muted-warm">
            Each tool takes the same base URL and API key.
          </p>
        </div>

        <div className="space-y-5">
          {tools.map((tool) => (
            <div key={tool.id}>
              <p className="text-[12px] font-medium text-ink">{tool.name}</p>
              <p className="mt-0.5 text-[11px] text-muted-warm">{tool.blurb}</p>
              <ol className="mt-2 list-inside list-decimal space-y-1 text-[11px] text-muted-warm">
                {tool.steps.map((step) => (
                  <li key={step}>{step}</li>
                ))}
              </ol>
              <div className="mt-2">
                <div className="flex items-center justify-between">
                  <p className="text-[10px] uppercase tracking-wide text-muted-warm">
                    {tool.configLabel}
                  </p>
                  <CopyButton value={tool.config} />
                </div>
                <pre className="mt-1 overflow-x-auto rounded-xl border border-line-warm bg-cream-card p-3">
                  <code className="font-mono text-[11px] leading-relaxed text-ink">
                    {tool.config}
                  </code>
                </pre>
              </div>
            </div>
          ))}
        </div>
      </DialogContent>
    </Dialog>
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