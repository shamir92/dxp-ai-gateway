# Pencil Design to Frontend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring all 18 screens in `docs/dxp-ai-gateway.pen` into the existing frontend while retaining current routes and application behavior.

**Architecture:** Keep the existing Next.js App Router and update its shared public and signed-in shells, design tokens, and route components to match the corresponding Pencil frames. Preserve existing API calls and client-side flows; do not add backend contracts or replace pages with static exports.

**Tech Stack:** Next.js 16.3.8, React 19.2.8, TypeScript, Tailwind CSS 4, existing shadcn/Base UI components, and lucide-react.

**Spec:** `docs/superpowers/specs/2026-10-12-pen-design-to-frontend.md`

## Global Constraints

- Keep existing Next.js routes, React components, and Tailwind v4 styling.
- Keep existing API calls, session handling, forms, loading/error states, and route destinations.
- Do not introduce backend contracts or fabricate data behavior for static screens.
- Adapt desktop canvas layouts to existing responsive conventions.
- Build on existing worktree edits in `docs/dxp-ai-gateway.pen`, `lp/frontend/src/app/profile/page.tsx`, `lp/frontend/src/app/settings/page.tsx`, and `lp/frontend/src/components/app-shell.tsx`; do not reset them.
- No backend endpoints, database schema changes, generated static HTML replacement, or new routes.

---

## File and Route Map

| Pencil frame | Route file |
| --- | --- |
| Landing page (`aHOSM`) | `lp/frontend/src/app/page.tsx` |
| Pricing — Token plans (`n5qFvS`) | `lp/frontend/src/app/pricing/page.tsx` |
| Models (`VU8Mz`) | `lp/frontend/src/app/models/page.tsx` |
| Contact (`F6FNk`) | `lp/frontend/src/app/contact/page.tsx` |
| API pricing (`WLdfJ`) | `lp/frontend/src/app/api-pricing/page.tsx` |
| Privacy Policy (`rpYof`) | `lp/frontend/src/app/privacy/page.tsx` |
| Terms of Use (`ScGz3`) | `lp/frontend/src/app/terms/page.tsx` |
| Login (`q0fyYD`) | `lp/frontend/src/app/login/page.tsx` |
| Register (`DDGrl`) | `lp/frontend/src/app/register/page.tsx` |
| Logout (`TBgep`) | `lp/frontend/src/app/logout/page.tsx` |
| Token plan — Single model (`GGJNQ`) | `lp/frontend/src/app/token-plan/page.tsx` |
| Token plan — Total consumption (`ZcLNt`) | `lp/frontend/src/app/consumption/page.tsx` |
| Profile & account (`Ww8Ku`) | `lp/frontend/src/app/profile/page.tsx` |
| API lifecycle (`t5dd6x`) | `lp/frontend/src/app/api-lifecycle/page.tsx` |
| Plugins (`WHb7v`) | `lp/frontend/src/app/plugins/page.tsx` |
| Settings (`ud6Xm`) | `lp/frontend/src/app/profile/page.tsx` (folded into Profile); `/settings` redirects to `/profile` |
| Documentation (`rTGze`) | `lp/frontend/src/app/documentation/page.tsx` |
| Support center (`mXEKz`) | `lp/frontend/src/app/support/page.tsx` |

The `Group · Non-login (public)` and `Group · After-login (app)` frames are navigation references, not separate pages.

## Tasks

### Task 1: Shared design tokens and shells

**Files:**
- Modify: `lp/frontend/src/app/layout.tsx`
- Modify: `lp/frontend/src/app/globals.css`
- Modify: `lp/frontend/src/components/public-shell.tsx`
- Modify: `lp/frontend/src/components/app-shell.tsx`
- Modify: `lp/frontend/src/components/brand.tsx`
- Reference: `docs/dxp-ai-gateway.pen`

**Interfaces:**
- Consumes: existing `PublicShell`, `AppShell`, and `Brand` call sites.
- Produces: the same shell interfaces with shared typography, colors, navigation, footer, and responsive layout aligned to Pencil.

- [ ] Read the installed Next.js 16 layout/font guidance at `lp/frontend/node_modules/next/dist/docs/` before modifying layout conventions.
- [ ] Inspect the `Landing page`, `Token plan — Single model`, and `Login` frames to confirm type, spacing, color, and shell boundaries.
- [ ] Align the global tokens with the Pencil variables. Preserve current token names where possible; the existing palette includes `ink: #151515`, `flame: #FF5A2A`, and `cream: #F7F4EE`.
- [ ] Update public header/footer and signed-in sidebar/header to match their shared frame structure without changing shell props or route destinations.
- [ ] Preserve current logo assets and ensure shell layouts remain usable below desktop widths.

