# lib

Small framework-agnostic helpers.

- `utils.ts` — `cn(...inputs)`: merges Tailwind classes via `clsx` + `tailwind-merge`. Standard shadcn helper, used throughout `components/ui/*`.

Add new cross-cutting utilities (formatting, geo/coordinate helpers for the future map views, etc.) here rather than inline in components.
