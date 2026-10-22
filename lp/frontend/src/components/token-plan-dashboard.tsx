"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { AppShell } from "@/components/app-shell";
import {
  ApiKeyCard,
  PlanBenefitsCard,
  PlanMetaStrip,
  PlanUsageCard,
  TokenPricingCard,
  TopSearchBar,
} from "@/components/token-plan-ui";
import { ApiError } from "@/lib/api";
import {
  createKey,
  deleteKey,
  fetchEndpoints,
  fetchKeys,
  fetchSubscription,
  formatCredits,
  type ApiKey,
  type Endpoint,
  type Site,
  type SubscriptionUsage,
} from "@/lib/api-keys";

// A user who registered before provisioning existed has no external account yet.
// The backend enqueues the missing rows on this first request and the worker
// provisions within a few seconds, so we retry briefly instead of erroring.
const PROVISIONING_RETRY_DELAY_MS = 3000;
const PROVISIONING_MAX_ATTEMPTS = 6;

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export function TokenPlanDashboard({ site }: { site: Site }) {
  const router = useRouter();
  const [subscription, setSubscription] = useState<SubscriptionUsage | null>(null);
  const [endpoints, setEndpoints] = useState<Endpoint[]>([]);
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [createdKey, setCreatedKey] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [keyError, setKeyError] = useState<string | null>(null);

  const handleError = useCallback(
    (err: unknown) => {
      if (err instanceof ApiError && err.status === 401) {
        router.replace("/login");
        return;
      }
      setKeyError(err instanceof Error ? err.message : "Request failed");
    },
    [router],
  );

  const load = useCallback(
    async (target: Site) => {
      // State is only touched after the first await, so the effect body never
      // performs a synchronous update.
      for (let attempt = 0; ; attempt++) {
        // Endpoints come from the platform's public status, so they load
        // independently of the authenticated calls.
        const [subscriptionResult, keysResult, endpointsResult] = await Promise.allSettled([
          fetchSubscription(target),
          fetchKeys(target),
          fetchEndpoints(target),
        ]);

        if (endpointsResult.status === "fulfilled") {
          setEndpoints(endpointsResult.value.items ?? []);
        } else {
          setEndpoints([]);
        }

        if (subscriptionResult.status === "fulfilled") {
          setSubscription(subscriptionResult.value);
        } else {
          setSubscription(null);
        }
        if (keysResult.status === "fulfilled") {
          setKeys(keysResult.value.items);
        } else {
          setKeys([]);
        }

        const failure =
          subscriptionResult.status === "rejected"
            ? subscriptionResult.reason
            : keysResult.status === "rejected"
              ? keysResult.reason
              : null;

        if (!failure) {
          setKeyError(null);
          setLoading(false);
          return;
        }

        // 409 means the platform account is still being created for this user.
        if (
          failure instanceof ApiError &&
          failure.status === 409 &&
          attempt < PROVISIONING_MAX_ATTEMPTS
        ) {
          setKeyError("Setting up your account on this platform…");
          await sleep(PROVISIONING_RETRY_DELAY_MS);
          continue;
        }

        handleError(failure);
        setLoading(false);
        return;
      }
    },
    [handleError],
  );

  useEffect(() => {
    void load(site);
  }, [site, load]);

  const handleCreate = async (name: string) => {
    try {
      // Unlimited quota and no expiry, matching the platform's default form.
      const created = await createKey(site, {
        name,
        unlimited_quota: true,
        expired_time: -1,
      });
      setCreatedKey(created.key);
      const page = await fetchKeys(site);
      setKeys(page.items);
    } catch (err) {
      handleError(err);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteKey(site, id);
      setKeys((prev) => prev.filter((key) => key.id !== id));
    } catch (err) {
      handleError(err);
    }
  };

  const plans = subscription?.plans ?? [];
  const creditsLabel =
    plans.length === 0
      ? "—"
      : plans.every((plan) => plan.unlimited)
        ? "Unlimited credits"
        : `${formatCredits(
            plans.reduce((sum, plan) => sum + (plan.credits_total ?? 0), 0),
          )} credits`;

  return (
    <AppShell
      title="Token plan"
      subtitle="Manage your subscription, credentials, and detailed model usage."
      actions={<TopSearchBar />}
    >
      <div className="mx-auto max-w-[1204px] space-y-6">
        <PlanMetaStrip usage={subscription} />
        <PlanUsageCard usage={subscription} loading={loading} />

        <div className="grid gap-6 lg:grid-cols-2">
          <ApiKeyCard
            keys={keys}
            loading={loading}
            error={keyError}
            createdKey={createdKey}
            onCreate={handleCreate}
            onDelete={handleDelete}
            onDismissCreated={() => setCreatedKey(null)}
            endpoints={endpoints}
          />
          <PlanBenefitsCard credits={creditsLabel} />
        </div>

        <TokenPricingCard />

        {/* Usage details is hidden for now: the chart, filters and tabs are demo data.
            Re-enable once a usage-analytics endpoint exists (see the design spec). */}
      </div>
    </AppShell>
  );
}