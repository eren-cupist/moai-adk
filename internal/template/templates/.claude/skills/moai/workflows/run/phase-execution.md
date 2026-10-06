---
description: "Run phase — Plan Audit Gate, methodology routing (TDD/DDD), and brownfield delta markers"
user-invocable: false
metadata:
  parent: moai-workflow-run
  phase: "Plan Audit Gate and methodology routing"
---

# Phase 1: Plan Audit Gate

Every `/moai run` confirms that an independent plan-auditor has passed the SPEC's current plan artifacts before any code is written. An unreviewed or changed plan does not enter implementation.

## Reuse a prior PASS when nothing changed

Find the highest-numbered plan-phase review, `.moai/reports/plan-audit/<SPEC-ID>-review-<N>.md`. Its result can be reused, and the gate is met without a new audit, only when all three conditions of the skip policy in `.claude/rules/moai/workflow/spec-workflow.md` (section "Phase Transitions", Plan Audit Gate skip policy) hold: the verdict is PASS, the overall score meets the tier threshold (S 0.75, M 0.80, L 0.85), and the plan artifacts are unchanged since that review. The review records `plan_artifact_hash`; treat any edit to `spec.md`, `plan.md`, `acceptance.md`, `design.md`, `research.md`, `tasks.md` or `decision-index.md` after the review (check `git log` and `git status` on the SPEC directory) as a change. A review without hash or score fields is not reusable. Record the reuse in `progress.md` (`audit_cache_hit: true`, the review path).

## Otherwise run plan-auditor

Spawn plan-auditor once from the main session with the SPEC directory as its only input — not the conversation or implementation context, so its verdict stays independent:

> Audit the SPEC at `.moai/specs/<SPEC-ID>/` (spec.md, plan.md, acceptance.md, and tasks.md/design.md/research.md when present). Return a PASS/FAIL verdict and write the report to `.moai/reports/plan-audit/<SPEC-ID>-review-<N>.md`.

Route on the verdict in the report:

- **PASS** — the gate is met; continue to the Implementation Kickoff Approval.
- **FAIL** — do not implement. Show the report path and the must-pass failures, then ask through AskUserQuestion: revise the SPEC (recommended; fix the listed defects, directly or through manager-spec, and re-audit only the defect delta), override and proceed (recorded as BYPASSED with the user's reason), or stop.
- **INCONCLUSIVE** — a timeout, malformed report, missing verdict or write failure. INCONCLUSIVE is never treated as PASS. Retry the audit; after three attempts in one run, ask the user whether to proceed with acknowledgement (recorded as `inconclusive_acknowledged_by: <user.name>`) or stop.

Where the SPEC's `tier:` is `S`, one audit pass is enough: a PASS is final without a re-run for score. FAIL and INCONCLUSIVE route as above for every tier.

## --skip-audit

`--skip-audit` (or `MOAI_SKIP_PLAN_AUDIT=1`) skips the audit. Record the bypass with the user's name from `.moai/config/sections/user.yaml` and their reason — ask for it through AskUserQuestion in an interactive session, record `non-interactive` otherwise — and escape Markdown characters in the reason before writing it. The Implementation Kickoff Approval is still required.

## Record

Append to `.moai/specs/<SPEC-ID>/progress.md`:

```yaml
- audit_verdict: PASS            # or BYPASSED, INCONCLUSIVE (acknowledged)
- audit_report: .moai/reports/plan-audit/<SPEC-ID>-review-<N>.md
- audit_at: <ISO-8601 UTC>
```

The independent audit of the implementation itself is not part of run: sync-auditor performs it during `/moai sync`.

# Methodology routing

Read `constitution.development_mode` from `.moai/config/sections/quality.yaml`:

- `tdd` — new behavior is specified by a failing test first (RED-GREEN-REFACTOR). Conventions: the `moai-workflow-tdd` skill.
- `ddd` — for code with little or no test coverage: characterize current behavior with tests, then change it (ANALYZE-PRESERVE-IMPROVE). Conventions: the `moai-workflow-ddd` skill.

Within a single SPEC both can apply: new files and new functions follow TDD, while edits to existing untested code get characterization tests before they change, whatever the configured mode.

## Delta markers (brownfield SPECs)

When `spec.md` marks items with `[EXISTING]`, `[MODIFY]`, `[NEW]` or `[REMOVE]`, handle them in this order so changes land on a safety net:

1. `[EXISTING]` — context only; add characterization tests where the SPEC relies on that behavior, change no code.
2. `[MODIFY]` — characterization tests first, confirm they pass, then modify and keep them passing (or update the ones the SPEC intentionally changes).
3. `[NEW]` — the full TDD cycle.
4. `[REMOVE]` — find every caller and reference first; remove only what nothing still uses.

A SPEC without markers skips this ordering.
