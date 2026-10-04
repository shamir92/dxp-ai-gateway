import { ExternalLink, Puzzle, Search, Shield, Zap } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AppShell } from "@/components/app-shell";
import { StatTile, TopSearchBar } from "@/components/token-plan-ui";

const plugins = [
  {
    name: "Web search",
    status: "Enabled",
    desc: "Search and parse current public web content for grounded answers.",
    models: "All chat models",
    billing: "Usage based",
    action: "Configure",
  },
  {
    name: "Document fetch",
    status: "Beta",
    desc: "Retrieve and extract approved URLs, PDFs, and documents.",
    models: "Pro models",
    billing: "Free in beta",
    action: "Enable",
  },
  {
    name: "Code interpreter",
    status: "Disabled",
    desc: "Run isolated Python analysis, transformations, and charts.",
    models: "Pro models",
    billing: "$0.02 per run",
    action: "Enable",
  },
  {
    name: "Knowledge base",
    status: "Disabled",
    desc: "Search private indexed content with workspace permissions.",
    models: "All chat models",
    billing: "$0.001 per query",
    action: "Enable",
  },
];

export default function PluginsPage() {
  return (
    <AppShell
      title="Plugin management"
      subtitle="Extend model capabilities with controlled, billable tools."
      actions={
        <div className="flex items-center gap-3">
          <TopSearchBar />
          <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
            Plugin documentation
          </Button>
        </div>
      }
    >
      <div className="mx-auto max-w-[1204px] space-y-6">
        <Card className="border-line-warm bg-white">
          <CardContent className="flex flex-wrap items-center gap-6 py-6">
            <div className="flex size-12 items-center justify-center rounded-2xl bg-cream">
              <Puzzle className="size-5 text-flame" />
            </div>
            <div className="flex-1">
              <p className="text-[11px] font-medium uppercase tracking-wide text-muted-warm">
                Connected services
              </p>
              <h2 className="mt-1 text-[18px] font-semibold text-ink">
                Bring trusted tools into every request
              </h2>
              <p className="mt-1 text-[12px] text-muted-warm">
                Enable only the capabilities approved for your workspace.
              </p>
            </div>
            <StatTile label="Enabled" value="1 / 4" hint="Web search active" />
            <StatTile label="Month spend" value="$18.40" />
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader className="flex flex-wrap items-center justify-between gap-3">
            <CardTitle className="text-[15px]">Plugin services</CardTitle>
            <Tabs defaultValue="available">
              <TabsList className="rounded-full bg-cream p-1">
                <TabsTrigger value="available" className="rounded-full text-[11px]">
                  Available plugins
                </TabsTrigger>
                <TabsTrigger value="enabled" className="rounded-full text-[11px]">
                  Enabled
                </TabsTrigger>
                <TabsTrigger value="activity" className="rounded-full text-[11px]">
                  Activity log
                </TabsTrigger>
              </TabsList>
            </Tabs>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex flex-wrap gap-2">
              <div className="flex min-w-[220px] flex-1 items-center gap-2 rounded-full border border-line-warm px-3">
                <Search className="size-3.5 text-muted-warm" />
                <Input
                  placeholder="Search plugins"
                  className="h-9 border-0 bg-transparent p-0 text-[12px] shadow-none focus-visible:ring-0"
                />
              </div>
              <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
                All compatible models
              </Button>
            </div>

            {plugins.map((p) => (
              <div
                key={p.name}
                className="flex flex-wrap items-center gap-4 rounded-xl border border-line-warm p-4"
              >
                <div className="flex size-10 items-center justify-center rounded-xl bg-cream">
                  {p.name === "Web search" ? (
                    <Search className="size-4 text-flame" />
                  ) : p.name === "Code interpreter" ? (
                    <Zap className="size-4 text-flame" />
                  ) : p.name === "Knowledge base" ? (
                    <Shield className="size-4 text-flame" />
                  ) : (
                    <Puzzle className="size-4 text-flame" />
                  )}
                </div>
                <div className="min-w-[220px] flex-1">
                  <div className="flex items-center gap-2">
                    <p className="text-[13px] font-semibold text-ink">{p.name}</p>
                    <Badge
                      className={`rounded-full text-[10px] ${
                        p.status === "Enabled"
                          ? "bg-success/10 text-success"
                          : p.status === "Beta"
                            ? "bg-flame/10 text-flame"
                            : "bg-cream text-muted-warm"
                      }`}
                    >
                      {p.status}
                    </Badge>
                  </div>
                  <p className="mt-1 text-[12px] text-muted-warm">{p.desc}</p>
                  <p className="mt-1 text-[11px] text-muted-warm">
                    Models · {p.models} · Billing · {p.billing}
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="ghost" className="gap-1 text-[11px] text-flame">
                    API docs <ExternalLink className="size-3" />
                  </Button>
                  <Button
                    variant="outline"
                    className="rounded-full border-line-warm text-[12px]"
                  >
                    {p.action}
                  </Button>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <div className="grid gap-6 lg:grid-cols-2">
          <Card className="border-line-warm bg-white">
            <CardHeader>
              <CardTitle className="text-[15px]">Web search</CardTitle>
              <p className="text-[12px] text-muted-warm">
                Real-time search and page parsing with domain controls and transparent
                usage pricing.
              </p>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-3">
                <StatTile label="Requests this month" value="920" />
                <StatTile label="Successful calls" value="98.7%" />
                <StatTile label="Current spend" value="$18.40" />
                <StatTile label="Last used" value="2 minutes ago" />
              </div>
            </CardContent>
          </Card>
          <Card className="border-line-warm bg-white">
            <CardHeader>
              <CardTitle className="text-[15px]">Permissions</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {["Public web access", "Page content parsing", "Usage logging"].map((perm) => (
                <div
                  key={perm}
                  className="flex items-center justify-between rounded-xl border border-line-warm px-4 py-3 text-[12px] text-ink"
                >
                  {perm}
                  <Badge className="rounded-full bg-success/10 text-success">Allowed</Badge>
                </div>
              ))}
              <Button variant="outline" className="w-full rounded-full border-destructive/30 text-destructive">
                Disable plugin
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </AppShell>
  );
}
