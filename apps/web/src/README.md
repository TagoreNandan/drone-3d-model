# src

Application source.

- `main.tsx` — React entry, mounts `RouterProvider` from `router.tsx`
- `router.tsx` — route table (`createBrowserRouter`): `/` → `AppShell` (layout) with `Home` as the index route and a catch-all 404
- `app-shell.tsx` — shared layout: `ThemeProvider` (dark by default) → `Header` + `Outlet` + `Toaster`
- `index.css` — Tailwind entry point and CSS theme tokens (used by shadcn components)
- `components/` — shared UI components, see `components/README.md`
- `lib/` — small utilities, see `lib/README.md`
- `routes/` — route-level page components, see `routes/README.md`

No `pages` vs `features` split yet — this is intentionally flat since there's only one real route so far. Reassess the structure once map/telemetry/mission views are added (see `../MAPPING_ENGINES.md`).
