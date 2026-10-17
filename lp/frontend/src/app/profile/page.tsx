"use client";

import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import { AppShell } from "@/components/app-shell";
import { TopSearchBar } from "@/components/token-plan-ui";
import { apiJSON } from "@/lib/api";

type ProfileUser = {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
};

type Connection = { provider: string; connected: boolean; email: string };
type SessionRow = {
  id: string;
  device: string;
  user_agent: string;
  ip: string;
  location: string;
  last_seen_at: string;
  current: boolean;
};

type Settings = {
  notify_balance: boolean;
  notify_usage: boolean;
  notify_security: boolean;
  notify_product: boolean;
  login_alerts: boolean;
};

type ProfileResponse = {
  user: ProfileUser;
  connections: Connection[];
  session_count: number;
  settings?: { login_alerts: boolean };
};

function initials(name: string, email: string) {
  const src = name || email || "?";
  return src
    .split(/[\s@._-]+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((p) => p[0]!.toUpperCase())
    .join("");
}

export default function ProfilePage() {
  const [user, setUser] = useState<ProfileUser | null>(null);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [connections, setConnections] = useState<Connection[]>([]);
  const [sessions, setSessions] = useState<SessionRow[]>([]);
  const [settings, setSettings] = useState<Settings>({
    notify_balance: true,
    notify_usage: true,
    notify_security: true,
    notify_product: false,
    login_alerts: true,
  });
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    try {
      const [profile, preferences, sess] = await Promise.all([
        apiJSON<ProfileResponse>("/api/v1/profile"),
        apiJSON<Omit<Settings, "login_alerts">>("/api/v1/settings"),
        apiJSON<{ sessions: SessionRow[] }>("/api/v1/profile/sessions"),
      ]);
      setUser(profile.user);
      setName(profile.user.name);
      setEmail(profile.user.email);
      setConnections(profile.connections ?? []);
      setSettings({
        notify_balance: preferences.notify_balance ?? true,
        notify_usage: preferences.notify_usage ?? true,
        notify_security: preferences.notify_security ?? true,
        notify_product: preferences.notify_product ?? false,
        login_alerts: profile.settings?.login_alerts ?? true,
      });
      setSessions(sess.sessions ?? []);
    } catch {
      window.location.href = "/login";
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  function patchSettings(p: Partial<Settings>) {
    setSettings((prev) => ({ ...prev, ...p }));
    setSaved(false);
  }

  async function save() {
    setBusy(true);
    setError("");
    setSaved(false);
    try {
      await apiJSON("/api/v1/profile", {
        method: "PATCH",
        body: JSON.stringify({ name, email }),
      });
      await apiJSON("/api/v1/settings", {
        method: "PATCH",
        body: JSON.stringify({
          notify_balance: settings.notify_balance,
          notify_usage: settings.notify_usage,
          notify_security: settings.notify_security,
          notify_product: settings.notify_product,
        }),
      });
      await apiJSON("/api/v1/profile/security", {
        method: "PATCH",
        body: JSON.stringify({ login_alerts: settings.login_alerts }),
      });
      setSaved(true);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Save failed");
    } finally {
      setBusy(false);
    }
  }

  async function revokeSession(id: string) {
    await apiJSON(`/api/v1/profile/sessions/${id}`, { method: "DELETE" }).catch(() => {});
    await load();
  }

  async function revokeAll() {
    await apiJSON("/api/v1/profile/sessions/revoke-all", { method: "POST" }).catch(() => {});
    window.location.href = "/login";
  }

  async function unlink(provider: string) {
    await apiJSON(`/api/v1/profile/connections/${provider}`, { method: "DELETE" }).catch((e) => {
      setError(e instanceof Error ? e.message : "Unlink failed");
    });
    await load();
  }

  return (
    <AppShell
      title="Profile & account"
      subtitle="Manage your identity, connected accounts, security, and notifications."
      actions={
        <div className="flex items-center gap-3">
          <TopSearchBar />
          <Button
            onClick={() => void save()}
            disabled={busy}
            className="rounded-full bg-ink text-[12px] text-white hover:bg-ink/90"
          >
            Save changes
          </Button>
        </div>
      }
    >
      <div className="mx-auto max-w-[860px] space-y-6">
        {error ? (
          <p className="rounded-xl border border-flame/30 bg-flame/5 px-4 py-3 text-[12px] text-flame">{error}</p>
        ) : null}
        {saved ? (
          <p className="rounded-xl border border-success/30 bg-success/5 px-4 py-3 text-[12px] text-success">
            Saved.
          </p>
        ) : null}

        <Card className="border-line-warm bg-white">
          <CardContent className="flex flex-wrap items-center gap-5 py-6">
            <div className="flex size-16 items-center justify-center rounded-2xl bg-prompt text-[18px] font-semibold text-white">
              {user ? initials(user.name, user.email) : "··"}
            </div>
            <div className="flex-1">
              <p className="text-[18px] font-semibold text-ink">{user?.name || "—"}</p>
              <p className="text-[12px] text-muted-warm">Workspace owner</p>
              <p className="text-[12px] text-muted-warm">{user?.email}</p>
            </div>
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">Account information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label htmlFor="name">Display name</Label>
                <Input id="name" value={name} onChange={(e) => setName(e.target.value)} />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="email">Primary email</Label>
                <Input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
              </div>
            </div>
            <p className="text-[11px] text-muted-warm">
              Your primary email is used for security alerts, invoices, and account recovery.
            </p>
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">Connected accounts</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {connections.map((c) => (
              <div
                key={c.provider}
                className="flex flex-wrap items-center gap-3 rounded-xl border border-line-warm px-4 py-3"
              >
                <div className="flex-1">
                  <p className="text-[13px] font-medium text-ink capitalize">{c.provider}</p>
                  <p className="text-[11px] text-muted-warm">
                    {c.connected
                      ? c.provider === "email"
                        ? `Connected · ${c.email || email || "magic code"}`
                        : `Connected · ${c.email || "linked"}`
                      : "Not connected"}
                  </p>
                </div>
                {c.connected ? (
                  <Badge className="rounded-full bg-success/10 text-success">Connected</Badge>
                ) : null}
                {c.provider !== "email" ? (
                  <Button
                    variant="outline"
                    className="rounded-full border-line-warm text-[12px]"
                    onClick={() => {
                      if (c.connected) void unlink(c.provider);
                      else window.location.href = `/api/v1/profile/connections/${c.provider}/link`;
                    }}
                  >
                    {c.connected ? "Manage" : "Connect"}
                  </Button>
                ) : (
                  <span className="rounded-full border border-line-warm px-3 py-1 text-[11px] text-muted-warm">
                    Magic code
                  </span>
                )}
              </div>
            ))}
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">Security</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex flex-wrap items-center gap-3">
              <div className="flex-1">
                <p className="text-[13px] font-medium text-ink">Login alerts</p>
                <p className="text-[11px] text-muted-warm">Email on new device sign-in</p>
              </div>
              <Switch
                checked={settings.login_alerts}
                onCheckedChange={(v) => patchSettings({ login_alerts: v })}
              />
            </div>
            <Separator />
            <div>
              <div className="flex items-center justify-between">
                <p className="text-[13px] font-medium text-ink">Active sessions</p>
                <Button variant="ghost" className="text-[12px] text-flame" onClick={() => void revokeAll()}>
                  Sign out all
                </Button>
              </div>
              <div className="mt-3 space-y-2">
                {sessions.map((s) => (
                  <div
                    key={s.id}
                    className="flex items-center gap-3 rounded-xl border border-line-warm px-4 py-3"
                  >
                    <div className="flex-1">
                      <p className="text-[12px] font-medium text-ink">
                        {s.device || "Unknown device"}
                        {s.current ? " · Current session" : ""}
                      </p>
                      <p className="text-[11px] text-muted-warm">
                        {[s.location, s.ip].filter(Boolean).join(" · ") || "—"}
                      </p>
                    </div>
                    {!s.current ? (
                      <Button
                        variant="ghost"
                        className="text-[11px] text-flame"
                        onClick={() => void revokeSession(s.id)}
                      >
                        Revoke
                      </Button>
                    ) : null}
                  </div>
                ))}
                {sessions.length === 0 ? (
                  <p className="text-[12px] text-muted-warm">No active sessions.</p>
                ) : null}
              </div>
            </div>
          </CardContent>
        </Card>

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
                  checked={settings[key]}
                  onCheckedChange={(v) => patchSettings({ [key]: v } as Partial<Settings>)}
                />
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </AppShell>
  );
}
