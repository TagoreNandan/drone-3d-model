# @gcs/env

Typed environment variable access via [`@t3-oss/env-core`](https://env.t3.gg/) + Zod. Workspace package, imported as `@gcs/env` (currently only the `./web` export is declared in `package.json`).

## `src/web.ts` (exported as `@gcs/env/web`)

Client-side env for the Vite frontend. `clientPrefix: "VITE_"`, no required client vars declared yet — reads from `import.meta.env`. Not currently imported anywhere in `apps/web` (the frontend reads `import.meta.env.VITE_SERVER_URL` directly in `graph-backend-status.tsx` instead); wire this in if you want validated/typed env access there.

## `src/server.ts` — not wired up

Declares required `DATABASE_URL` (string) and `CORS_ORIGIN` (url), optional `NODE_ENV` (default `development`), read from `process.env`. **Not imported by `apps/server`** — that's a Go module and can't consume a TS package. This was likely intended for a Node/Bun-based backend that doesn't exist in this project graph; either repurpose it or leave unused until a Node service needs it.

## `src/native.ts` — dead code

Declares `EXPO_PUBLIC_SERVER_URL` for a mobile/Expo app. **No such app exists under `apps/`** — this is scaffold carried over from the generator's broader template library. Safe to delete unless a native app gets added later.

## Adding an export

Update the `exports` map in `package.json` (currently only `"./web": "./src/web.ts"`) if you start using `server.ts` or `native.ts` from somewhere.
