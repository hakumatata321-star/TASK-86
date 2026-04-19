# District Materials Commerce & Logistics Portal

**Project type: fullstack** — Go backend with server-rendered HTML UI (HTMX + Alpine.js + Bootstrap).

A web-based portal for managing district-wide distribution of educational materials.
Role-aware workflows cover students (browsing, ordering, favourites), instructors (course plans,
approvals), clerks (distribution, ledger), moderators (comment queue), and administrators
(users, analytics, settings).  Built with Go 1.22, Fiber, SQLite, HTMX, and Alpine.js.

---

## Quick Start (Docker — primary workflow)

Docker is the **only supported startup method**.  No local Go, Node, or GCC installation
is required.

```bash
# 1. Generate a .env file with cryptographically-random secrets
cat > .env << EOF
PORT=3000
DB_PATH=/app/data/portal.db
ENCRYPTION_KEY=$(openssl rand -hex 32)
SESSION_SECRET=$(openssl rand -hex 32)
APP_ENV=production
TIMEZONE=UTC
EOF

# 2. Start the portal
docker-compose up
```

> Run detached with `docker-compose up -d`, then tail logs with `docker-compose logs -f portal`.

The server listens on **http://localhost:3000**.

### First-boot admin password

On the very first start the admin account bootstrap is auto-rotated.
Retrieve the one-time password from the container log:

```bash
docker-compose logs portal | grep SECURITY
```

The `temporary_password` field in that line is the admin password.  Log in as `admin` and
change it immediately — the account is flagged `must_change_password = 1`.

### Other useful Docker commands

```bash
docker-compose up -d              # start detached
docker-compose logs -f portal     # tail logs
docker-compose down               # stop and remove containers
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up   # live template reload
```

---

## Demo Credentials

### Seeding demo accounts

Run the seed script **after** `docker-compose up` and after you have retrieved the admin
password (see above).

```bash
ADMIN_PASSWORD=<your-admin-password> bash scripts/seed_demo.sh
```

This script creates one account per role via the admin API.

### All-roles credentials table

| Username          | Email                          | Password                               | Role        |
|-------------------|--------------------------------|----------------------------------------|-------------|
| `admin`           | `admin@portal.local`           | `<SECURITY.temporary_password>`        | admin       |
| `demo_student`    | `demo_student@portal.local`    | `Demo1234!` | student     |
| `demo_instructor` | `demo_instructor@portal.local` | `Demo1234!` | instructor  |
| `demo_clerk`      | `demo_clerk@portal.local`      | `Demo1234!` | clerk       |
| `demo_moderator`  | `demo_moderator@portal.local`  | `Demo1234!` | moderator   |
| `demo_manager`    | `demo_manager@portal.local`    | `Demo1234!` | manager     |

> `SECURITY.temporary_password` means the exact `temporary_password` value from:
> `docker-compose logs portal | grep SECURITY`.
>
> All demo accounts (non-admin) are created by `scripts/seed_demo.sh` with the fixed
> password `Demo1234!`. They persist across server restarts.

---

## Verification

### Health & auth (curl)

```bash
# 1. Health check — should return {"status":"ok"}
curl http://localhost:3000/health

# 2. Login (replace credentials as needed)
curl -c cookies.txt -X POST http://localhost:3000/login \
  -d "username=demo_student&password=Demo1234!" \
  -L

# 3. List materials (authenticated)
curl -b cookies.txt http://localhost:3000/materials

# 4. View favorites lists (authenticated)
curl -b cookies.txt http://localhost:3000/favorites

# 5. Admin user list (requires admin session)
curl -c admin_cookies.txt -X POST http://localhost:3000/login \
  -d "username=admin&password=<admin-pass>" -L
curl -b admin_cookies.txt http://localhost:3000/admin/users

# 6. Distribution export CSV (admin only)
curl -b admin_cookies.txt \
  "http://localhost:3000/analytics/export/distribution" -o distribution.csv
```

### UI checklist

| Step | Action | Expected result |
|------|--------|-----------------|
| 1 | Open `http://localhost:3000/login` | Login form renders |
| 2 | Log in as `demo_student` / `Demo1234!` | Redirects to `/materials` |
| 3 | Click any material | Detail page shows rating + comments |
| 4 | Add to cart, place order | Order confirmation with ID |
| 5 | Open `http://localhost:3000/favorites` | Favorites lists page loads |
| 6 | Log in as `demo_instructor` / `Demo1234!` | Redirects to instructor dashboard |
| 7 | Open `/courses` | Course list renders |
| 8 | Log in as `demo_clerk` / `Demo1234!` | Redirects to `/distribution` |
| 9 | Log in as `admin` / `<admin-pass>` | Redirects to admin dashboard |
| 10 | Open `/analytics/map` (admin) | Leaflet map page with offline tiles |

---

## Running Tests

All tests (Go + frontend) run inside Docker — no local toolchain required.

```bash
./run_tests.sh            # run all suites
./run_tests.sh -v         # verbose output
./run_tests.sh -v -race   # verbose + Go race detector
```

`run_tests.sh` automatically launches a `golang:1.22-bookworm` container (with Node.js
installed for the Vitest frontend suite) and mounts the repo.  Exit code 0 = all pass.

**Test suites run:**

