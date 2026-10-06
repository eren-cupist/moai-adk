---
name: moai-workflow-ddd
description: >
  MoAI conventions for the DDD (ANALYZE-PRESERVE-IMPROVE) cycle: characterize
  existing behavior with tests, then change it, used when quality.yaml
  development_mode is ddd or a change touches untested code. Not for purely
  new behavior (use moai-workflow-tdd).
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Write, Edit, Bash(git:*), Bash(pytest:*), Bash(ruff:*), Bash(npm:*), Bash(npx:*), Bash(node:*), Bash(uv:*), Bash(make:*), Bash(cargo:*), Bash(go:*), Bash(mix:*), Bash(bundle:*), Grep, Glob
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-01-16"
  modularized: "false"
  tags: "workflow, refactoring, ddd, behavior-preservation, characterization-tests"
  author: "MoAI-ADK Team"
  related-skills: "moai-workflow-testing, moai-foundation-quality"
---

# DDD in MoAI

DDD here means characterize, then change: before modifying code that has little or no test coverage, pin its current behavior with tests, then change it in small steps while those tests keep passing. Selected by `.moai/config/sections/quality.yaml`:

```yaml
constitution:
  development_mode: ddd    # or tdd
  ddd_settings:
    require_existing_tests: true
    characterization_tests: true
    behavior_snapshots: true
    max_transformation_size: small
    preserve_before_improve: true
```

## The cycle

- **ANALYZE** — read the code the SPEC changes and its callers; note which behavior the SPEC keeps and which it changes, and the coupling that makes a change risky. Read the `@MX:ANCHOR` and `@MX:WARN` tags in those files.
- **PRESERVE** — confirm the existing tests pass, then add characterization tests for the paths you will touch: run the code, record what it actually does (not what it should do), and assert that. Name them `test_characterize_<component>_<scenario>` (or the project's equivalent, e.g. `describe('characterize: <component>')`). For large outputs use snapshot tests. A surprising behavior you discover is reported, not silently fixed.
- **IMPROVE** — make one small change at a time (`max_transformation_size: small`) and run the affected tests after each; revert a step that breaks a characterization test unless the SPEC changes that behavior, in which case update that test deliberately and say so. Independent packages can progress in parallel; within one package, one change at a time.

Behavior the SPEC adds is built test-first (`moai-workflow-tdd`) on top of the safety net.

## Done

- Characterization tests cover the changed paths and passed before the change.
- After the change, every characterization test passes, except those the SPEC intentionally changes, which were updated deliberately and are listed in the report.
- The full relevant suite passes and no existing test was deleted or weakened.
- The report names the behaviors preserved, the behaviors changed (with the acceptance criterion behind each), and any surprising behavior found while characterizing.
- `@MX:ANCHOR` tags are updated where fan-in changed, and `@MX:NOTE` added for business rules discovered during ANALYZE.
