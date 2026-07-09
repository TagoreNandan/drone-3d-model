# apps/web — GCS frontend

React + Vite SPA for the Ground Control Station. Currently a minimal shell — routing, theming, and a backend health widget — with no actual GCS features (map view, telemetry dashboard, vehicle list, mission editor) built yet.

## Stack

- **React 19** + **react-router 8** (`createBrowserRouter`)
- **Vite 8** — dev server on port **5173**
- **Tailwind CSS 4** (via `@tailwindcss/vite`) + **shadcn/radix-ui** components (`components.json`, "nova" style)
- **Zustand** — installed, not yet used anywhere
- **@tanstack/react-form** + **@hookform/resolvers** — installed, not yet used
- **next-themes** — dark/light theme provider (defaults to dark)
- **sonner** — toast notifications
- **@gcs/env** (workspace package) — typed client env access
- **Testing**: Vitest 4 + Testing Library + jsdom (`pnpm test`)

## Directory Map

```text
apps/web/
├── index.html
├── vite.config.ts
├── components.json           # shadcn config
├── GRAPH_BACKEND.md          # backend URL/health-check reference
├── MAPPING_ENGINES.md        # IMPORTANT: planned map stack, read before building map views
├── .env                      # VITE_SERVER_URL=http://localhost:8080
└── src/
    ├── main.tsx               # mounts RouterProvider
    ├── router.tsx             # route table
    ├── app-shell.tsx          # ThemeProvider + Header + Outlet + Toaster
    ├── index.css              # Tailwind entry + theme tokens
    ├── components/            # Header, theme toggle, backend status widget, shadcn ui/*
    ├── lib/                   # utils.ts (cn() helper)
    └── routes/                # route-level page components (currently just home.tsx)
```

## Running

```sh
pnpm install
pnpm dev        # vite dev, http://localhost:5173
```

Other scripts: `pnpm build` (vite build), `pnpm serve` (vite preview), `pnpm check-types` (tsc --noEmit), `pnpm test` (vitest run).

## Current implementation state

- `main.tsx` → `router.tsx` (routes: `/` → `Home`, `*` → 404) → `app-shell.tsx` (theme provider + header + outlet + toaster)
- `routes/home.tsx` renders an ASCII banner and `<GraphBackendStatus />`, which pings `${VITE_SERVER_URL}/health` and shows connected/disconnected
- shadcn `ui/*` primitives (button, card, checkbox, dropdown-menu, input, label, skeleton, sonner) exist but are largely unused boilerplate — no map, telemetry, or mission UI yet

## Before building the Fly View / 3D View

Read **`MAPPING_ENGINES.md`** first. Two separate mapping engines are planned (not merged into one map layer):

- **mapcn.dev (MapLibre GL)** for the primary "Fly View" — `pnpm add maplibre-gl`, then install mapcn.dev's copy-in components via its own CLI
- **CesiumJS** (`pnpm add cesium resium`) for a lazy-loaded "3D View"
- **msgpackr** (`pnpm add msgpackr`) for decoding the binary WebSocket telemetry stream (once the backend's `/ws` route + MAVLink ingestion exist — see `apps/server/internal/realtime` and `apps/server/internal/mavlink`)

None of `maplibre-gl`, `cesium`, `resium`, or `msgpackr` are installed yet.
