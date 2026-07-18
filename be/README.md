# UMSRMS — Backend API

Go REST API for the **University Maintenance & Service Request Management System**
(MIT 8333). Handles authentication, role-based access, service-request lifecycle,
assignments, notifications, reporting and an audit trail, documented with Swagger.

## Tech stack

| Concern | Choice |
|---------|--------|
| Language | **Go 1.26** |
| HTTP framework | **Gin** |
| Database | **PostgreSQL** (via `database/sql` + `pgx/v5` stdlib driver) |
| Query builder | **Squirrel** |
| Auth | **JWT** (`golang-jwt/v5`) + **bcrypt** password hashing |
| Validation | **go-playground/validator** (via Gin binding tags) |
| API docs | **Swagger** (`swaggo/swag` + `gin-swagger`) |
| Migrations | **golang-migrate** |
| Live reload | **air** |

## Requirements

- Go 1.26+
- PostgreSQL 14+
- CLI tools (installed via Make): `swag`, `migrate` (with the `postgres` driver), `air`

```sh
make install-swagger
make install-migrate
```

## Getting started

```sh
# 1. Configure environment
cp .env.example .env          # then edit DB creds + JWT_SECRET

# 2. Create the database (once)
#    createdb umsrms          # or via your Postgres client

# 3. Run — migrations are applied automatically on startup
make dev                      # live reload (air)
# or
make run                      # run once
```

> **Migrations run automatically on boot** (embedded in the binary via `//go:embed`).
> Set `RUN_MIGRATIONS=false` to disable, or apply them manually with `make migrate-up`.
> The app only needs the database to *exist* — it creates the schema itself.

The API serves on `http://localhost:8080`, base path **`/api/v1`**.
Swagger UI: **http://localhost:8080/swagger/index.html**

## Environment variables

See [`.env.example`](.env.example). Key ones:

| Variable | Purpose |
|----------|---------|
| `PORT` | HTTP port (default 8080) |
| `JWT_SECRET` | Token signing secret (**required in production**) |
| `DB_*` | PostgreSQL connection (`DB_SSLMODE` = `disable` locally, `require` for most cloud DBs) |
| `CORS_ALLOWED_ORIGINS` | `*` for all, or a comma-separated allowlist |
| `RATE_LIMIT_*` | Global per-IP token bucket (auth routes are additionally capped at 5/min) |
| `UPLOAD_DIR` / `MAX_UPLOAD_MB` | Evidence upload location and size limit |

## Project structure

```
be/
├── main.go                 # entrypoint, DI wiring, router
├── internal/
│   ├── config/             # env + Postgres connection
│   ├── models/             # DB entity structs
│   ├── dto/                # request/response payloads (+ binding validation tags)
│   ├── repository/         # SQL data access (Squirrel)
│   ├── service/            # business logic
│   ├── handler/            # Gin HTTP handlers (+ Swagger annotations)
│   ├── middleware/         # CORS, rate limit, JWT auth, RBAC, audit
│   └── utils/              # JWT, bcrypt, responses, validation, token ban list
├── migrations/             # golang-migrate SQL files
├── docs/                   # generated Swagger (swag init)
├── uploads/                # runtime evidence files (gitignored)
└── Makefile
```

Layering: **handler → service → repository → PostgreSQL**, with cross-cutting
concerns in `middleware`.

## Roles & access control

Three roles, enforced by JWT auth + `RequireRoles` middleware (and re-checked in
services):

- **student_staff** — submit and track their own requests, upload evidence
- **maintenance_officer** — view assigned requests, update status
- **admin** — manage users, assign requests, view reports & audit log

## API overview

Base path `/api/v1`. Full, interactive reference at `/swagger/index.html`.

| Group | Endpoints |
|-------|-----------|
| **Auth** | `POST /auth/register`, `POST /auth/login`, `GET /auth/registration-data`, `POST /auth/logout` 🔒, `GET /auth/me` 🔒 |
| **Requests** | `POST /requests` (student), `GET /requests`, `GET /requests/:id`, `PUT /requests/:id/status` (officer/admin), `POST /requests/:id/assign` (admin), `DELETE /requests/:id` (admin), `POST /requests/:id/evidence` (owner) |
| **Users** (admin) | `GET /users`, `GET /users/officers`, `PUT /users/:id/role`, `DELETE /users/:id` |
| **Categories** | `GET /categories` |
| **Notifications** | `GET /notifications`, `PUT /notifications/:id/read`, `PUT /notifications/read-all` |
| **Reports** (admin) | `GET /reports/summary` |
| **Audit** (admin) | `GET /audit-logs` |
| **Health** | `GET /ping` |

All responses use a consistent envelope: `{ success, message, data | error }`.

## Advanced features

- **JWT auth** with a server-side token ban list (real logout / revocation)
- **RBAC** middleware
- **File/image upload** for fault evidence (type + size validated)
- **In-app notifications** emitted on request create / assign / status change
- **Search, filter, pagination** on the request list
- **Audit trail** — middleware records mutating actions with the acting user
- **Reporting** — totals and counts by status/category
- **CORS** + **rate limiting** (global per-IP; stricter 5/min on auth)
- **Swagger** API documentation

## Database

PostgreSQL with 8 core tables — `roles`, `users`, `request_categories`,
`service_requests`, `assignments`, `status_logs`, `notifications`, `audit_logs`.
Relationships: one-to-many (user → requests, request → status logs), many-to-one
(request → category), and many-to-many (officers ↔ requests via `assignments`).
Managed with golang-migrate; roles and categories are seeded by migrations.

## Make targets

```sh
make help              # list all targets
make dev               # live-reload server (air)
make run               # run once
make build             # build ./bin/api
make test              # go test ./...
make swagger           # regenerate docs/ from annotations
make migrate-up        # apply migrations
make migrate-down      # roll back one
make migrate-create name=add_x   # new migration
```

## Testing

```sh
make test
```

Middleware (CORS, rate limiting, JWT auth) has unit tests under
`internal/middleware`. Endpoints were verified end-to-end against a live
PostgreSQL instance.
