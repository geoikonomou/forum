# Forum

A full-stack community discussion platform built with **Go**. Users can register and log in (email/username/password or GitHub/Google OAuth), create topics, comment, upvote/downvote content, browse categories, and receive **real-time notifications** via Server-Sent Events.

The project follows **Clean Architecture** with **CQRS** and runs as two Go processes: a JSON **backend API** and a server-side-rendered **frontend client** that acts as a Backend-for-Frontend (BFF).

---

## Features

- **Authentication**
  - Email or username + password login with bcrypt hashing (cost 12)
  - Strong password validation (8+ chars, upper/lower/digit/special)
  - OAuth2 login via GitHub and Google
  - Session-based auth with rotating access/refresh tokens (HttpOnly, SameSite=Lax, Secure cookies)
- **Topics** — Create, read, update, delete with title, rich content, and optional image upload (JPEG/PNG/GIF, max 20 MB)
- **Categories** — Organize topics, assign multiple categories per topic, CRUD operations
- **Comments** — Nested discussion on topics with full CRUD
- **Voting** — Upvote/downvote (+1/-1) on both topics and comments, with toggle behavior and vote score calculation
- **Real-time Notifications** — Server-Sent Events (SSE) for replies, mentions, likes, and dislikes; view list, mark single/all as read
- **Activity Feed** — See your own created topics, voted items, and comments
- **Search & Filtering** — Search topics/categories, sort by date/title/vote score, filter by category
- **Pagination** — Configurable page size on topics and categories
- **Security** — Rate limiting (sliding window), input validation, TLS support, strict CORS
- **Infrastructure** — SQLite (WAL mode), Docker, GitHub Actions CI with linting, tests, and benchmark comparison

---

## Architecture

The application runs as **two separate Go servers**:

```
┌───────────────┐   HTTP/JSON    ┌──────────────────┐
│   Client      │ ─────────────► │   Backend API    │
│   (:3001)     │                │   (:8080)        │
│  SSR + BFF    │ ◄───────────── │  JSON REST + SSE │
└───────────────┘                └─────────┬────────┘
                                           │ SQLite
                                    ┌──────▼──────┐
                                    │   forum.db  │
                                    └─────────────┘
```

- **Backend API** (`cmd/server/`) — A pure JSON REST API. It owns the domain, all business logic, the database, and real-time notification streams.
- **Frontend Client** (`cmd/client/`) — Server-side rendered HTML using Go templates and vanilla JavaScript. It proxies requests to the backend (forwarding cookies, IP headers, and session data), acting as a **Backend-for-Frontend (BFF)**. This keeps the backend isolated and allows the two to be deployed/scaled independently.

### Clean Architecture & CQRS

The backend is organized into three layers, following the dependency rule (inner layers never depend on outer layers):

```
┌─────────────────────────────────────────────┐
│  Infrastructure (internal/infra)            │
│  HTTP handlers, middleware, SQLite repos,   │
│  session store, logger, notification pub/sub│
├─────────────────────────────────────────────┤
│  Application (internal/app)                 │
│  Use cases as Commands (writes) &           │
│  Queries (reads) — CQRS                     │
├─────────────────────────────────────────────┤
│  Domain (internal/domain)                   │
│  Pure entities + repository interfaces      │
└─────────────────────────────────────────────┘
```

1. **Domain layer** — Pure business entities (user, topic, comment, category, vote, session, oauth, activity, notification) and their repository interfaces. No external dependencies.
2. **Application layer** — Use cases implementing business workflows. Organized into **Commands** (write operations) and **Queries** (read operations). Depends only on the domain layer.
3. **Infrastructure layer** — HTTP handlers, middleware, SQLite repositories, session manager, logger, and the SSE notification service. Implements domain interfaces and depends on the application layer.

---

## Tech Stack

