import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { AppShell } from "@/components/app-shell";
import { TopSearchBar } from "@/components/token-plan-ui";

export default function ProfilePage() {
  return (
    <AppShell
      title="Profile & account"
      subtitle="Manage your identity, connected accounts, and security preferences."
      actions={
        <div className="flex items-center gap-3">
          <TopSearchBar />
          <Button className="rounded-full bg-ink text-[12px] text-white hover:bg-ink/90">
            Save changes
          </Button>
        </div>
      }
    >
      <div className="mx-auto max-w-[860px] space-y-6">
        <Card className="border-line-warm bg-white">
          <CardContent className="flex flex-wrap items-center gap-5 py-6">
            <div className="flex size-16 items-center justify-center rounded-2xl bg-prompt text-[18px] font-semibold text-white">
              SK
            </div>
            <div className="flex-1">
              <p className="text-[18px] font-semibold text-ink">Shamir K.</p>
              <p className="text-[12px] text-muted-warm">Workspace owner</p>
              <p className="text-[12px] text-muted-warm">sha•••@f•••••p.com</p>
              <p className="text-[11px] text-muted-warm">Xiaomi ID · 6892304915</p>
            </div>
            <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
              Change photo
            </Button>
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
                <Input id="name" defaultValue="Shamir K." />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="email">Primary email</Label>
                <Input id="email" defaultValue="sha•••@f•••••p.com" />
              </div>
              <div className="space-y-1.5">
                <Label>Phone number</Label>
                <div className="flex items-center gap-2">
                  <Input value="Not connected" readOnly />
                  <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
                    Connect
                  </Button>
                </div>
              </div>
              <div className="space-y-1.5">
                <Label>Region</Label>
                <Input defaultValue="Singapore (SG)" readOnly />
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
            {[
              ["Google", "Connected · sha•••@gmail.com", "Manage", true],
              ["GitHub", "Not connected", "Connect", false],
              ["Open Platform", "Primary email not bound", "Connect", false],
            ].map(([name, status, action, connected]) => (
              <div key={name as string} className="flex flex-wrap items-center gap-3 rounded-xl border border-line-warm px-4 py-3">
                <div className="flex-1">
                  <p className="text-[13px] font-medium text-ink">{name}</p>
                  <p className="text-[11px] text-muted-warm">{status}</p>
                </div>
                {connected ? (
                  <Badge className="rounded-full bg-success/10 text-success">Connected</Badge>
                ) : null}
                <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
                  {action}
                </Button>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">Security</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {[
              ["Two-factor authentication", "Authenticator app enabled", "Enabled"],
              ["Login alerts", "Email on new device sign-in", "Enabled"],
              ["Password", "Last changed 34 days ago", "Update"],
            ].map(([title, desc, action]) => (
              <div key={title as string} className="flex flex-wrap items-center gap-3">
                <div className="flex-1">
                  <p className="text-[13px] font-medium text-ink">{title}</p>
                  <p className="text-[11px] text-muted-warm">{desc}</p>
                </div>
                <Button variant="outline" className="rounded-full border-line-warm text-[12px]">
                  {action}
                </Button>
              </div>
            ))}
            <Separator />
            <div>
              <div className="flex items-center justify-between">
                <p className="text-[13px] font-medium text-ink">Active sessions</p>
                <Button variant="ghost" className="text-[12px] text-flame">
                  Sign out all
                </Button>
              </div>
              <div className="mt-3 space-y-2">
                {[
                  ["MacBook Pro", "Singapore · Current session", "Now"],
                  ["Chrome on Windows", "Singapore", "2 days ago"],
                ].map(([device, loc, when]) => (
                  <div key={device as string} className="flex items-center gap-3 rounded-xl border border-line-warm px-4 py-3">
                    <div className="flex-1">
                      <p className="text-[12px] font-medium text-ink">{device}</p>
                      <p className="text-[11px] text-muted-warm">{loc}</p>
                    </div>
                    <p className="text-[11px] text-muted-warm">{when}</p>
                  </div>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </AppShell>
  );
}
