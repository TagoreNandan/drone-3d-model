# @gcs/config

Shared TypeScript config for the JS workspace.

- `tsconfig.base.json` — strict TS, ESNext target, bundler module resolution. Consumed via `extends` from the root `tsconfig.json`, `apps/web/tsconfig.json`, and `packages/env/tsconfig.json`.

No source code, no build scripts — this package only exists to be extended, not run or imported at runtime.
