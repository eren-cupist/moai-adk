---
paths: "**/*.ts,**/*.tsx,**/tsconfig.json"
---

# TypeScript

Commands (Node toolchain, also used for JavaScript):

- Use the package manager the lockfile names: `pnpm-lock.yaml` means pnpm. In a pnpm workspace never run `npm`, because it cannot resolve `workspace:*` links.
- Prefer the workspace's own scripts (`typecheck`, `lint`, `test`) over invoking `tsc`, `eslint` or `vitest` directly.
- When `turbo.json` exists at the repo root, verify with `pnpm turbo run <task> --affected` for changed packages or `pnpm --filter <pkg> <task>` for one package. A bare root `npm test` or `vitest run` skips turbo's `^build` and `^typecheck` ordering and fails on cross-package imports that have nothing to do with the change. The `turbo` binary is usually a local devDependency, so prefix it with `pnpm`.
- `moai gate` resolves the typecheck step in this order: `gate.typecheck.command` in `.moai/config/sections/gate.yaml`, then a `typecheck` script in package.json, then `npx --no-install tsc --noEmit` when a non-solution tsconfig.json is present. A monorepo root should define a `typecheck` script that delegates to turbo.

Conventions:

- No `any`; use `unknown` and narrow it. Use `@ts-expect-error` with a reason instead of `@ts-ignore`.
- Import only from packages declared in the importing package's own package.json. pnpm's strict layout makes imports of transitive dependencies, including type-only imports, fail in CI even when they resolve locally.
- Follow the tsconfig's module settings (`verbatimModuleSyntax`, path aliases) rather than rewriting import style.
- Re-run the typecheck for every package whose files changed before reporting done.

Frontend framework guidance (Next.js, React, Tailwind) is in the moai-domain-frontend skill.
