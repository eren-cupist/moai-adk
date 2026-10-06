---
name: moai-workflow-tdd
description: >
  MoAI conventions for the TDD (RED-GREEN-REFACTOR) cycle used when
  quality.yaml development_mode is tdd or a run adds new behavior. Not for
  changing untested legacy code (use moai-workflow-ddd).

when_to_use: >
  Use when implementing new behavior test-first in a /moai run or /moai fix,
  or when a manager-develop brief sets cycle_type=tdd.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Write, Edit, Bash(pytest:*), Bash(ruff:*), Bash(npm:*), Bash(npx:*), Bash(node:*), Bash(jest:*), Bash(vitest:*), Bash(go:*), Bash(cargo:*), Bash(mix:*), Bash(uv:*), Bash(bundle:*), Bash(php:*), Bash(phpunit:*), Grep, Glob
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-03"
  modularized: "false"
  tags: "workflow, tdd, test-driven, red-green-refactor, test-first"
  author: "MoAI-ADK Team"
  related-skills: "moai-workflow-ddd, moai-workflow-testing, moai-foundation-quality"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000
---

# TDD in MoAI

Selected by `.moai/config/sections/quality.yaml`:

```yaml
constitution:
  development_mode: tdd    # or ddd
  test_coverage_target: 85
  tdd_settings:
    test_first_required: true
    min_coverage_per_commit: 80
```

`tdd` is the default for new behavior, including new code inside an existing project. Edits to existing code that has no tests get characterization tests first (`moai-workflow-ddd`), whichever mode is configured.

## The cycle, with MoAI's evidence requirements

- **RED** — add a test for one behavior of an acceptance criterion, named for the behavior (`rejects an expired token`, not `returns false`). Run it and keep its failing output: that output, captured before the implementation exists, is the run's proof that the test came first and is reported in `progress.md` and in a manager-develop report. A new test that passes on its first run tested nothing new; find out why before moving on.
- **GREEN** — implement the general behavior the test describes. Do not special-case the fixture values.
- **REFACTOR** — remove duplication and clarify names with the suite green; run the affected tests after each step.

If implementation code already exists before its failing test, delete it and redo it test-first; a test written afterwards cannot produce the RED evidence.

## Done

- Every acceptance criterion in scope has at least one test that was observed failing and now passes.
- The full relevant suite passes; no test was skipped, disabled or weakened to get there.
- Coverage of changed code meets `test_coverage_target` (per-commit minimum `min_coverage_per_commit`) where the project measures coverage.
- New code carries `@MX` tags where the protocol calls for them (`@MX:ANCHOR` on functions with expected fan-in of 3 or more, `@MX:WARN` on concurrency or high complexity); `@MX:TODO` tags added during RED are removed once their tests pass.
