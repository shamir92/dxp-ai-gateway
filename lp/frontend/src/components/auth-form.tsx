"use client";

import { useState } from "react";
import Link from "next/link";
import { Brand } from "@/components/brand";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

function OAuthButtons({ mode }: { mode: "login" | "register" }) {
  const label = mode === "login" ? "Continue with Google" : "Sign up with Google";
  const gh = mode === "login" ? "Continue with GitHub" : "Sign up with GitHub";

  return (
    <div className="space-y-2.5">
      <Button
        type="button"
        variant="outline"
        className="h-11 w-full rounded-full border-line-warm bg-white text-[13px] text-ink"
        onClick={() => {
          window.location.href = "/api/auth/signin/google";
        }}
      >
        <svg className="mr-2 size-4" viewBox="0 0 24 24" aria-hidden>
          <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
          <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
          <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
          <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
        </svg>
        {label}
      </Button>
      <Button
        type="button"
        variant="outline"
        className="h-11 w-full rounded-full border-line-warm bg-white text-[13px] text-ink"
        onClick={() => {
          window.location.href = "/api/auth/signin/github";
        }}
      >
        <svg className="mr-2 size-4" viewBox="0 0 24 24" fill="currentColor" aria-hidden>
          <path d="M12 .3a12 12 0 0 0-3.79 23.4c.6.1.82-.26.82-.58v-2.02c-3.34.73-4.04-1.61-4.04-1.61-.55-1.39-1.34-1.76-1.34-1.76-1.09-.75.08-.73.08-.73 1.21.08 1.85 1.24 1.85 1.24 1.07 1.84 2.81 1.31 3.5 1 .1-.78.42-1.31.76-1.61-2.67-.3-5.47-1.34-5.47-5.96 0-1.32.47-2.39 1.24-3.23-.13-.3-.54-1.52.11-3.18 0 0 1.01-.32 3.3 1.23a11.4 11.4 0 0 1 6 0c2.29-1.55 3.3-1.23 3.3-1.23.65 1.66.24 2.88.12 3.18.77.84 1.23 1.91 1.23 3.23 0 4.63-2.8 5.65-5.48 5.95.43.37.81 1.1.81 2.22v3.29c0 .32.21.69.82.57A12 12 0 0 0 12 .3z" />
        </svg>
        {gh}
      </Button>
    </div>
  );
}

function Divider() {
  return (
    <div className="my-5 flex items-center gap-3">
      <div className="h-px flex-1 bg-line-warm" />
      <span className="text-[11px] text-muted-warm">or</span>
      <div className="h-px flex-1 bg-line-warm" />
    </div>
  );
}

export function AuthForm({ mode }: { mode: "login" | "register" }) {
  const [step, setStep] = useState<"email" | "code">("email");
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [code, setCode] = useState("");

  return (
    <div className="flex min-h-screen items-center justify-center bg-cream px-6 py-12">
      <div className="w-full max-w-sm rounded-2xl border border-line-warm bg-white p-8">
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
          <OAuthButtons mode={mode} />
          <Divider />

          {step === "email" ? (
            <form
              className="space-y-4"
              onSubmit={(e) => {
                e.preventDefault();
                setStep("code");
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
              <Button
                type="submit"
                className="h-11 w-full rounded-full bg-ink text-[13px] text-white hover:bg-ink/90"
              >
                {mode === "login" ? "Email me a sign-in code" : "Email me a sign-up code"}
              </Button>
            </form>
          ) : (
            <form
              className="space-y-4"
              onSubmit={(e) => {
                e.preventDefault();
                window.location.href = "/token-plan";
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
              <Button
                type="submit"
                className="h-11 w-full rounded-full bg-ink text-[13px] text-white hover:bg-ink/90"
              >
                Verify & continue
              </Button>
              <div className="flex items-center justify-between text-[11px]">
                <button
                  type="button"
                  className="text-muted-warm hover:text-ink"
                  onClick={() => setStep("email")}
                >
                  ← Use a different email
                </button>
                <button type="button" className="text-flame">
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