### Task 2: Landing page

**Files:**
- Modify: `lp/frontend/src/app/page.tsx`
- Modify if needed: `lp/frontend/src/components/api-playground.tsx`

**Interfaces:**
- Consumes: `PublicShell`, `ApiPlayground`, current links and model/capability data.
- Produces: the designed landing route with the existing playground behavior intact.

- [ ] Read frame `aHOSM` section-by-section with `Get` and compare its screenshot to the current route.
- [ ] Match the announcement, split hero, metrics, model preview, developer example, capabilities, compatibility, CTA, and footer order shown in the frame.
- [ ] Keep `ApiPlayground` interactive and preserve all existing internal controls.
- [ ] Keep CTA destinations at `/register`, `/models`, `/login`, and `/pricing`.

### Task 3: Token-plan pricing page

**Files:**
- Modify: `lp/frontend/src/app/pricing/page.tsx`
- Modify if needed: `lp/frontend/src/components/pricing-shared.tsx`

**Interfaces:**
- Consumes: existing `PlanGrid`, `ModelsSection`, plan data, and `PublicShell`.
- Produces: `/pricing` matching frame `n5qFvS` while retaining plan links and FAQ interactions.

- [ ] Read frame `n5qFvS` and map its hero, plan cards, included models, comparison, setup guide, benefits, FAQ, and CTA to the current page sections.
- [ ] Update shared pricing components only for patterns used by this route or the API-pricing route.
- [ ] Preserve all subscription destinations and accordion behavior.

### Task 4: Models page

**Files:**
- Modify: `lp/frontend/src/app/models/page.tsx`
- Reuse: `lp/frontend/src/components/pricing-shared.tsx`

**Interfaces:**
- Consumes: existing model data and shared public shell.
- Produces: `/models` matching frame `VU8Mz` with existing model names and links preserved.

- [ ] Read frame `VU8Mz`, including the model catalog and detail hierarchy.
- [ ] Match its title area, catalog grouping, model metadata, and CTA/footer spacing.
- [ ] Preserve the current model copy and navigation targets unless the frame supplies replacement labels.

### Task 5: API-pricing page

**Files:**
- Modify: `lp/frontend/src/app/api-pricing/page.tsx`
- Reuse: `lp/frontend/src/components/pricing-shared.tsx`

**Interfaces:**
- Consumes: current pricing data and shared public shell.
- Produces: `/api-pricing` matching frame `WLdfJ`.

- [ ] Read frame `WLdfJ` and map its pricing overview, rate tables, explanatory sections, FAQ, and CTA.
- [ ] Keep displayed price/model data sourced from the current page constants; do not invent backend pricing calls.
- [ ] Reuse shared pricing styles only where the visual pattern is identical.

### Task 6: Contact and legal pages

**Files:**
- Modify: `lp/frontend/src/app/contact/page.tsx`
- Modify: `lp/frontend/src/app/privacy/page.tsx`
- Modify: `lp/frontend/src/app/terms/page.tsx`

**Interfaces:**
- Consumes: `PublicShell` and current legal/contact content.
- Produces: `/contact`, `/privacy`, and `/terms` matching frames `F6FNk`, `rpYof`, and `ScGz3`.

- [ ] Read the three target frames before editing; use `F6FNk` for the contact form layout and the legal frames for heading hierarchy and long-form text width.
- [ ] Preserve contact submission behavior and all existing legal copy; adjust structure and styles to match the design.
- [ ] Keep all cross-links within the current public routes.

### Task 7: Authentication and logout screens

**Files:**
- Modify: `lp/frontend/src/app/login/page.tsx`
- Modify: `lp/frontend/src/app/register/page.tsx`
- Modify: `lp/frontend/src/app/logout/page.tsx`
- Modify if needed: `lp/frontend/src/components/auth-form.tsx`

**Interfaces:**
- Consumes: existing auth form state machine, auth API, and route transitions.
- Produces: login, registration, and logout screens matching `q0fyYD`, `DDGrl`, and `TBgep`.

