---
name: moai-foundation-core
description: >
  Shared MoAI working contract for the plan/run/sync agents: the SPEC lifecycle and
  where its artifacts live, the TRUST 5 checks every change must pass, and how work is
  delegated and carried across sessions. Not for SPEC authoring detail
  (moai-workflow-spec) or quality scoring (moai-foundation-quality).

when_to_use: >
  Background context for agents working inside a /moai plan, run or sync workflow.

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
  tags: "foundation, spec-lifecycle, trust-5, delegation"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 1500
---

# MoAI Foundation Core

The contract shared by the agents that author, implement, document and commit SPEC work. Each topic below names the file that owns it; read that file when you need more than this summary.

## SPEC lifecycle

Every SPEC moves through plan, run and sync. Artifacts live in `.moai/specs/SPEC-{ID}/`.

| Phase | Command | Agent | What it produces |
|-------|---------|-------|------------------|
| plan | `/moai plan` | manager-spec, then plan-auditor | `spec.md` and `plan.md` at every tier, `acceptance.md` from Tier M, `design.md` and `research.md` at Tier L; `spec.md` frontmatter `status: draft` |
| run | `/moai run SPEC-{ID}` | manager-develop | code and tests; `status: draft` to `in-progress` |
| sync | `/moai sync SPEC-{ID}` | manager-docs, sync-auditor, and manager-git when the route opens a PR | docs; `in-progress` to `implemented` to `completed` on the sync commit |

Run starts only after the user approves the Implementation Kickoff. The run cycle follows `development_mode` in `.moai/config/sections/quality.yaml`: `tdd` by default, `ddd` for existing code with little test coverage, and `autofix` for `/moai fix`.

Only manager-spec edits SPEC body content. manager-develop and manager-docs change only the `status:` and `updated:` frontmatter fields; when a SPEC body needs to change, they return a blocker report so the orchestrator can send it back to manager-spec. Phase progress is recorded in `.moai/specs/SPEC-{ID}/progress.md`, which `/moai run` and `/moai sync` read.

Where the details live: tiers and the commit/PR routes in `.claude/rules/moai/workflow/spec-workflow.md`; frontmatter keys, the status enum and transition ownership in `.claude/rules/moai/development/spec-frontmatter-schema.md`; requirement notation (GEARS) in the moai-workflow-spec skill.

## TRUST 5

The five questions every change answers before it is called done:

- Tested: the suite is green and coverage on changed code meets `constitution.test_coverage_target` in quality.yaml (85 by default). Untested code gets characterization tests before it is changed.
- Readable: the linter is clean, names say what things are, and comments follow the `code_comments` language in `.moai/config/sections/language.yaml`.
- Unified: the project's formatter has run and the change follows the conventions of the file it touches.
- Secured: input is validated at trust boundaries and no secret is written to a file. Use moai-ref-owasp-checklist when reviewing.
- Trackable: commits follow Conventional Commits and reference the SPEC ID.

The thresholds are configuration, not judgment: `lsp_quality_gates` in quality.yaml allows zero errors, type errors and lint errors in run, and at most 10 warnings in sync. `moai gate` (also run by the git pre-commit hook) detects the language from marker files and runs vet, typecheck, lint and test, skipping tools that are not installed. The full rubric and how sync-auditor scores it are in moai-foundation-quality.

## Delegation

MoAI subagents have no Agent tool and cannot ask the user anything. When a decision belongs to the user, a subagent stops and returns a blocker report that states the question and the options; the orchestrator asks.

## Context and session continuity

Before a `/clear` or a session boundary, in-flight state goes into the SPEC's `progress.md` and the orchestrator emits a paste-ready resume message. Thresholds and format are in `.claude/rules/moai/workflow/context-window-management.md` and `.claude/rules/moai/workflow/session-handoff.md`.
