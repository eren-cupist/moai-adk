---
paths: "**/.moai/specs/**,**/.moai/config/sections/quality.yaml"
---

# SPEC Workflow

Non-trivial work moves through three phases — plan, run, sync — around one SPEC directory, `.moai/specs/SPEC-<ID>/`. This rule is the reference for what each phase produces, how SPECs are tiered, and the plan audit gate between plan and run. The workflow instructions themselves live in `.claude/skills/moai/workflows/plan.md`, `run.md`, and `sync.md`.

## Phase Overview

| Phase | Command | Agents | Produces |
|-------|---------|--------|----------|
| Plan | `/moai plan` | manager-spec writes, plan-auditor audits | the SPEC directory, a plan-audit report, the Implementation Kickoff Approval |
| Run | `/moai run SPEC-<ID>` | manager-develop (cycle_type from `constitution.development_mode` in quality.yaml: `ddd` or `tdd`) | code and tests, `progress.md` §E.2 and §E.3 |
| Sync | `/moai sync SPEC-<ID>` | manager-docs, sync-auditor, manager-git | updated docs, `status: completed`, `progress.md` §E.4, a PR on the PR route |

Per the retained catalog, the MoAI agent catalog consists of exactly 7 retained agents (`manager-spec`, `manager-develop`, `manager-docs`, `manager-git`, `plan-auditor`, `sync-auditor`, plus the built-in `Explore`).

## SPEC Phase Discipline

Step 1 (plan) runs in the main checkout, not in a worktree. Plan artifacts are markdown with nothing to conflict, and a plan written in the main checkout stays visible to other SPECs and to plan-auditor's cross-SPEC checks.

Steps 2 (run) and 3 (sync) run in the same place: the same branch, or the same worktree when the user chose one (`moai cc -w SPEC-<ID>` enters it). Sync rewrites docs in the tree that run modified, so a fresh worktree at sync would lose that state.

How a phase ends depends on the route, chosen by tier and the `--pr` flag:

- **Route A — trunk (default for Tier S and M):** each phase ends with its commits on the working branch; there is no per-phase PR. Pushing follows the git-strategy config and needs the user's approval.
- **Route B — PR route (Tier L, or `--pr`):** manager-git opens a PR per phase, and the phase ends when that PR merges. The merge method is `git_strategy.<mode>.merge_method` (`squash` | `merge` | `rebase`, default `squash`).

Step 4 (cleanup) applies to **Route B only**. It MUST happen ONLY after BOTH run AND sync PRs are merged, and only when a worktree was created: `moai worktree done SPEC-<ID>`, run from the host checkout. An unmerged worktree branch is the only copy of the work, so disposing of it early loses the work.

## Subcommand Classification

`/moai fix` is a pipeline — localize, repair, validate, with no model choosing the next step — and ignores `--mode` (`MODE_FLAG_IGNORED_FOR_UTILITY`). `/moai plan`, `/moai run`, and `/moai sync` are multi-agent workflows; they reject `--mode pipeline` with `MODE_PIPELINE_ONLY_UTILITY`. `/moai review` is not classified.

### Mode Dispatch

For `/moai run`, the mode comes from the `--mode` flag first, then `workflow.default_mode` in `.moai/config/sections/workflow.yaml`, then the default `autopilot`. Valid values are `autopilot`, `loop`, and `team`; any other value fails with `MODE_UNKNOWN`. `team` is experimental and used only on explicit request — see Agent Teams Variant.

## SPEC Complexity Tier (S/M/L)

Plan classifies every SPEC into a tier before writing artifacts. The tier sets the artifact set and the plan-auditor PASS threshold. `moai spec lint` reads the "Artifact set" column of this table to check each SPEC directory, so keep the table's shape.

| Tier | Scope guidance (LOC) | Files affected | Artifact set | plan-auditor PASS threshold |
|------|----------------------|----------------|--------------|------------------------------|
| S (Simple) | < 300 LOC | < 5 files | **2 files**: spec.md + plan.md (acceptance criteria inline in the spec) | 0.75 |
| M (Medium) | 300 - 1000 LOC | 5 - 15 files | **3 files**: spec.md + plan.md + acceptance.md | 0.80 |
| L (Large) | > 1000 LOC or constitutional | > 15 files | **5 files**: spec.md + plan.md + acceptance.md + design.md + research.md | 0.85 |

The LOC and file counts are guidance; a small change to a security boundary or a shared contract can still warrant Tier L. `progress.md` is written at every tier and is not counted. The tier is stored in the optional `tier:` frontmatter field; a SPEC without it is treated as Tier L.

Requirement and acceptance-criterion ceilings, each counted separately: S 8 and 8, M 16 and 16, L 25 and 25. Split a SPEC that exceeds either ceiling, or move it up a tier, rather than letting it grow — the auditor and the implementer both have to hold every requirement in view.

## Plan Phase

Plan produces a SPEC a developer can approve and `/moai run` can execute: GEARS requirements, acceptance criteria covering every requirement, the tier's artifact set, and an independent audit. Planning writes no implementation code.

