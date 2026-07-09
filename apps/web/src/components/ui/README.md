# components/ui

[shadcn](https://ui.shadcn.com/)-style primitives, copied in via the `shadcn` CLI per `../../components.json` (style: "nova"), built on `radix-ui`/`@base-ui/react` + `class-variance-authority` + `tailwind-merge`.

Present: `button`, `card`, `checkbox`, `dropdown-menu`, `input`, `label`, `skeleton`, `sonner` (toast).

Most of these are currently unused boilerplate — only `button` and `dropdown-menu` are actually referenced (by `mode-toggle.tsx`) and `sonner` (by `app-shell.tsx`'s `<Toaster />`). Add more primitives with the shadcn CLI as real UI is built, rather than hand-writing new ones — keeps styling/variants consistent with what's already here.
