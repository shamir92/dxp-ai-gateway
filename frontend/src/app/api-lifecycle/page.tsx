import { Plus, Wallet } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { AppShell } from "@/components/app-shell";
import { StatTile, TopSearchBar } from "@/components/token-plan-ui";

const steps = [
  { n: "01", title: "API keys", desc: "Authenticate requests" },
  { n: "02", title: "Balance", desc: "Fund consumption" },
  { n: "03", title: "Batch jobs", desc: "Process workloads" },
  { n: "04", title: "Billing", desc: "Review and reconcile" },
];

const keys = [
  { name: "Production", key: "sk_live_••••8K2p", used: "2 min ago" },
  { name: "Batch worker", key: "sk_live_••••2Xf9", used: "18 min ago" },
  { name: "Development", key: "sk_test_••••4Nh7", used: "3 hours ago" },
];

const jobs = [
  ["embed-docs-01", "mimo-v2.6-flash", "12,400", "$42.10", "Running"],
  ["eval-suite-04", "mimo-v2.6-pro", "3,200", "$28.60", "Running"],
  ["summarize-batch", "mimo-v2.5", "48,000", "$61.20", "Queued"],
  ["classify-batch", "mimo-v2.6-flash", "22,000", "$19.40", "Completed"],
];

export default function ApiLifecyclePage() {
  return (
    <AppShell
      title="API lifecycle"
      subtitle="Create credentials, fund usage, run batch workloads, and reconcile billing in one flow."
      actions={<TopSearchBar />}
    >
      <div className="mx-auto max-w-[1204px] space-y-8">
        <div className="flex items-center gap-4">
          <Button variant="ghost" className="rounded-full border border-line-warm text-[12px]">
            Lifecycle guide
          </Button>
          <div className="flex flex-1 flex-wrap gap-3">
            {steps.map((s) => (
              <div key={s.n} className="flex min-w-[160px] flex-1 items-center gap-3 rounded-xl border border-line-warm bg-white px-4 py-3">
                <span className="text-[12px] font-semibold text-flame">{s.n}</span>
                <div>
                  <p className="text-[12px] font-semibold text-ink">{s.title}</p>
                  <p className="text-[10px] text-muted-warm">{s.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <Card className="border-line-warm bg-white">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-[15px]">1. API keys & endpoints</CardTitle>
              <p className="mt-1 text-[12px] text-muted-warm">
                Credentials that initiate metered gateway usage.
              </p>
            </div>
            <Button className="rounded-full bg-ink text-[12px] text-white hover:bg-ink/90">
              <Plus className="mr-1 size-3.5" /> Create API key
            </Button>
          </CardHeader>
          <CardContent>
            <div className="overflow-hidden rounded-xl border border-line-warm">
              <table className="w-full text-left text-[12px]">
                <thead className="bg-cream-card text-[10px] uppercase tracking-wide text-muted-warm">
                  <tr>
                    <th className="px-4 py-3">Name</th>
                    <th className="px-4 py-3">Key</th>
                    <th className="px-4 py-3">Last used</th>
                    <th className="px-4 py-3">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {keys.map((k) => (
                    <tr key={k.name} className="border-t border-line-warm">
                      <td className="px-4 py-3 font-medium text-ink">{k.name}</td>
                      <td className="px-4 py-3 font-mono text-muted-warm">{k.key}</td>
                      <td className="px-4 py-3 text-muted-warm">{k.used}</td>
                      <td className="px-4 py-3">
                        <Badge className="rounded-full bg-success/10 text-success">Active</Badge>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <div className="mt-5 rounded-xl border border-line-warm p-4">
              <p className="text-[11px] font-medium text-muted-warm">Gateway endpoint</p>
              <p className="mt-1 text-[11px] text-muted-warm">
                Use with any OpenAI-compatible SDK.
              </p>
              <code className="mt-2 block font-mono text-[13px] text-ink">
                https://api.dxp.ai/v1
              </code>
              <p className="mt-2 text-[11px] text-muted-warm">Region · Singapore (sgp-1)</p>
            </div>
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-[15px]">2. Balance & funding</CardTitle>
              <p className="mt-1 text-[12px] text-muted-warm">
                Prepaid funds power real-time and batch inference.
              </p>
            </div>
            <Button className="rounded-full bg-ink text-[12px] text-white hover:bg-ink/90">
              <Wallet className="mr-1 size-3.5" /> Recharge
            </Button>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-3">
              <StatTile label="Available balance" value="$936.00" hint="Cash $900 · Bonus $36" />
              <StatTile label="Alert threshold" value="$200.00" hint="Notifications enabled" />
              <StatTile label="30-day burn" value="$312.40" hint="10.4% lower than Sep" />
            </div>
            <div className="mt-5">
              <p className="text-[11px] font-medium text-muted-warm">Quick recharge</p>
              <div className="mt-2 flex flex-wrap gap-2">
                {["$50", "$100", "$200", "$500", "Custom"].map((a, i) => (
                  <Button
                    key={a}
                    variant="outline"
                    className={`rounded-full border-line-warm text-[12px] ${i === 1 ? "border-flame text-flame" : ""}`}
                  >
                    {a}
                  </Button>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-[15px]">3. Batch jobs</CardTitle>
              <p className="mt-1 text-[12px] text-muted-warm">
                Jobs consume funded balance and generate billable model usage.
              </p>
            </div>
            <Button className="rounded-full bg-ink text-[12px] text-white hover:bg-ink/90">
              <Plus className="mr-1 size-3.5" /> New batch job
            </Button>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-4">
              <StatTile label="Running" value="2" hint="Across 18.4K requests" />
              <StatTile label="Queued" value="3" hint="Est. start in 8 min" />
              <StatTile label="Completed" value="128" hint="99.2% success rate" />
              <StatTile label="Month cost" value="$184.60" hint="59% of total spend" />
            </div>
            <div className="mt-5 overflow-hidden rounded-xl border border-line-warm">
              <table className="w-full text-left text-[12px]">
                <thead className="bg-cream-card text-[10px] uppercase tracking-wide text-muted-warm">
                  <tr>
                    <th className="px-4 py-3">Job</th>
                    <th className="px-4 py-3">Model</th>
                    <th className="px-4 py-3">Requests</th>
                    <th className="px-4 py-3">Cost</th>
                    <th className="px-4 py-3">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {jobs.map((j) => (
                    <tr key={j[0]} className="border-t border-line-warm">
                      <td className="px-4 py-3 font-medium text-ink">{j[0]}</td>
                      <td className="px-4 py-3 font-mono text-muted-warm">{j[1]}</td>
                      <td className="px-4 py-3 text-muted-warm">{j[2]}</td>
                      <td className="px-4 py-3 text-muted-warm">{j[3]}</td>
                      <td className="px-4 py-3">
                        <Badge
                          className={`rounded-full ${
                            j[4] === "Completed"
                              ? "bg-success/10 text-success"
                              : j[4] === "Queued"
                                ? "bg-cream text-muted-warm"
                                : "bg-flame/10 text-flame"
                          }`}
                        >
                          {j[4]}
                        </Badge>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </div>
    </AppShell>
  );
}
