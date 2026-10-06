---
description: "Boundary verification — catching integration defects by reading both sides of a component boundary"
paths: "**/*_test.go,**/*test*.py,**/*test*.ts,**/*test*.rs,**/test/**,**/tests/**"
---

<!-- Source: revfactory/harness — Apache License 2.0 — see .claude/rules/moai/NOTICE.md -->

# Boundary Verification

Some defects live between components: each side passes its own tests and the pair fails at runtime. Type casts and generic parameters often hide them from the compiler. Find them by reading both sides of the boundary together and testing the pair, not each side alone.

Boundaries to check when a change touches one side:

- **API response ↔ client type.** The route's actual payload (`NextResponse.json(...)`, the handler's return) against the type the hook or client expects — including wrapping (`{ projects: [...] }` versus a bare array), required versus optional fields, and every shape an async flow returns (an initial `202 { status, taskId }` and the final result need a discriminated union, not one type).
- **Field naming across layers.** Database `snake_case`, API and frontend `camelCase`: the transform has to happen at exactly one layer, consistently.
- **Links ↔ routes.** Every `href`, `router.push()` and `redirect()` target maps to a real page file; route groups `(group)` drop out of the URL, dynamic segments `[id]` are variables.
- **Client calls ↔ API routes.** Every `fetch` path and method has a route that actually handles it, not an empty stub.
- **State map ↔ state updates.** Every allowed transition has code that performs it, every performed transition is allowed, and every intermediate state has an exit.
- **Casts.** `as T`, `as unknown as T` and unconstrained generics such as `fetchJson<T>()` assert a shape without checking it; each one needs a reason, and preferably a runtime parse (for example a schema validator) at the boundary.

Prove a boundary with a test that exercises both sides — an integration or end-to-end test that sends a real request through the route and consumes the real response — rather than mocking one side to match the other.
