# routes

Route-level page components, wired up in `../router.tsx`.

- `home.tsx` — the only route so far (`/`). Renders an ASCII-art banner and an "API Status" card showing `<GraphBackendStatus />` plus a placeholder note ("No API client selected for this frontend.").

This is where the planned Fly View (MapLibre/mapcn.dev), 3D View (CesiumJS), vehicle list, and mission editor routes will go — see `../MAPPING_ENGINES.md` before starting on any map-related route. Add each new page as its own file here and register it in `../router.tsx`.