## Run Phase

Run implements the SPEC with the methodology in `constitution.development_mode` (quality.yaml), through manager-develop:

- **DDD Mode** (ANALYZE-PRESERVE-IMPROVE) — for code with little test coverage: understand the existing behavior, pin it with characterization tests, then change it in small steps.
- **TDD Mode** (RED-GREEN-REFACTOR, the default) — write a failing test for the behavior, make it pass simply, then clean up with the tests green. On existing code, read the current behavior before writing the first test.

`moai init --mode <ddd|tdd>` and the `MOAI_DEVELOPMENT_MODE` environment variable override the config. Run is done when every acceptance criterion is met with evidence and the quality gates pass. When run stalls — several iterations with no new criterion met, or a requirement that cannot be met as written — manager-develop returns a blocker report, and the orchestrator asks the user whether to continue, revise the SPEC, or change approach.

## Sync Phase

Sync updates documentation to match the implementation, records the close in `progress.md` §E.4, moves the SPEC to `completed`, and opens the PR on Route B.

## Phase Transitions

Plan to Run:
- Pre-condition: `progress.md` §E.1 records `plan_complete_at` and `plan_status: audit-ready`, the latest plan-auditor verdict is PASS, and the user gave the Implementation Kickoff Approval. On Route B the plan PR is merged.
- `/moai run` starts with Phase 1, the Plan Audit Gate.
- Plan Audit Gate skip policy: run may reuse the existing plan-audit verdict instead of re-auditing when all three hold — (1) the verdict is `PASS`; (2) its overall score meets the SPEC's tier threshold (S 0.75, M 0.80, L 0.85; `internal/runtime.SkipEligibleByScore`); (3) the plan-artifact hash is unchanged since that verdict. Otherwise Phase 1 audits again. Reusing the verdict never skips the Kickoff Approval.

Run to Sync: run's commits are on the branch (Route A) or the run PR is merged (Route B), with tests passing. Sync runs where run ran.

Sync (close): the single sync commit carries the `implemented → completed` transition and fills `sync_commit_sha` in `progress.md` §E.4. A commit cannot contain its own hash, so the sync commit writes the placeholder `pending-backfill` (the suffixed form `pending-backfill-sync` is also accepted) and a following commit backfills the real SHA. Do not leave the slot empty: an empty value is neither a SHA nor a recognized placeholder, and once the SPEC is `completed` nothing else comes back to fill it.

## Phase 1: Plan Audit Gate

The first step of every `/moai run SPEC-<ID>`: before any implementation, the plan-auditor verdict for the current plan artifacts must be PASS. No harness level disables it.

Before the audit, check `depends_on:` in spec.md. Every listed SPEC must have `status: completed`. If one does not, ask the user to wait, override (`--ignore-deps`, with the unfulfilled IDs and the reason logged to `.moai/logs/depends-on-override.log`), or abort.

### Verdicts

| Verdict | Meaning | Action |
|---------|---------|--------|
| `PASS` | No must-pass failure, no blocking finding, score at or above the tier threshold | Record in progress.md and continue |
| `FAIL` | A must-pass failure or a blocking finding | Stop, show the report, ask the user how to proceed |
| `BYPASSED` | The user passed `--skip-audit` or set `MOAI_SKIP_PLAN_AUDIT=1` | Record the bypass and continue |
| `INCONCLUSIVE` | The audit timed out, errored, or returned output that does not parse | Never treated as PASS; ask the user to retry, proceed, or abort |

### Report Persistence

- **Plan-phase review reports** — `.moai/reports/plan-audit/<SPEC-ID>-review-<N>.md`, one file per audit round, written by plan-auditor. `runtime.ResolveLatestPlanAudit` and `moai plan render-html` read the highest `<N>`. The report's machine lines are listed in `.claude/agents/moai/plan-auditor.md`, section "Report".
- **Run-gate record** — `.moai/reports/plan-audit/<SPEC-ID>-<YYYY-MM-DD>.md`, appended by the run gate (`internal/runtime/audit_report.go`) on every gate call. It is history only, never the source of a verdict or a cache key.

The plan-artifact hash (`internal/runtime/audit_cache.go` `ComputeHash`) covers the exact bytes of whichever of these exist in the SPEC directory: `acceptance.md`, `decision-index.md`, `design.md`, `plan.md`, `research.md`, `spec.md`, `tasks.md`. Editing any of them — including an in-place amendment of a completed SPEC — invalidates a cached PASS and forces a fresh audit on the next run.

Both report streams are local, gitignored artifacts.

### Grace Window

For 7 days after the time recorded in `.moai/state/audit-gate-merge-at.txt` (ISO-8601), a FAIL verdict only warns (`FAIL_WARNED`); after that it blocks.

## Agent Teams Variant

Agent Teams is experimental. It is used only on an explicit `--team` or `--mode team` request and needs `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`; nothing selects it automatically. The default is one manager-develop agent working serially, with read-only research fanned out to parallel subagents where the tracks are genuinely independent.
