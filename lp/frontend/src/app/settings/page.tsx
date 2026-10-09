import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { AppShell } from "@/components/app-shell";
import { TopSearchBar } from "@/components/token-plan-ui";

export default function SettingsPage() {
  return (
    <AppShell
      title="Settings"
      subtitle="Workspace preferences, notifications, and API defaults."
      actions={
        <div className="flex items-center gap-3">
          <TopSearchBar />
          <Button className="rounded-full bg-ink text-[12px] text-white hover:bg-ink/90">
            Save changes
          </Button>
        </div>
      }
    >
      <div className="mx-auto max-w-[760px] space-y-6">
        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">Workspace</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label>Default region</Label>
                <Select defaultValue="sgp">
                  <SelectTrigger className="rounded-xl">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="sgp">Singapore (sgp-1)</SelectItem>
                    <SelectItem value="us">United States (us-1)</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <Label>Default model</Label>
                <Select defaultValue="flash">
                  <SelectTrigger className="rounded-xl">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="mimo-v2.6-flash">mimo-v2.6-flash</SelectItem>
                    <SelectItem value="mimo-v2.6-pro">mimo-v2.6-pro</SelectItem>
                    <SelectItem value="deepseek-v4.1-flash">deepseek-v4.1-flash</SelectItem>
                    <SelectItem value="deepseek-v4-flash">deepseek-v4-flash</SelectItem>
                    <SelectItem value="glm-5.3-flash">glm-5.3-flash</SelectItem>
                    <SelectItem value="glm-5.3">glm-5.3</SelectItem>
                    <SelectItem value="kimi-k3">kimi-k3</SelectItem>
                    <SelectItem value="minimax-m3">minimax-m3</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">Notifications</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {[
              ["Balance alerts", "Email when balance falls below threshold", true],
              ["Usage digest", "Weekly usage and spend summary", true],
              ["Security alerts", "New device and credential changes", true],
              ["Product updates", "Model launches and API changes", false],
            ].map(([title, desc, on]) => (
              <div key={title as string} className="flex items-center justify-between gap-4">
                <div>
                  <p className="text-[13px] font-medium text-ink">{title}</p>
                  <p className="text-[11px] text-muted-warm">{desc}</p>
                </div>
                <Switch defaultChecked={on as boolean} />
              </div>
            ))}
          </CardContent>
        </Card>

        <Card className="border-line-warm bg-white">
          <CardHeader>
            <CardTitle className="text-[15px]">API defaults</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="text-[13px] font-medium text-ink">Stream responses by default</p>
                <p className="text-[11px] text-muted-warm">
                  Prefer streaming for chat completions when supported.
                </p>
              </div>
              <Switch defaultChecked />
            </div>
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="text-[13px] font-medium text-ink">Store request logs</p>
                <p className="text-[11px] text-muted-warm">
                  Keep metadata for usage analytics (payloads excluded).
                </p>
              </div>
              <Switch defaultChecked />
            </div>
          </CardContent>
        </Card>
      </div>
    </AppShell>
  );
}
