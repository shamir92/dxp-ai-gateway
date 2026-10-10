"use client";

import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { AppShell } from "@/components/app-shell";
import { TopSearchBar } from "@/components/token-plan-ui";
import { apiJSON } from "@/lib/api";

type Settings = {
  notify_balance: boolean;
  notify_usage: boolean;
  notify_security: boolean;
  notify_product: boolean;
};

export default function SettingsPage() {
  const [s, setS] = useState<Settings | null>(null);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    try {
      const data = await apiJSON<Settings>("/api/v1/settings");
      setS({
        notify_balance: data.notify_balance,
        notify_usage: data.notify_usage,
        notify_security: data.notify_security,
        notify_product: data.notify_product,
      });
    } catch {
      window.location.href = "/login";
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  function patch(p: Partial<Settings>) {
    setS((prev) => (prev ? { ...prev, ...p } : prev));
    setSaved(false);
  }

  async function save() {
    if (!s) return;
    setBusy(true);
    setError("");
    try {
      await apiJSON("/api/v1/settings", {
        method: "PATCH",
        body: JSON.stringify(s),
      });
      setSaved(true);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Save failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <AppShell
      title="Settings"
      subtitle="Workspace preferences and notifications."
      actions={
        <div className="flex items-center gap-3">
          <TopSearchBar />
          <Button
            onClick={() => void save()}
            disabled={busy || !s}
            className="rounded-full bg-ink text-[12px] text-white hover:bg-ink/90"
          >
            Save changes
          </Button>
        </div>
      }
    >
      <div className="mx-auto max-w-[760px] space-y-6">
        {error ? (
          <p className="rounded-xl border border-flame/30 bg-flame/5 px-4 py-3 text-[12px] text-flame">{error}</p>
        ) : null}
        {saved ? (
          <p className="rounded-xl border border-success/30 bg-success/5 px-4 py-3 text-[12px] text-success">
            Settings saved.
          </p>
        ) : null}

        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">Notifications</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {(
              [
                ["notify_balance", "Balance alerts", "Email when balance falls below threshold"],
                ["notify_usage", "Usage digest", "Weekly usage and spend summary"],
                ["notify_security", "Security alerts", "New device and credential changes"],
                ["notify_product", "Product updates", "Model launches and API changes"],
              ] as const
            ).map(([key, title, desc]) => (
              <div key={key} className="flex items-center justify-between gap-4">
                <div>
                  <p className="text-[13px] font-medium text-ink">{title}</p>
                  <p className="text-[11px] text-muted-warm">{desc}</p>
                </div>
                <Switch
                  checked={s?.[key] ?? false}
                  onCheckedChange={(v) => patch({ [key]: v } as Partial<Settings>)}
                />
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </AppShell>
  );
}
