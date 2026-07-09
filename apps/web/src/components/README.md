# components

Shared, non-route-specific UI.

- `header.tsx` — top nav bar (currently just a "Home" link) + `<ModeToggle />`
- `mode-toggle.tsx` — light/dark/system theme dropdown, built on `theme-provider.tsx` + shadcn `Button`/`DropdownMenu`
- `theme-provider.tsx` — thin re-export wrapper around `next-themes`' `ThemeProvider`/`useTheme`
- `graph-backend-status.tsx` — `<GraphBackendStatus />`: polls `${VITE_SERVER_URL}/health` once on mount and renders a checking/connected/disconnected dot + label. This is currently the only piece of the frontend that talks to the backend.
- `loader.tsx` — generic centered spinner (`Loader2` from lucide-react), not currently used by any route
- `ui/` — shadcn/radix primitives, see `ui/README.md`

None of these are GCS-domain components yet (no vehicle card, telemetry readout, map marker, etc.) — this directory is still just app-shell chrome.
