# DXP AI Gateway

Control-plane API, web app, and provisioning worker for the DXP AI Gateway.

## Layout

| Path | What it is |
| --- | --- |
| `lp/backend` | Go / Gin control-plane API, listens on `:8080` |
| `lp/frontend` | Next.js 16 web app, listens on `:3000` |
| `lp/worker` | Provisioning worker (separate Go module) |
| `lp/vercel.json` | Production routing: both apps deployed as Vercel services |
| `Makefile` | Root task runner — `make help` lists everything |

## Local development

### Prerequisites

- Go (the toolchain is fetched automatically from `lp/backend/go.mod`)
- Node.js and npm
- Docker + Docker Compose, for the local Postgres

### 1. Configure the backend

```bash
cp lp/backend/.env.example lp/backend/.env
```

The server reads `lp/backend/.env` on start (real environment variables win over
the file). These values are **required** — the process exits on startup without
them:

- `DATABASE_URL` — the default matches the Compose Postgres, so no change needed
- `TPAI_BASE_URL` and `TPAI_ADMIN_TOKEN`
- `APAI_BASE_URL` and `APAI_ADMIN_TOKEN`
- `EXTERNAL_PASSWORD_ENCRYPTION_KEY` — a base64-encoded 32-byte key:

  ```bash
  openssl rand -base64 32
  ```

### 2. Start Postgres

```bash
make db-up      # postgres:16-alpine on localhost:5432
```

`make db-down` stops it, `make db-reset` wipes the volume, `make psql` opens a shell.

### 3. Start the backend

```bash
make run        # http://localhost:8080
```

Migrations run automatically on start; set `SKIP_MIGRATE=1` to skip them.
Confirm it is up with `make smoke`, which curls `/healthz`.

### 4. Start the frontend

```bash
make frontend-install   # once, per checkout
make frontend-dev       # http://localhost:3000
```

Run the backend and frontend in **two separate terminals** — `make dev-all`
only starts Postgres and the backend, then prints these steps.

### How the two are connected

The frontend never hardcodes a backend origin. Every call goes through
`apiFetch` (`lp/frontend/src/lib/api.ts`) as a **relative** path such as
`/api/v1/auth/me`. Something has to map that to the Go service:

| Environment | Mechanism |
| --- | --- |
| Production | `lp/vercel.json` rewrites `/api/(.*)` and `/healthz` to the `backend` service at Vercel's edge, before Next.js sees the request |
| Local | `lp/frontend/next.config.ts` rewrites the same paths to `BACKEND_URL`, defaulting to `http://localhost:8080` |

The local rewrite is gated on `NODE_ENV === "development"`, so a production
build emits no rewrites and Vercel's edge rules stay authoritative. The default
is enough for a normal checkout; override it if the backend is not on loopback:

```bash
BACKEND_URL=http://backend:8080 npm run dev
```

### Logging in locally

Set `MAIL_BACKEND=console` (the `.env.example` default) to use magic-code login
without configuring SMTP. The email body, including the code, is printed to the
**backend's** stdout — look for a `[mail:console]` line in the terminal running
`make run`.

Google and GitHub login additionally need `GOOGLE_CLIENT_ID` /
`GOOGLE_CLIENT_SECRET` and `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET`, with the
redirect URLs pointing at `http://localhost:8080/api/v1/auth/oauth/<provider>/callback`.
Leaving them blank does not stop the server from starting; those buttons simply
will not work.

### Verifying the wiring

With both servers running, a proxied path returns JSON from Go. A Next.js HTML
404 means the rewrite did not load:

```bash
curl -i http://localhost:8080/healthz       # direct to the backend
curl -i http://localhost:3000/api/v1/auth/me # through the frontend proxy
```

An unauthenticated `/api/v1/auth/me` should answer `401` with a JSON body —
not a `404`.

### Gotchas

- **Use `http://localhost`, not `127.0.0.1` or a LAN IP.** The session cookie is
  set with `Secure`, and browsers only accept a `Secure` cookie over plain HTTP
  for `localhost`. Other hostnames silently drop it, which shows up as a login
  loop with no error.
- **`localhost` is resolved by the Next.js server**, not the browser, because
  the rewrite is a server-side proxy. If the frontend runs in a container,
  `localhost` no longer means the backend — set `BACKEND_URL`.
- The frontend defaults to port 3000 and the backend to 8080. If either is
  taken, `next dev` picks another port and the backend honors `HTTP_ADDR`.