- [ ] Read each auth frame and compare it to the shared `AuthForm` modes.
- [ ] Match the compact centered form, logo treatment, fields, separators, and button states.
- [ ] Preserve email-code login/registration, busy/error/info states, and logout redirect behavior.

### Task 8: Single-model token-plan page

**Files:**
- Modify: `lp/frontend/src/app/token-plan/page.tsx`
- Modify if needed: `lp/frontend/src/components/token-plan-ui.tsx`

**Interfaces:**
- Consumes: existing view state, shared token-plan components, and `AppShell`.
- Produces: `/token-plan` matching frame `GGJNQ`.

- [ ] Read frame `GGJNQ` and map the sidebar/header, summary tiles, selected model content, usage values, and actions.
- [ ] Preserve existing view selection and route links.
- [ ] Keep desktop columns responsive without changing data semantics.

### Task 9: Total-consumption page

**Files:**
- Modify: `lp/frontend/src/app/consumption/page.tsx`
- Reuse: `lp/frontend/src/components/token-plan-ui.tsx`
- Reuse if present: `lp/frontend/src/components/ui/chart.tsx`

**Interfaces:**
- Consumes: existing chart components and app shell.
- Produces: `/consumption` matching frame `ZcLNt`.

- [ ] Read frame `ZcLNt` and identify metric cards, date controls, charts, tables, and filter states.
- [ ] Match chart labels and visual hierarchy using installed Recharts/shadcn chart APIs already in the project.
- [ ] Preserve any existing local/static data shape; do not add API calls.

### Task 10: API lifecycle page

**Files:**
- Modify: `lp/frontend/src/app/api-lifecycle/page.tsx`
- Reuse: `lp/frontend/src/components/token-plan-ui.tsx`
- Reuse: `lp/frontend/src/components/ui/table.tsx`

**Interfaces:**
- Consumes: current lifecycle records and app shell.
- Produces: `/api-lifecycle` matching frame `t5dd6x`.

- [ ] Read frame `t5dd6x` and map lifecycle summaries, timeline/table, status badges, and action controls.
- [ ] Preserve current record labels and action affordances; do not imply new server behavior.
- [ ] Keep wide tables usable on narrow viewports.

### Task 11: Plugins page

**Files:**
- Modify: `lp/frontend/src/app/plugins/page.tsx`

**Interfaces:**
- Consumes: current plugin catalog and app shell.
- Produces: `/plugins` matching frame `WHb7v`.

- [ ] Read frame `WHb7v` and map plugin list, status, search/tabs, detail area, permissions, and usage panels.
- [ ] Match control styling while preserving existing status values and nonfunctional affordances; do not create plugin endpoints.

### Task 12: Documentation and support pages

**Files:**
- Modify: `lp/frontend/src/app/documentation/page.tsx`
- Modify: `lp/frontend/src/app/support/page.tsx`

**Interfaces:**
- Consumes: current documentation/support content and app shell.
- Produces: `/documentation` and `/support` matching `rTGze` and `mXEKz`.

- [ ] Read both frames and map search, shortcuts, guides, updates, support channels, topics, tickets, and service status.
- [ ] Preserve existing labels, links, and local data; do not add search or ticket APIs.
- [ ] Keep cards and lists responsive at the existing breakpoint conventions.

### Task 13: Profile and settings pages

**Files:**
- Modify: `lp/frontend/src/app/profile/page.tsx`
- Modify: `lp/frontend/src/app/settings/page.tsx`
- Modify if needed: `lp/frontend/src/components/user-avatar.tsx`

**Interfaces:**
- Consumes: existing profile/session/settings APIs, current modified notification state, and app shell.
- Produces: `/profile` and `/settings` matching frames `Ww8Ku` and `ud6Xm`.

- [ ] Re-read the current worktree versions before editing and retain all uncommitted notification and navigation changes.
- [ ] Match the profile frame's account, connections, security, sessions, and notification sections while keeping existing API request payloads and error handling.
- [ ] Fold Settings frame notification controls into Profile and preserve `/settings` redirect behavior; do not add a Settings sidebar item.
- [ ] Preserve session revocation and connection-removal actions.

## Completion Review

- [ ] Confirm all 18 route/frame mappings above have corresponding page updates.
- [ ] Confirm no route, API contract, or auth flow was removed while applying visual changes.
- [ ] Review representative public, auth, and signed-in pages for layout consistency and responsiveness.
- [ ] Inspect the final diff and confirm unrelated and pre-existing worktree changes remain intact.
