---
name: moai-workflow-testing
description: >
  How MoAI projects organize, run and measure tests: coverage targets and test-quality
  settings from quality.yaml, characterization tests before changing untested code, and
  the evidence a run or sync audit expects. Not for choosing a test mix (moai-ref-testing-pyramid)
  or the TDD/DDD cycle itself (moai-workflow-tdd, moai-workflow-ddd).

when_to_use: >
  Load when writing or changing tests, measuring coverage, or judging whether a change is
  adequately tested.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Write, Edit, Bash(pytest:*), Bash(ruff:*), Bash(npm:*), Bash(npx:*), Bash(node:*), Bash(jest:*), Bash(vitest:*), Bash(go:*), Bash(cargo:*), Bash(mix:*), Bash(uv:*), Bash(bundle:*), Bash(php:*), Bash(phpunit:*), Grep, Glob
user-invocable: false
metadata:
  version: "3.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-07-10"
  modularized: "false"
  tags: "workflow, testing, coverage, characterization-tests"
  author: "MoAI-ADK Team"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 1500
---

# Testing in MoAI Projects

## Where the numbers come from

Coverage and test-quality expectations are configuration in `.moai/config/sections/quality.yaml` (under `constitution:`). Read the project's file rather than assuming defaults:

- `test_coverage_target` (85 by default): coverage required on the code a change touches. manager-develop reports it per changed file or package; sync-auditor checks it.
- `tdd_settings.min_coverage_per_commit` (80 by default) and `tdd_settings.test_first_required`: in `tdd` mode the failing test comes before the implementation.
- `ddd_settings`: `require_existing_tests`, `characterization_tests`, `behavior_snapshots`, `preserve_before_improve`, and `max_transformation_size` (small by default) govern refactoring of existing code.
- `coverage_exemptions`: off by default. When enabled, an exemption needs a written justification and the exempt share stays under `max_exempt_percentage`.
- `test_quality`: tests assert specified behavior through public interfaces, carry meaningful assertions, and avoid coupling to implementation details.

## Organizing tests

Follow the layout the project already uses. Absent one, use the language convention: Go `*_test.go` in the same package, TypeScript `*.test.ts(x)` next to the source, Python `tests/test_*.py`. Name tests after the behavior they pin down. Each test sets up its own fixtures (temporary directories, fresh state) so tests can run in parallel and in any order.

## Running tests

Run the project's own scripts (package.json, Makefile, turbo pipeline) in preference to invoking a runner directly, and scope runs to what changed while iterating. The per-language commands, including the monorepo rule for turbo/pnpm workspaces, are in the language rules under `.claude/rules/moai/languages/`.

`moai gate` is the commit-time gate (the git pre-commit hook calls it). It detects the language from marker files and runs vet, typecheck, lint and test; step timeouts and the pre-commit scope are in `.moai/config/sections/gate.yaml`. A gate pass is not a substitute for running the tests that cover your change.

## Characterization tests (ddd cycle)

Before changing code that has little or no test coverage, capture what it does now:

1. Write tests that record current outputs for representative and edge inputs, including current error behavior, even where that behavior looks wrong. Note suspected bugs instead of fixing them in the same step.
2. Run them against the unchanged code; all must pass. This is the baseline.
3. Make the change in small transformations, re-running the characterization tests after each one.
4. When a behavior change is intended, change the characterization test in the same commit and say why in the commit body.

For complex outputs, snapshot or golden-file comparisons are acceptable when the snapshot is reviewed, not just regenerated. The naming convention and cycle detail are in moai-workflow-ddd.

## What a reviewer checks

- Every changed behavior has a test that would fail without the change (tdd), or characterization tests pass before and after (ddd).
- Error and boundary paths are tested, not only the happy path.
- Coverage on changed code meets `test_coverage_target`, shown with the tool's output.
- No skipped or disabled test without a stated reason; no test that only passes on retry.
- Concurrent code is exercised with the race detector where the language has one (`go test -race`).
- The report names the commands that were run and their result.