| Category | Technology |
|---|---|
| Language | **Go 1.24.4** |
| Frontend | Server-side rendered Go `html/template`, vanilla JavaScript, CSS (no framework) |
| Database | **SQLite3** via `mattn/go-sqlite3` (WAL journal mode) |
| Authentication | bcrypt (`golang.org/x/crypto`), session tokens, OAuth2 (GitHub + Google) |
| Realtime | Server-Sent Events (SSE) |
| UUIDs | `google/uuid` |
| Containerization | Docker (multi-stage), Docker Compose |
| CI/CD | GitHub Actions (lint, test, coverage, benchmark comparison) |
| Linting | golangci-lint v2, staticcheck |

---

## Project Structure

```
forum/
├── cmd/
│   ├── server/main.go            # Backend API server entry point
│   └── client/                   # Frontend BFF server
│       ├── main.go
│       ├── config/               # Client configuration
│       ├── server/               # HTTP handlers (pages + API proxies)
│       ├── middleware/           # Auth, required-auth, IP extraction
│       ├── domain/               # Client-side types (home, categories, activity...)
│       └── helpers/              # Template rendering, validation, response decode
├── internal/
│   ├── config/                   # Server config (env loading)
│   ├── app/                      # Application layer (use cases)
│   │   ├── user/                 #   Commands + Queries
│   │   ├── topics/
│   │   ├── comments/
│   │   ├── categories/
│   │   ├── votes/
│   │   ├── activities/
│   │   ├── oauth/
│   │   └── health/
│   ├── domain/                   # Domain entities + repository interfaces
│   │   ├── user/ topic/ comment/ category/ vote/
│   │   ├── session/ oauth/ activity/ notification/
│   ├── infra/                    # Infrastructure layer
│   │   ├── http/                 #   HTTP handlers by domain
│   │   ├── middleware/           #   Auth, CORS, rate limiting
│   │   ├── storage/sqlite/       #   SQLite repositories + migrations
│   │   ├── storage/sessionstore/ #   Session manager
│   │   ├── storage/notifications/#   Notification service + SSE
│   │   └── logger/               #   Leveled logger
│   └── pkg/                      # Shared packages
│       ├── bcrypt/               #   Password hashing
│       ├── uuid/
│       ├── validator/            #   Validation framework
│       ├── oAuth/                #   OAuth providers (GitHub, Google)
│       ├── helpers/              #   JSON responses, env parsing
│       └── testing/              #   Mocks + test helpers
├── frontend/
│   ├── html/
│   │   ├── layouts/base.html     # Base layout (navbar, footer, scripts)
│   │   ├── pages/                # Page templates
│   │   └── partials/             # Reusable fragments
│   └── static/
│       ├── css/                  # 9 stylesheets
│       ├── js/                   # script, topic, votes, notifications, create-post
│       └── images/
├── db/
│   ├── migrations/               # schema.sql, indexes.sql
│   └── seeds/dev_data.sql        # Development seed data
├── docs/
│   ├── backend_api_reference.md  # Full API documentation
│   └── diagrams/                 # Mermaid flow diagrams
├── certs/                        # TLS certificates (localhost)
├── .github/workflows/            # CI/CD pipelines
├── Dockerfile
├── docker-compose.yml            # Production orchestration
├── docker-compose.dev.yml        # Development overlay
├── entrypoint.sh                 # Starts both servers
├── Makefile
└── .env.example                  # Environment variable reference
```

---

## Getting Started

### Prerequisites

- **Go 1.24.4+**
- **gcc / CGO** enabled (required by the `go-sqlite3` driver)
- Optional: Docker + Docker Compose, SQLite CLI

### Local Setup

```bash
# 1. Configure environment
cp .env.example .env

# 2. Run the server (backend + frontend) — or use the Makefile
go run ./cmd/server            # backend on :8080
go run ./cmd/client            # frontend on :3001 (in a second terminal)
```

> On first start, migrations and seed data are applied automatically (controlled by `DB_MIGRATE_ON_START` / `DB_SEED_ON_START`).

### Access Points

| Service | URL |
|---|---|
| Frontend | http://localhost:3001 |
| Backend API | http://localhost:8080/api/v1 |
| Health Check | http://localhost:8080/api/v1/health |

### Docker

See the dedicated [DOCKER_SETUP.md](./DOCKER_SETUP.md) for full instructions.

