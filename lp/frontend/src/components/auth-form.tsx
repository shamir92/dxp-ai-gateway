"use client";

import { useState } from "react";
import Link from "next/link";
import { Brand } from "@/components/brand";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { apiFetch } from "@/lib/api";

export function AuthForm({ mode }: { mode: "login" | "register" }) {
  const [step, setStep] = useState<"email" | "code">("email");
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [errorAction, setErrorAction] = useState<"login" | "register" | null>(null);
  const [info, setInfo] = useState("");
  const [busy, setBusy] = useState(false);

  const purpose = mode === "register" ? "register" : "login";

  async function requestCode() {
    setBusy(true);
    setError("");
    setErrorAction(null);
    setInfo("");
    try {
      const res = await apiFetch("/api/v1/auth/magic/request", {
        method: "POST",
        body: JSON.stringify({
          email,
          purpose,
          ...(mode === "register" ? { name } : {}),
        }),
      });
      const body = await res.json().catch(() => null);
      if (!res.ok) {
        const errBody = body as { error?: { code?: string; message?: string } } | null;
        const code = errBody?.error?.code;
        if (code === "no_account") {
          setError("No account for this email.");
          setErrorAction("register");
          return;
        }
        if (code === "account_exists") {
          setError("Account already exists.");
          setErrorAction("login");
          return;
        }
        throw new Error(errBody?.error?.message ?? `Could not send code (${res.status})`);
      }
      setInfo(`Code sent to ${email}. Check your inbox.`);
      setStep("code");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Could not send code");
    } finally {
      setBusy(false);
    }
  }

  async function verifyCode() {
    setBusy(true);
    setError("");
    try {
      const res = await apiFetch("/api/v1/auth/magic/verify", {
        method: "POST",
        body: JSON.stringify({ email, code, purpose }),
      });
      if (!res.ok) {
        const body = await res.json().catch(() => null);
        throw new Error(body?.error?.message ?? "Invalid code");
      }
      window.location.href = "/token-plan";
    } catch (e) {
      setError(e instanceof Error ? e.message : "Invalid code");
    } finally {
      setBusy(false);
    }
  }

  async function resendCode() {
    setBusy(true);
    setError("");
    try {
      const res = await apiFetch("/api/v1/auth/magic/resend", {
        method: "POST",
        body: JSON.stringify({ email, purpose }),
      });
      if (!res.ok) {
        const body = await res.json().catch(() => null);
        throw new Error(body?.error?.message ?? "Could not resend");
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : "Could not resend");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="auth-shell flex min-h-screen items-center justify-center bg-cream px-6 py-12">
      <div className="w-full max-w-[432px] rounded-2xl border border-line-warm bg-white p-8">
        <Brand href="/" />
        <h1 className="mt-8 font-heading text-[24px] font-semibold text-ink">
          {mode === "login" ? "Welcome back" : "Create your account"}
        </h1>
        <p className="mt-1 text-[13px] text-muted-warm">
          {mode === "login"
            ? "Sign in to your DXP AI Gateway workspace. No password needed."
            : "Get your DXP AI Gateway workspace in minutes."}
        </p>

        <div className="mt-6">
          {step === "email" ? (
            <form
              className="space-y-4"
              onSubmit={(e) => {
                e.preventDefault();
                void requestCode();
              }}
            >
              {mode === "register" ? (
                <div className="space-y-1.5">
                  <Label htmlFor="name">Full name</Label>
                  <Input
                    id="name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Ada Lovelace"
                    autoComplete="name"
                    required
                  />
                </div>
              ) : null}
              <div className="space-y-1.5">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="you@company.com"
                  autoComplete="email"
                />
              </div>
              {mode === "register" ? (
                <p className="text-[11px] leading-relaxed text-muted-warm">
                  By continuing you agree to our{" "}
                  <Link href="/terms" className="text-flame">
                    Terms
                  </Link>{" "}
                  and{" "}
                  <Link href="/privacy" className="text-flame">
                    Privacy Policy
                  </Link>
                  .
                </p>
              ) : null}
              {error ? (
                <div className="rounded-xl border border-flame/30 bg-flame/5 px-3 py-2.5 text-[12px] text-flame">
                  <p>{error}</p>
                  {errorAction === "register" ? (
                    <p className="mt-1">
                      <Link href="/register" className="font-medium underline">
                        Create an account
                      </Link>
                    </p>
                  ) : null}
                  {errorAction === "login" ? (
                    <p className="mt-1">
                      <Link href="/login" className="font-medium underline">
                        Log in instead
                      </Link>
                    </p>
                  ) : null}
                </div>
              ) : null}
              <Button
                type="submit"
                disabled={busy}
                className="h-11 w-full rounded-full bg-ink text-[13px] text-white hover:bg-ink/90"
              >
                {busy
                  ? "Sending…"
                  : mode === "login"
                    ? "Email me a sign-in code"
                    : "Email me a sign-up code"}
              </Button>
            </form>
          ) : (
            <form
              className="space-y-4"
              onSubmit={(e) => {
                e.preventDefault();
                void verifyCode();
              }}
            >
              <p className="text-[12px] leading-relaxed text-muted-warm">
                We sent a 6-digit code to <span className="text-ink">{email || "your email"}</span>.
                Enter it below to {mode === "login" ? "sign in" : "finish creating your account"}.
              </p>
              <div className="space-y-1.5">
                <Label htmlFor="code">Magic code</Label>
                <Input
                  id="code"
                  inputMode="numeric"
                  pattern="[0-9]*"
                  maxLength={6}
                  required
                  value={code}
                  onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                  placeholder="123456"
                  className="text-center font-mono text-lg tracking-[0.4em]"
                  autoComplete="one-time-code"
                />
              </div>
              {info ? <p className="text-[12px] text-success">{info}</p> : null}
              {error ? <p className="text-[12px] text-flame">{error}</p> : null}
              <Button
                type="submit"
                disabled={busy}
                className="h-11 w-full rounded-full bg-ink text-[13px] text-white hover:bg-ink/90"
              >
                {busy ? "Verifying…" : "Verify & continue"}
              </Button>
              <div className="flex items-center justify-between text-[11px]">
                <button
                  type="button"
                  className="text-muted-warm hover:text-ink"
                  onClick={() => {
                    setStep("email");
                    setError("");
                  }}
                >
                  ← Use a different email
                </button>
                <button type="button" className="text-flame" onClick={() => void resendCode()} disabled={busy}>
                  Resend code
                </button>
              </div>
            </form>
          )}

          <p className="mt-6 text-center text-[12px] text-muted-warm">
            {mode === "login" ? (
              <>
                New to DXP?{" "}
                <Link href="/register" className="text-flame">
                  Create account
                </Link>
              </>
            ) : (
              <>
                Already have an account?{" "}
                <Link href="/login" className="text-flame">
                  Log in
                </Link>
              </>
            )}
          </p>
        </div>
      </div>
    </div>
  );
}
