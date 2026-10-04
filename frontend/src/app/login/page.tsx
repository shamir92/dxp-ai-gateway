import { Brand } from "@/components/brand";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-cream px-6">
      <div className="w-full max-w-sm rounded-2xl border border-line-warm bg-white p-8">
        <Brand />
        <h1 className="mt-8 font-heading text-[24px] font-semibold text-ink">Welcome back</h1>
        <p className="mt-1 text-[13px] text-muted-warm">Sign in to your DXP AI Gateway workspace.</p>
        <div className="mt-6 space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" placeholder="you@company.com" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" placeholder="••••••••" />
          </div>
          <Button className="w-full rounded-full bg-ink text-white hover:bg-ink/90">Log in</Button>
        </div>
      </div>
    </div>
  );
}