| Suite | Package(s) |
|-------|------------|
| Unit tests | `./unit_tests/...` |
| API functional tests | `./API_tests/...` |
| Integration tests (9 modules) | `./internal/integration/` |
| Config & scheduler unit tests | `./internal/config/...` `./internal/scheduler/...` |
| Middleware unit tests | `./internal/middleware/...` |
| Service unit tests | `./internal/services/...` |
| Frontend (Vitest — app.js, map.js) | `web/` via `npm test` |

---

## Environment Variables

All configuration is loaded from environment variables (or from `.env` via
`github.com/joho/godotenv`).

| Variable | Required | Default | Description |
|---|---|---|---|
| `ENCRYPTION_KEY` | **yes** | — | 64-char hex string (32 bytes) for AES-256-GCM field encryption. Generate: `openssl rand -hex 32`. |
| `SESSION_SECRET` | **yes** | — | Arbitrary secret string for session token signing. Generate: `openssl rand -hex 32`. |
| `PORT` | no | `3000` | TCP port the HTTP server listens on. |
| `DB_PATH` | no | `data/portal.db` | SQLite database file path. In Docker use `/app/data/portal.db`. |
| `APP_ENV` | no | `development` | Set to `production` to disable template hot-reload and enable stricter security. |
| `BANNED_WORDS` | no | *(empty)* | Comma-separated list of words blocked in material comments. |
| `TIMEZONE` | no | `UTC` | IANA timezone for Do-Not-Disturb window evaluation (e.g. `America/New_York`). |

`.env.example` is tracked with placeholder values only — never commit real secrets.

---

## Available Roles

| Role | Capabilities |
|---|---|
| `student` | Browse materials, place orders, manage favourites, inbox |
| `instructor` | Course plans, approve orders, approve/reject return & refund requests, inbox |
| `manager` | Approve/reject return & refund requests (same privileges as `instructor`) |
| `clerk` | Distribution events, ledger, backorder management, inbox |
| `moderator` | Review and act on reported comments, inbox |
| `admin` | Full access: user management, analytics, all settings, all of the above |

> Routes `GET /admin/returns`, `POST /admin/returns/:id/approve`, and
> `POST /admin/returns/:id/reject` accept `instructor`, `manager`, and `admin`.

---

## Project Structure

```
w2t86/
├── cmd/
│   └── server/
│       └── main.go              # Entry point: wires repos, services, handlers, routes
│   └── seed/
│       └── main.go              # Demo-account seeder (run once after first boot)
├── internal/
│   ├── config/                  # Environment-based configuration + validation
│   ├── crypto/                  # Password hashing + AES-256-GCM helpers
│   ├── db/                      # SQLite open + migration runner
│   ├── handlers/                # HTTP handlers (one file per domain)
│   │   ├── admin.go             # User management, custom fields, audit log
│   │   ├── analytics.go         # Dashboard stats, exports, geospatial map
│   │   ├── auth.go              # Login / logout
│   │   ├── courses.go           # Course plans (instructor)
│   │   ├── distribution.go      # Issue, return, exchange, reissue, ledger
│   │   ├── materials.go         # Browse, detail, rating, comments, favourites, share
│   │   ├── messages.go          # Inbox, SSE, DND, subscriptions
│   │   ├── moderation.go        # Reported comment review queue
│   │   └── orders.go            # Place, pay, cancel, returns, admin views
│   ├── middleware/
│   │   ├── auth.go              # Session validation, GetUser helper
│   │   ├── ratelimit.go         # Sliding-window rate limiter (comments)
│   │   └── rbac.go              # Role-based access control (RequireRole)
│   ├── models/                  # Go structs for every DB table
│   ├── observability/           # Structured loggers, request logger, metrics
│   ├── repository/              # Data access layer (one file per domain)
│   ├── scheduler/
│   │   └── scheduler.go         # Cron: auto-close stale orders every minute
│   ├── services/                # Business logic (one file per domain)
│   └── testutil/                # Shared in-memory DB helper for tests
├── migrations/                  # Numbered SQL migration files (001–016)
├── scripts/
│   └── seed_demo.sh             # Creates demo accounts via admin API
├── API_tests/                   # Black-box HTTP API tests (Fiber test runner)
├── internal/integration/        # Full end-to-end integration tests
├── unit_tests/                  # Pure unit tests (state machine, validation, etc.)
├── web/
│   ├── package.json             # Frontend test setup (Vitest)
│   ├── vitest.config.js
│   ├── tests/                   # Frontend unit tests (app.js, map.js)
│   ├── static/
│   │   ├── css/                 # Bootstrap, Bootstrap Icons, Leaflet, app.css
│   │   └── js/                  # HTMX, Alpine.js, Bootstrap bundle, Leaflet, app.js, map.js
│   └── templates/               # Go html/template files
├── Dockerfile
├── docker-compose.yml
├── docker-compose.dev.yml       # Development overrides (live template reload)
├── run_tests.sh                 # Docker-based full test suite runner
├── go.mod / go.sum
├── .env                         # Local secrets — gitignored, never committed
└── .env.example                 # Placeholder template — safe to commit
```

---

Frontend assets (HTMX, Alpine.js, Bootstrap, Leaflet) are vendored in `web/static/` — no npm
build step is needed for the server. The frontend Vitest suite is executed via `./run_tests.sh`.
