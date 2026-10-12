# Pencil Design to Frontend

## Goal

Implement the screens in `docs/dxp-ai-gateway.pen` in the existing frontend, keeping the current Next.js routes and application behavior.

## Scope

Translate all 18 designed screens:

| Design | Route |
| --- | --- |
| Landing page | `/` |
| Pricing — Token plans | `/pricing` |
| Models | `/models` |
| Contact | `/contact` |
| API pricing | `/api-pricing` |
| Privacy Policy | `/privacy` |
| Terms of Use | `/terms` |
| Login | `/login` |
| Register | `/register` |
| Logout | `/logout` |
| Token plan — Single model | `/token-plan` |
| Token plan — Total consumption | `/consumption` |
| Profile & account | `/profile` |
| API lifecycle | `/api-lifecycle` |
| Plugins | `/plugins` |
| Settings | `/settings` |
| Documentation | `/documentation` |
| Support center | `/support` |

The two group frames in the design organize public and signed-in screens; they are not routes.

## Approach

- Keep the existing Next.js App Router, page routes, React components, and Tailwind v4 styling.
- Use `PublicShell` and `AppShell` as the shared implementation points for the public header/footer and signed-in navigation/header.
- Match the design's colors, type, spacing, and layout using the tokens and typography already present in the Pencil file and frontend.
- Update each route's page content to match its corresponding frame. Reuse existing UI components and extract focused shared components only where the design repeats a real pattern.
- Keep existing API calls, session handling, forms, loading/error states, and route destinations. Do not introduce backend contracts or fabricate data behavior for static screens.
- Adapt the desktop canvas layouts to existing responsive conventions.

## Implementation order

1. Review the relevant installed Next.js 16 documentation and the existing route/component behavior.
2. Align global design tokens and shared public and signed-in shells.
3. Implement public pages and auth screens.
4. Implement signed-in pages, preserving existing API interactions.
5. Review the routes for visual consistency, responsiveness, and preserved behavior.

## Existing worktree changes

At design time, the worktree already contained user changes in `docs/dxp-ai-gateway.pen`, `lp/frontend/src/app/profile/page.tsx`, `lp/frontend/src/app/settings/page.tsx`, and `lp/frontend/src/components/app-shell.tsx`. Implementation must build on those changes and must not reset them.

## Validation

Use the project's available lint/build checks and visually inspect representative public, auth, and signed-in routes in the browser. No backend or database changes are in scope.

## Out of scope

- Backend endpoints, database schema, and authentication contract changes.
- Replacing the App Router with generated static HTML.
- Adding routes not represented in the Pencil design.
