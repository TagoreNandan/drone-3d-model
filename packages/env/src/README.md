# src

- `web.ts` — client env for Vite (`VITE_` prefix), exported as `@gcs/env/web`. See `../README.md` for usage status.
- `server.ts` — server env (`DATABASE_URL`, `CORS_ORIGIN`, `NODE_ENV`) — not currently consumed anywhere (the backend is Go, not this package).
- `native.ts` — Expo/mobile env (`EXPO_PUBLIC_SERVER_URL`) — dead code, no native app exists in this repo.

Each file is a standalone `createEnv(...)` call from `@t3-oss/env-core`; there's no shared/base config between them.
