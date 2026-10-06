---
description: "Sync Phases 1-10 — pre-sync check, deployment readiness (migrations, environment, dependencies, compatibility), quality verification, security review, @MX validation, and coverage."
user-invocable: false
metadata:
  parent: moai-workflow-sync
  phase: "Phase 1-10: Checks and Quality Verification"
---

# Sync Phases 1-10: Checks

These phases establish that what `/moai run` shipped works and is ready to document and deliver. In `status` mode they run read-only: no auto-fix, no generated tests, no file or git writes — print the report after Phase 10 and stop.

Scope follows the mode: `auto` checks the changed packages and what imports them (widen only with a stated reason); `force` and `project` run full-repository checks.

## Reusing recorded results

Before running a check, query the shared diagnostic snapshot with `moai verify check --key-current`. The key covers HEAD, the working-tree status, and the diff, so a new commit or edit invalidates it. When a fresh snapshot (matching key, within its TTL) already covers a check — recorded by run's pre-review gate or an earlier sync phase on the unchanged tree — cite it instead of re-running: give the snapshot path, key, original command, and exit code as the evidence. On a key mismatch or expiry, run the check and record the result with `moai verify record`. A stale snapshot is never cited.

## Language detection

Resolve every language present from repository markers and changed-file suffixes, and route each — a monorepo can have several. Common markers: `go.mod` (Go), `tsconfig.json` or `package.json` (TypeScript/JavaScript), `pyproject.toml` / `requirements.txt` / `setup.py` (Python), `Cargo.toml` (Rust); other languages by their usual build files. Use the project's own runners and scripts: in a pnpm workspace with `turbo.json`, run tasks through turbo (`pnpm turbo run test --affected`, or `pnpm --filter <pkg> test` for one package) rather than a bare test runner, which skips turbo's build ordering. A tool that is not installed is reported as skipped; an installed tool that fails is a failure.

## Phase 1: Pre-sync check (`gate-sync-1`)

The tree should contain only the expected changes, and lint, format, type check, and tests should pass. Run them in parallel. Outside `status` mode, apply the formatter's and linter's automatic fixes. If errors remain, ask the user: fix them (recommended — delegate to `manager-develop` with the failures, injecting the SPEC's cycle skill `moai-workflow-tdd` or `moai-workflow-ddd`), skip to Phase 2 knowing they will surface later, or abort.

## Phases 2-6: Deployment readiness

- **Phase 2 — Database.** Schema changes, new or altered models, migration files: list them as deployment prerequisites for the PR body.
- **Phase 3 — Tests.** The full suite for each detected language passes with zero failures. On failure, ask: fix and retry (recommended, via `manager-develop`), continue with a warning, or abort.
- **Phase 4 — Configuration and environment.** Environment variables the code reads that `.env.example` or the docs do not mention, new config files or keys, changed defaults: summarize for the PR body.
- **Phase 5 — Dependencies.** Changed manifests and lockfiles: list added, removed, and upgraded packages for the PR body. (The Stop hook separately records `deps_modified=1` when the HEAD commit changed a manifest; that is informational, not a vulnerability scan.)
- **Phase 6 — Backward compatibility.** Removed or changed public endpoints, exported functions and types, response shapes, and dependency changes that affect consumers. Classify each as breaking (needs a major version and documentation), deprecation (needs a migration note), or compatible. A breaking change needs the user's explicit acknowledgment before sync continues.

Summarize readiness as READY, NEEDS_ATTENTION, or BLOCKED. On BLOCKED, show the blocking items and stop unless the user overrides.

## Phase 7: Quality verification

Run tests, lint, and type check per detected language in parallel (or cite fresh snapshot results). LSP gates for sync come from `quality.yaml` `lsp_quality_gates.sync`: `max_errors: 0`, `max_warnings: 10`, `require_clean_lsp: true`. If tests fail here, ask: continue with the failure recorded, or abort.

Then run the independent audit described in `sync.md`, section "The independent audit (Phase 7)". The auditor may cite these results instead of re-running them.

## Phase 8: Security review

Applies when the change touches security-sensitive code: authentication, sessions, tokens, permissions, roles; database queries, models, migrations, repositories; API handlers, routes, middleware; user input handling and validation; configuration that carries credentials. Otherwise record "Security review skipped: no security-sensitive files changed."

When `sync-auditor` runs, its security area covers this review; brief it to look closely at these files. When it does not run, review them yourself with `Skill("moai-ref-owasp-checklist")` (and `Skill("moai-ref-secops")` for CI/CD, container, or infrastructure changes).

Critical and high findings block; medium and low findings are advisory and go in the sync report. For a blocking finding, ask: fix now (recommended — `manager-develop`, then re-check), continue under an approved exception (record finding ID, rationale, scope, approver, expiry, and review condition), or abort. An exception approved here never clears an audit FAIL.

## Phase 9: @MX tag validation

Skip with `--skip-mx` and note "MX validation skipped by user flag" in the report. Otherwise check the changed source files against the @MX protocol (`.claude/skills/moai/references/mx-tag.md`):

- P1 (blocking): an exported function with three or more callers lacks `@MX:ANCHOR`.
- P2 (blocking): a goroutine or async pattern lacks `@MX:WARN` — for example a Go goroutine without cancellation, a TypeScript async function or `Promise.all` without error handling, a Python `async def` or thread without error handling.
- P3 (advisory): a long exported function lacks `@MX:NOTE`.
- P4 (advisory): an untested public function lacks `@MX:TODO`.

Sync reports tags but does not edit source. P1 or P2 violations stop the sync before Phase 10 — the coverage run is the expensive part, so a missing tag should not abort it halfway. Show the violations with file:line and the message "Run /moai run to add missing tags, or use --skip-mx to bypass". P3 and P4 are listed in the report and do not stop anything.

## Phase 10: Coverage

Measure coverage for each detected language (`go test -coverprofile=coverage.out ./...` then `go tool cover -func=coverage.out`; `vitest run --coverage` or `jest --coverage`; `pytest --cov`) and compare with `quality.yaml` `test_coverage_target` (85 by default) or the evaluator profile's threshold.

Below target, and not in `status` mode: have `manager-develop` add tests for the riskiest gaps first — public API, functions with three or more callers or `@MX:ANCHOR`, error-handling paths — in the style the SPEC's development mode uses, then re-run the suite and re-measure. If the target is still not met, record the remaining gaps and continue; coverage does not block sync. Report before and after percentages, the tests added, and the remaining gaps per package.
