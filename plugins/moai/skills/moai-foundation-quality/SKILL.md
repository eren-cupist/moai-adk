---
name: moai-foundation-quality
description: >
  TRUST 5 quality checklist (Tested, Readable, Unified, Secured, Trackable)
  with the thresholds and per-language toolchains MoAI checks against. Used
  by sync-auditor and code review to judge a change. Not a linter or a
  library.

when_to_use: >
  Use when judging whether a change meets MoAI's quality bar: auditing a
  SPEC implementation, reviewing code, checking coverage and lint
  thresholds, or picking the right toolchain command for a language.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Grep, Glob
user-invocable: false
metadata:
  version: "4.0.0"
  category: "foundation"
  status: "active"
  updated: "2026-07-10"
  modularized: "false"
  tags: "foundation, quality, testing, validation, trust-5, code-review"
  aliases: "moai-foundation-quality"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000
---

# TRUST 5

TRUST 5 is the checklist MoAI applies to a change. MoAI ships no quality library: the checks are the project's own toolchain run against the change, plus judgment. Ask each question of the change in front of you, and when one does not apply, say why ("docs-only change, no executable code") — an unexplained "not applicable" is an unverified item, not a pass.

## Checklist

**Tested**
- The changed behavior has tests, and they pass (run them; quote the output).
- The tests assert the behavior — they would fail if the code were wrong — rather than just executing it.
- Coverage meets `quality.yaml` `test_coverage_target` (85% by default) or the evaluator profile's threshold.
- A bug fix comes with a test that reproduced the bug before the fix.
- Existing untested code the change touches has characterization tests capturing its prior behavior.

**Readable**
- Names say what things are; nothing needs decoding.
- Comments explain the non-obvious why, in the configured code-comment language.
- The implementation is the simplest correct one — no abstraction without a second use.

**Unified**
- The change follows the file's and codebase's existing conventions: naming, error handling, imports, structure.
- Formatter and linter pass; suppressions are per-line with a reason, never global.

**Secured**
- External input is validated at trust boundaries.
- Web and API code is free of the OWASP-relevant issues: injection, broken authentication or authorization, XSS, CSRF, SSRF, unsafe deserialization, sensitive data in logs or responses. Load `moai-ref-owasp-checklist` for the details.
- No credentials in code or history; secrets come from the environment or a secrets manager.
- Destructive operations require explicit confirmation.

**Trackable**
- Commits follow the project's commit convention and name the SPEC or issue they implement.
- The change traces back to a requirement through the commit, the SPEC, or a linked issue.
- @MX tags are present where the protocol calls for them (`.claude/skills/moai/references/mx-tag.md`).

## How the checklist is applied

`sync-auditor` judges four areas — Functionality (maps to Tested), Security (Secured), Craft (Readable and Unified), Consistency (Unified and Trackable) — against the evaluator profile in `.moai/config/evaluator-profiles/` (selected by the SPEC's `evaluator_profile` or `harness.yaml` `harness.default_profile`). The profile names the must-pass areas and thresholds; by default Functionality and Security are must-pass. Whether the auditor runs at all depends on the harness level: `harness.yaml` `levels.<level>.evaluator` is false at `minimal` and true at `standard` and `thorough`.

`/moai review` uses the same questions to find defects in a diff, reporting every finding with severity and confidence before ranking them. `manager-develop` builds the Tested and Unified items in during `/moai run` through its `cycle_type` (`tdd`, `ddd`, or `autofix` for `/moai fix`). The agents themselves are listed in CLAUDE.md, which holds the roster for the 7-agent catalog.

## Thresholds

LSP gates from `quality.yaml` `lsp_quality_gates`:

| Phase | Requirement |
|-------|-------------|
| plan | capture the LSP baseline |
| run | zero errors, zero type errors, zero lint errors |
| sync | zero errors, at most 10 warnings, clean LSP |

Coverage: `test_coverage_target` (default 85%).

## Toolchains

Detect each language present from its markers and use the project's own scripts when it has them (`package.json` scripts, a `Makefile`, turbo tasks). A tool that is not installed is reported as skipped, not as passing.

| Language | Marker | Lint | Format | Type check | Test / coverage |
|----------|--------|------|--------|-----------|------------------|
| Go | `go.mod` | `go vet`, `golangci-lint run` | `gofmt -l` | (compiler) | `go test ./...`, `go test -coverprofile=coverage.out ./...` |
| TypeScript | `tsconfig.json` | `eslint` | `prettier --check` | `tsc --noEmit` | `vitest run --coverage` or `jest --coverage` |
| JavaScript | `package.json` | `eslint` | `prettier --check` | — | `vitest` / `jest` |
| Python | `pyproject.toml`, `requirements.txt` | `ruff check` | `ruff format --check` or `black --check` | `mypy` | `pytest --cov` |

In a pnpm workspace with `turbo.json`, run tasks through turbo (`pnpm turbo run lint test --affected`, or `pnpm --filter <pkg> <task>`) so build dependencies run first. Other languages: use the build file's standard commands. Further detail on toolchains and the evaluator profiles: [references/reference.md](references/reference.md).
