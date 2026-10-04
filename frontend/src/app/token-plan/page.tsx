"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { AppShell } from "@/components/app-shell";
import {
  ApiKeyCard,
  PlanBenefitsCard,
  PlanMetaStrip,
  PlanUsageCard,
  TopSearchBar,
  UsageBarChart,
  UsageLegend,
} from "@/components/token-plan-ui";

const dayLabels = [
  "Oct 01",
  "Oct 02",
  "Oct 03",
  "Oct 04",
  "Oct 05",
  "Oct 06",
  "Oct 07",
  "Oct 08",
  "Oct 09",
  "Oct 10",
];

const totalSeries = [
  { name: "mimo-v2.6-flash", color: "#FF5A2A", values: [420, 510, 480, 620, 590, 700, 650, 720, 680, 740] },
  { name: "mimo-v2.6-pro", color: "#FF9B86", values: [180, 220, 200, 260, 240, 300, 280, 310, 290, 320] },
  { name: "mimo-v2.5-pro", color: "#141F32", values: [90, 110, 100, 130, 120, 150, 140, 160, 150, 170] },
  { name: "mimo-v2.5", color: "#6D6962", values: [40, 50, 45, 60, 55, 70, 65, 75, 70, 80] },
];

const singleSeries = [
  { name: "Total", color: "#FF5A2A", values: [120, 180, 150, 220, 200, 260, 240, 280, 250, 300] },
  { name: "Cache hit", color: "#38A169", values: [70, 100, 90, 130, 110, 150, 140, 160, 145, 170] },
  { name: "Input", color: "#FF9B86", values: [30, 50, 40, 60, 55, 70, 65, 75, 70, 80] },
  { name: "Output", color: "#141F32", values: [20, 30, 20, 30, 35, 40, 35, 45, 35, 50] },
];

export default function TokenPlanPage() {
  const [view, setView] = useState("single");

  return (
    <AppShell
      title="Token plan"
      subtitle="Manage your subscription, credentials, and detailed model usage."
      actions={<TopSearchBar />}
    >
      <div className="mx-auto max-w-[1204px] space-y-6">
        <PlanMetaStrip />
        <PlanUsageCard />

        <div className="grid gap-6 lg:grid-cols-2">
          <ApiKeyCard />
          <PlanBenefitsCard />
        </div>

        <Card className="border-line-warm bg-white">
          <CardHeader className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <CardTitle className="text-[15px]">Usage details</CardTitle>
              <p className="mt-1 text-[11px] text-muted-warm">
                All timestamps are displayed in UTC. Data updates in near-real time (≤5 min
                delay).
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              <Select defaultValue="2026-10">
                <SelectTrigger className="w-[120px] rounded-full text-[12px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="2026-10">Month 2026-10</SelectItem>
                  <SelectItem value="2026-09">Month 2026-09</SelectItem>
                </SelectContent>
              </Select>
              <Select defaultValue="all">
                <SelectTrigger className="w-[130px] rounded-full text-[12px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All models</SelectItem>
                  <SelectItem value="flash">mimo-v2.6-flash</SelectItem>
                  <SelectItem value="pro">mimo-v2.6-pro</SelectItem>
                </SelectContent>
              </Select>
              <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
                Export
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap items-center justify-between gap-3">
              <p className="text-[12px] text-muted-warm">
                {view === "single"
                  ? "mimo-v2.6-flash · 8.4B tokens · Oct 01–30"
                  : "Total token consumption · 178,381,464 tokens · Oct 01–30"}
              </p>
              <Tabs value={view} onValueChange={setView}>
                <TabsList className="rounded-full bg-cream p-1">
                  <TabsTrigger value="single" className="rounded-full text-[11px]">
                    Single model token consumption
                  </TabsTrigger>
                  <TabsTrigger value="total" className="rounded-full text-[11px]">
                    Total token consumption
                  </TabsTrigger>
                </TabsList>
              </Tabs>
            </div>

            {view === "single" ? (
              <>
                <div className="mt-4 flex flex-wrap items-center gap-3">
                  <span className="text-[10px] font-medium text-muted-warm">
                    SELECTED MODEL
                  </span>
                  <span className="rounded-full bg-cream px-3 py-1 font-mono text-[12px] text-ink">
                    mimo-v2.6-flash
                  </span>
                </div>
                <UsageLegend
                  items={[
                    { name: "Total tokens", value: "8.4B", share: "100%", color: "#FF5A2A" },
                    { name: "Cache hit", value: "4.7B", share: "56%", color: "#38A169" },
                    { name: "Input", value: "2.4B", share: "29%", color: "#FF9B86" },
                    { name: "Output", value: "1.3B", share: "15%", color: "#141F32" },
                  ]}
                />
                <UsageBarChart series={singleSeries} labels={dayLabels} max={350} />
              </>
            ) : (
              <>
                <UsageLegend
                  items={[
                    { name: "mimo-v2.6-flash", value: "8.4B", share: "47%", color: "#FF5A2A" },
                    { name: "mimo-v2.6-pro", value: "5.1B", share: "29%", color: "#FF9B86" },
                    { name: "mimo-v2.5-pro", value: "2.8B", share: "16%", color: "#141F32" },
                    { name: "mimo-v2.5", value: "1.5B", share: "8%", color: "#6D6962" },
                  ]}
                />
                <UsageBarChart series={totalSeries} labels={dayLabels} max={800} />
                <div className="mt-8">
                  <p className="text-[12px] text-muted-warm">
                    Number of requests · 1,774 requests · 30-day period
                  </p>
                  <div className="mt-3 grid grid-cols-3 gap-3">
                    {[
                      ["mimo-v2.6-flash", "1,024"],
                      ["mimo-v2.6-pro", "612"],
                      ["Other", "138"],
                    ].map(([k, v]) => (
                      <div key={k} className="rounded-xl border border-line-warm p-3">
                        <p className="text-[11px] text-muted-warm">{k}</p>
                        <p className="mt-1 text-[16px] font-semibold text-ink">{v}</p>
                      </div>
                    ))}
                  </div>
                </div>
              </>
            )}
          </CardContent>
        </Card>
      </div>
    </AppShell>
  );
}