```bash
make docker-up     # Build + start (production)
make docker-dev    # Start development environment with live mounts
```

---

## Configuration

All settings are read from environment variables, with a `.env` file and sensible defaults as fallback. See [.env.example](./.env.example) for the full list.

| Category | Variable | Default | Description |
|---|---|---|---|
| Server | `SERVER_HOST` | `localhost` | Backend bind host |
| | `SERVER_PORT` | `8080` | Backend port |
| | `SERVER_ENVIRONMENT` | `development` | Environment name |
| Client | `CLIENT_HOST` | `localhost` | Frontend bind host |
| | `CLIENT_PORT` | `3001` | Frontend port |
| | `BACKEND_URL` | `http://localhost:8080/api/v1` | Backend URL used by the client |
| Database | `DB_DRIVER` | `sqlite3` | Database driver |
| | `DB_PATH` | `db/data/forum.db` | Database file path |
| | `DB_MIGRATE_ON_START` | `true` | Run migrations on startup |
| | `DB_SEED_ON_START` | `true` | Load seed data on startup |
| | `DB_PRAGMA` | `_foreign_keys=on&_journal_mode=WAL` | SQLite pragmas |
| Session | `SESSION_DEFAULT_EXPIRY` | `1600` | Access token expiry (seconds) |
| | `SESSION_SECURE_COOKIE` | `false` | Set `Secure` flag on cookies (set `true` in prod) |
| | `SESSION_COOKIE_NAME` | `session_id` | Cookie name |
| | `SESSION_MAX_SESSIONS_PER_USER` | `5` | Max concurrent sessions per user |
| Rate Limit | `RATE_LIMIT_ENABLED` | `true` | Enable rate limiting |
| | `RATE_LIMIT_REQUESTS` | `100` | Requests per window |
| | `RATE_LIMIT_WINDOW_SECONDS` | `60` | Window size |
| OAuth | `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET` | — | GitHub OAuth credentials |
| | `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | — | Google OAuth credentials |
| TLS | `SERVER_TLS_CERT_FILE` / `SERVER_TLS_KEY_FILE` | — | Backend TLS certificates |
| | `CLIENT_TLS_CERT_FILE` / `CLIENT_TLS_KEY_FILE` | — | Frontend TLS certificates |

---

## API Reference

The backend exposes a versioned JSON REST API under `/api/v1`, documented in detail in [docs/backend_api_reference.md](./docs/backend_api_reference.md).

### Endpoint Groups

| Group | Endpoints |
|---|---|
| Health | `GET /health` |
| Auth | `POST /register`, `POST /login/email`, `POST /login/username`, `POST /logout`, `GET /me` |
| OAuth | `GET /auth/{github,google}/login`, `GET /auth/{github,google}/callback` |
| Topics | `POST /topics/create`, `PUT /topics/update`, `DELETE /topics/delete`, `GET /topic`, `GET /topics/all` |
| Comments | `POST /comments/create`, `PUT /comments/update`, `DELETE /comments/delete`, `GET /comments/get`, `GET /comments/topic` |
| Categories | `POST /category/create`, `PUT /category/update`, `DELETE /category/delete`, `GET /category`, `GET /categories/all` |
| Votes | `POST /vote/cast`, `DELETE /vote/delete`, `GET /vote/counts` |
| Activity | `GET /user/activity` |
| Notifications | `GET /notifications/stream` (SSE), `GET /notifications`, `GET /notifications/unread-count`, `POST /notifications/mark-read`, `POST /notifications/mark-all-read` |

Responses use a consistent wrapper: success returns `{"info": ..., "data": ...}`; errors return `{"error": "message"}`.

---

## Database

SQLite schema with 9 tables (see [db/migrations/schema.sql](./db/migrations/schema.sql)):

```
users 1───* oauth_providers
users 1───* sessions
users 1───* topics  *───* categories  (via topic_categories junction)
users 1───* comments *───1 topics
users 1───* votes    *───1 topics    (nullable)
users 1───* votes    *───1 comments  (nullable)
users 1───* notifications
```

- **users** — email/username/password (bcrypt hash, nullable for OAuth users)
- **oauth_providers** — linked provider accounts (GitHub/Google)
- **sessions** — access + refresh tokens with expiry
- **categories** — forum categories with name, description, image, color, slug
- **topics** — posts with title, content, optional image
- **topic_categories** — many-to-many junction
- **comments** — topic discussion
- **votes** — upvote/downvote with partial unique indexes to prevent double voting
- **notifications** — real-time notification records (reply/mention/like/dislike)

Migrations and indexes are applied on startup (when `DB_MIGRATE_ON_START=true`), wrapped in transactions.

---

## Testing

```bash
make test              # Run tests with race detection + coverage
make test-short        # Run quick tests only
make lint              # Run all linters (staticcheck + golangci-lint)
make format            # Format code with goimports + gofmt
make check-format      # Verify formatting (CI)
```

Tests use hand-written mocks (`internal/pkg/testing/mocks.go`) and table-driven style covering the application use cases, middleware, and utilities.

### Benchmarks

```bash
make ci-bench          # Run all benchmarks
make bench-compare     # Compare current branch vs main (benchstat)
make bench-profile     # CPU/memory profiling
make bench-flame       # Generate CPU flame graph
```

Benchmarks run automatically on PRs via GitHub Actions (`.github/workflows/bench_ci.yml`).

---

## CI/CD

GitHub Actions workflows in [.github/workflows](./.github/workflows):

- **ci.yml** — Runs on push to `main` and PRs: module verification, formatting check, staticcheck, golangci-lint, and tests with coverage (uploaded to Codecov).
- **bench_ci.yml** — Runs benchmarks on PRs and produces a benchstat comparison between the PR branch and `main`.

The complete pipeline is available via `make ci`.

---

## Docker Deployment

The application ships as a single container running both the backend and frontend, orchestrated via Docker Compose.

| Command | Description |
|---|---|
| `make docker-build` | Build the image |
| `make docker-up` | Start services (production) |
| `make docker-down` | Stop services |
| `make docker-logs` | View logs |
| `make docker-restart` | Restart |
| `make docker-dev` / `make docker-dev-build` | Development environment |
| `make docker-clean` | Remove everything |

Highlights:
- **Multi-stage** builds for a small, optimized image (non-root user `appuser`)
- **Persistent volumes** for the database (`forum-db`) and uploads (`forum-uploads`)
- **Health checks** on the backend
- Dev mode adds volume mounts for live code changes and seed data

For production, set `SESSION_SECURE_COOKIE=true` and consider a reverse proxy for HTTPS/TLS.

---

## OAuth Setup

To enable GitHub / Google login:

1. Register an application with the provider to obtain a Client ID and Client Secret.
2. For **GitHub**: create an OAuth App (scopes `user:email`) with the redirect URL set to the backend callback (e.g. `http://localhost:8080/api/v1/auth/github/callback`).
3. For **Google**: create an OAuth 2.0 Client ID with the authorized redirect URI set to the backend callback.
4. Set the corresponding environment variables (see [Configuration](#configuration)) in your `.env` or `docker-compose.yml`.

OAuth state tokens are managed in-memory with a 10-minute expiry for CSRF protection.

---

## Makefile Reference

| Command | Description |
|---|---|
| `make env` | Show environment info |
| `make tools` | Install CI tools (goimports, staticcheck, golangci-lint) |
| `make ci` | Full CI: mod verify, format check, lint, test |
| `make clean` | Remove coverage output |
| `make db-clean` | Remove the database directory |
| `make docker-*` | All Docker workflow commands |

---

## Contributing

1. Fork the repository and create a feature branch.
2. Write tests for new functionality (table-driven style, using the existing mocks).
3. Run `make format` and `make lint` to conform to the project's style (import ordering and all default linters are enforced by `.golangci.yml`).
4. Run `make test` to ensure all tests pass.
5. Open a pull request against `main`.

The project follows Go's standard file layout and enforces strict linting (all linters enabled by default, with noisy ones selectively disabled).

---

## License

This project is licensed under the **MIT License**. See the [LICENSE](./LICENSE) file for details.
