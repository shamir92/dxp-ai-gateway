import { apiJSON } from "@/lib/api";

export type Site = "tpai" | "apai";

/**
 * Public product slug per platform. The gateway never names the upstream
 * platform in a URL, so the token plan is served under `tp_*` endpoints; the
 * API plan's `ap_*` endpoints land with that page.
 */
const PRODUCT_SLUG: Record<Site, string> = { tpai: "tp", apai: "ap" };

const profilePath = (site: Site, suffix: string) =>
  `/api/v1/profile/${PRODUCT_SLUG[site]}_${suffix}`;

/** One API key as returned by the gateway. `key` is always masked here. */
export type ApiKey = {
  id: number;
  name: string;
  key: string;
  status: number;
  created_time: number;
  accessed_time: number;
  expired_time: number;
  remain_quota: number;
  unlimited_quota: boolean;
  model_limits_enabled: boolean;
  model_limits: string;
  allow_ips: string | null;
  used_quota: number;
  group: string;
  cross_group_retry: boolean;
  auto_groups: string[];
};

export type KeyPage = {
  items: ApiKey[];
  page: number;
  page_size: number;
  total: number;
};

/** One plan instance. A user may hold several active plans at once. */
export type SubscriptionView = {
  /** Identifies the instance; not unique per plan_id, since the same plan can be held twice. */
  subscription_id: number;
  plan_id: number;
  status: string;
  amount_total: number;
  amount_used: number;
  credits_total: number | null;
  credits_used: number;
  unlimited: boolean;
  start_time: number;
  end_time: number;
  next_reset_time: number;
};

/** The whole subscription state for one site. Plans keep the platform's priority order. */
export type SubscriptionUsage = {
  billing_preference: string;
  plans: SubscriptionView[];
};

export type CreatedKey = {
  token: ApiKey;
  /** Plaintext key as the platform generated it, returned exactly once. */
  key: string;
};

/**
 * Keys are displayed with an `sk-` prefix, mirroring the platform's own
 * dashboard: it stores the bare value and shows the prefix. Nothing is sent or
 * stored with the prefix — it is presentation only.
 */
export const KEY_DISPLAY_PREFIX = "sk-";

export function formatKeyForDisplay(key: string): string {
  return key ? `${KEY_DISPLAY_PREFIX}${key}` : "";
}

/** One advertised API base address. Works for OpenAI- and Anthropic-compatible clients. */
export type Endpoint = {
  label: string;
  url: string;
};

export type KeyInput = {
  name: string;
  unlimited_quota?: boolean;
  expired_time?: number;
  remain_quota?: number;
};

export function fetchSubscription(site: Site) {
  return apiJSON<SubscriptionUsage>(profilePath(site, "subscription"));
}

export function fetchEndpoints(site: Site) {
  return apiJSON<{ items: Endpoint[] }>(profilePath(site, "endpoints"));
}

export function fetchKeys(site: Site, page = 1, pageSize = 20) {
  const query = new URLSearchParams({
    p: String(page),
    page_size: String(pageSize),
  });
  return apiJSON<KeyPage>(`${profilePath(site, "api_keys")}?${query.toString()}`);
}

export function createKey(site: Site, input: KeyInput) {
  return apiJSON<CreatedKey>(profilePath(site, "api_keys"), {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateKey(site: Site, id: number, input: KeyInput) {
  return apiJSON<{ token: ApiKey }>(`${profilePath(site, "api_keys")}/${id}`, {
    method: "PUT",
    body: JSON.stringify(input),
  });
}

export function deleteKey(site: Site, id: number) {
  return apiJSON<{ ok: boolean }>(`${profilePath(site, "api_keys")}/${id}`, {
    method: "DELETE",
  });
}

export function formatCredits(value: number): string {
  if (!Number.isFinite(value)) return "—";
  if (Math.abs(value) >= 1e9) return `${(value / 1e9).toFixed(1)}B`;
  if (Math.abs(value) >= 1e6) return `${(value / 1e6).toFixed(1)}M`;
  if (Math.abs(value) >= 1e3) return `${(value / 1e3).toFixed(1)}K`;
  return String(value);
}

export function formatDateTime(unixSeconds: number): string {
  if (!unixSeconds) return "—";
  return new Date(unixSeconds * 1000).toISOString().slice(0, 10);
}