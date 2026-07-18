# MIVA FixIt — Frontend

Vue 3 + TypeScript + Vite SPA for the University Maintenance & Service Request
Management System (MIT 8333). Vercel-style dark UI, role-based dashboards for
**Student/Staff**, **Maintenance Officer**, and **Administrator**.

## Tech

- **Vue 3** (Composition API, `<script setup>`) + **TypeScript**
- **Vite** + **Tailwind CSS v4** (Vercel-style dark theme)
- **Pinia** + **pinia-plugin-persistedstate** (persisted auth session)
- **Vue Router** with role-based navigation guards
- **Axios** (JWT interceptor, ready for the Go backend)
- **Yup** form validation
- **vue3-toastify** notifications
- **Vitest** + Vue Test Utils (uses `happy-dom`)

## Mock backend

Until the Go API is ready, the UI runs against an in-memory mock
(`src/api/mock`, `src/api/services.ts`) persisted to `localStorage`, so every
flow is fully demoable. When the backend exists, set `VITE_USE_MOCK=false` and
point `VITE_API_BASE_URL` at it — the service layer swaps to the axios client in
`src/api/http.ts`.

## Demo accounts

All use password **`password123`** (also selectable on the login screen):

| Role | Email |
|------|-------|
| Student / Staff | `student@miva.edu` |
| Maintenance Officer | `officer@miva.edu` |
| Administrator | `admin@miva.edu` |

## Scripts

```sh
npm run dev          # dev server
npm run build        # type-check + production build
npm run test:unit    # unit tests (Vitest)
npm run lint         # lint
```

> **Node:** requires **20.19+** or **22.12+** (Vite 8). On older Node the dev
> server still runs but prints a warning; tests use `happy-dom` to stay
> compatible.

## Structure

```
src/
├── api/            # http client, mock db, service layer
├── components/     # ui/ (design system), layout/, requests/
├── lib/            # constants, validation (yup), format, csv/pdf export
├── router/         # routes + role guards
├── stores/         # pinia: auth (persisted), notifications
├── types/          # domain types
└── views/          # auth/, student/, officer/, admin/, RequestDetailView
```
