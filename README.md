# University Maintenance & Service Request Management System (UMSRMS)

A full-stack web application for **MIT 8333 — Advanced Web Application Development**.
It replaces the university's manual maintenance workflow (phone calls, paper forms,
WhatsApp, office visits) with a digital platform where students/staff submit and
track service requests, maintenance officers work assigned jobs, and administrators
manage users, assign tasks, monitor status, and generate reports.

## Roles

| Role | Can |
|------|-----|
| **Student / Staff** | Register, submit requests, upload fault evidence, track their requests, receive notifications |
| **Maintenance Officer** | View assigned requests, update progress, mark jobs complete |
| **Administrator** | Manage users & roles, assign requests to officers, monitor all requests, view reports & audit trail |

## Tech stack

**Backend** — Go (Gin) · PostgreSQL (Squirrel + pgx) · JWT + bcrypt · Swagger · golang-migrate
**Frontend** — Vue 3 + TypeScript · Vite · Pinia (persisted) · Tailwind CSS v4 · Yup · vue3-toastify (Vercel-style dark UI)

## Repository layout

```
miva_projects/
├── be/          # Go REST API  — see be/README.md
├── fe/          # Vue 3 + TS SPA — see fe/README.md
└── README.md    # you are here
```

Each package has its own README with detailed setup:
- **[Backend →](be/README.md)**
- **[Frontend →](fe/README.md)**

## Quick start

### 1. Backend

```sh
cd be
cp .env.example .env            # edit DB creds + JWT_SECRET
createdb umsrms                 # create the (empty) database once
make dev                        # migrations run automatically on boot
                                # http://localhost:8080  (Swagger at /swagger/index.html)
```

### 2. Frontend

```sh
cd fe
npm install
npm run dev                     # http://localhost:5173
```

By default the frontend runs against an **in-memory mock** so it works without the
backend. To use the real API, set `VITE_USE_MOCK=false` in `fe/.env` (with the
backend running).

## Feature coverage (assignment requirements)

| Requirement | Where |
|-------------|-------|
| Registration/login, role dashboards, request form, tracking, admin UI | Frontend (`fe/`) |
| Auth & authorization, REST APIs, CRUD, RBAC, request assignment, validation | Backend (`be/`) |
| Relational DB with 6+ entities & relationships | PostgreSQL (8 tables) |
| JWT auth · RBAC · file/image upload · in-app notifications · search/filter/pagination · audit trail · Swagger docs · reports | Implemented (8 of the 9 advanced options) |
| Testing | Frontend (Vitest) + backend (Go tests) |

## Architecture

```
┌──────────────────────┐      HTTPS/JSON      ┌───────────────────────────┐      SQL      ┌────────────┐
│  Vue 3 + TS (Vite)   │ ───────────────────▶ │  Go API (Gin)             │ ────────────▶ │ PostgreSQL │
│  Pinia · Tailwind    │                      │  JWT · RBAC · Swagger     │               │            │
└──────────────────────┘                      └───────────────────────────┘               └────────────┘
```

## Demo accounts (mock mode)

When `VITE_USE_MOCK=true` (default), sign in with password **`password123`**:

| Role | Email |
|------|-------|
| Student / Staff | `student@miva.edu` |
| Maintenance Officer | `officer@miva.edu` |
| Administrator | `admin@miva.edu` |

Against the real backend, register a new account (the sign-up form lets you pick an
account type).

## Course

MIT 8333 — Advanced Web Application Development (Virtual Lab), MIVA Open University.
