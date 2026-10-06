---
name: moai-domain-frontend
description: >
  Frontend implementation guidance for Next.js (App Router), React and Tailwind CSS in
  TypeScript, usually inside a pnpm/turbo monorepo: current framework facts that differ
  from older habits, design direction, and how to verify UI work. Not for component and
  state review criteria (moai-ref-react-patterns) or backend work.

when_to_use: >
  Load when building or changing pages, layouts, components, styling or client-side data
  flow in a web frontend.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Grep, Glob
user-invocable: false
metadata:
  version: "3.0.0"
  category: "domain"
  status: "active"
  updated: "2026-03-28"
  tags: "frontend, nextjs, react, tailwind, typescript"
  author: "MoAI-ADK Team"
---

# Frontend Guidance

## Start from what the project has

Read `package.json` (and the workspace's `pnpm-workspace.yaml` and `turbo.json`) for the framework versions actually installed before relying on any API below. In a monorepo, look for shared packages (a UI kit, design tokens, config presets, a `cn` helper) and use them instead of writing local copies. Reuse existing components and hooks; a new primitive that duplicates one in the shared package is a finding.

## Design direction

Match the project's existing design system first: its tokens, spacing scale, type scale, component variants and motion. When a design file or reference is given, implement it faithfully, including spacing and states. When no visual direction is given, make a deliberate choice that fits the product and avoid the patterns that read as default generated UI: cream or off-white page backgrounds, an italic accent word inside a headline, numbered 01/02/03 section labels, monospace eyebrow labels, and pill-shaped buttons everywhere.

## Next.js App Router (16.x)

- `params` and `searchParams` are Promises: `await` them in server components and route handlers, and unwrap them with React's `use()` in client components. `cookies()`, `headers()` and `draftMode()` are async as well. The global `PageProps<'/route/[slug]'>`, `LayoutProps` and `RouteContext` helpers type them; `next typegen` regenerates them outside dev and build.
- `middleware.ts` is deprecated in favor of `proxy.ts` exporting `proxy`; it runs on the Node.js runtime. A project that needs the edge runtime keeps `middleware.ts`.
- Turbopack is the default bundler for `next dev` and `next build`. `next lint` is gone; lint with the ESLint CLI.
- Caching is opt-in through `cacheComponents: true` and the `'use cache'` directive with `cacheLife` and `cacheTag`. `revalidateTag(tag, profile)` takes a cache-life profile as its second argument; inside a server action, `updateTag(tag)` gives read-your-own-writes and `refresh()` re-renders the client router.
- Components are server components unless marked `'use client'`. Keep the client boundary as low in the tree as possible and pass only serializable props across it.
- Server actions and route handlers are public endpoints: validate input and check authorization in each one (see moai-ref-owasp-checklist).

## React (19.x)

- Use `use()`, server components, or the project's query library for data; `useEffect` plus `fetch` for initial data is a finding.
- Forms can use actions with `useActionState` and `useFormStatus`; optimistic UI uses `useOptimistic`.
- `useEffectEvent` (19.2) holds non-reactive logic used inside an effect; `<Activity mode="hidden">` keeps state for UI that is off screen.
- When the React Compiler is enabled in the project, do not add `useMemo`, `useCallback` or `memo` by hand without a measured reason.

## Tailwind CSS (4.x)

Configuration is CSS-first: `@import "tailwindcss";` plus design tokens declared in `@theme` blocks, rather than a `tailwind.config.js` theme. Use the project's tokens and utilities; arbitrary values such as `bg-[#1a1a1a]` or `mt-[13px]` where a token exists are findings. Check how the project defines its spacing and color scales before assuming the defaults.

## Verifying UI work

Type-check and lint with the workspace's scripts (see the TypeScript language rule for the turbo/pnpm commands). For visible changes, run the app and check the changed screens at mobile and desktop widths, in the themes the project supports, and in loading, empty and error states; a passing type-check does not show that a layout is right. State in the report what was checked in the browser and what was not.
