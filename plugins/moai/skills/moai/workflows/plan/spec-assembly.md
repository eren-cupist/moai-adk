---
description: "/moai plan reference — SPEC artifact checklist, plan-audit mechanics, Implementation Kickoff Approval, and the optional branch and issue steps"
user-invocable: false
metadata:
  parent: moai-workflow-plan
  phase: "SPEC assembly, audit, and kickoff"
---

# Plan reference: assembly, audit, kickoff

Reference for `workflows/plan.md`. It holds the facts the plan workflow needs — file names, fields, report formats, config keys, commands — not a procedure.

## SPEC artifacts

manager-spec writes all of these; nobody else writes into `.moai/specs/SPEC-<ID>/` while it works.

| Tier | Files |
|------|-------|
| S | `spec.md` (with acceptance criteria inline), `plan.md`, `progress.md` |
| M | `spec.md`, `plan.md`, `acceptance.md`, `progress.md` |
| L | `spec.md`, `plan.md`, `acceptance.md`, `design.md`, `research.md`, `progress.md` |

Plus `decision-index.md` when `interview.decision_gate` is `on` in `.moai/config/sections/interview.yaml` (default `off`).

Before handing the SPEC to the auditor, check that:

- the SPEC ID matches `^SPEC(-[A-Z][A-Z0-9]*)+-[0-9]{3}$` and is a directory, not a flat file;
- `spec.md` frontmatter has the 12 fields of `.claude/rules/moai/development/spec-frontmatter-schema.md`, section "Canonical 12 Required Fields" — `id`, `title`, `version` (quoted semver), `status: draft`, `created`, `updated`, `author`, `priority`, `phase` (a release target, never `plan`/`run`/`sync`/`mx`), `module`, `lifecycle`, `tags` (comma-separated string) — plus `tier:`, with no `created_at`/`updated_at`/`labels`/`spec_id` aliases;
- only `spec.md` has a `status:` field;
- the exclusions section has a `### Out of Scope — <topic>` heading with bullets;
- every requirement is mapped by an acceptance criterion;
- `progress.md` has the four `§E.1`–`§E.4` headings, filled only at §E.1;
- `moai spec lint SPEC-<ID>` reports no errors.

`issue_number` is added only when a GitHub issue is created.

## Plan audit

Invoke plan-auditor with the SPEC directory path, the user's original request verbatim, and the round number `<N>` (1 for the first audit), plus the previous report path on later rounds. Pass nothing else — no reasoning, no drafts — so the audit stays independent.

The report lands at `.moai/reports/plan-audit/<SPEC-ID>-review-<N>.md`. Its header lines are machine-read by the run gate (`runtime.ResolveLatestPlanAudit`), the verdict admission code (`internal/auditverdict`), and `moai plan render-html`:

```
Verdict: PASS | FAIL
Overall Score: <0.00-1.00>
must_pass_failed: <count>
blocking_count: <count>
```

PASS requires zero must-pass failures, zero blocking findings, and a score at or above the tier threshold (S 0.75, M 0.80, L 0.85). Findings are listed as `D<n>.` lines with severity, confidence, and blocking/optional class. An audit that times out, errors, or returns a report without a parseable verdict is INCONCLUSIVE — treat it as not passed and run it again or ask the user.

Re-audit rounds per tier are capped by `harness.plan_audit_tier_ceilings` in `.moai/config/sections/harness.yaml` (shipped as S 1, M 2, L 3). A round counts each time plan-auditor is invoked.

<!-- moai:contract-mode-start id="contract-audit-retry" -->
Where `workflow.autonomy.mode: contract` — do not present the three options; repair and re-audit automatically up to the contract's `budget.audit_retries` (or `workflow.autonomy.escalation.budget_default.audit_retries` when the draft has no budget), then stop with an escalation report. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Gate disposition".

<!-- moai:contract-mode-end -->
<!-- moai:contract-mode-start id="contract-quality-gate" -->
Where `workflow.autonomy.mode: contract` — a WARNING or FAIL is repaired automatically up to the same retry cap (`budget.audit_retries`, or `workflow.autonomy.escalation.budget_default.audit_retries`), then reported as an escalation instead of a question. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Gate disposition".

<!-- moai:contract-mode-end -->

<!-- moai:contract-mode-start id="contract-draft" -->
Where `workflow.autonomy.mode: contract` — the manager-spec delegation also asks for a draft `contract.yaml` (no `signature` block) beside the SPEC artifacts, carrying the card, acceptance binding, invariants, ownership, approach, actions, reobserve list, review, and escalation triggers. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Scope and activation".

<!-- moai:contract-mode-end -->

## Implementation Kickoff Approval

The plan→run gate. After the audit passes, the audit-ready signal is written, and no `[NEEDS CLARIFICATION` marker remains, ask the user with AskUserQuestion whether to start implementation. Three options, recommended first:

- start `/moai run SPEC-<ID>` now (recommended);
- revise the SPEC first;
- stop here and keep the SPEC as a draft.

The approval is asked regardless of the plan-auditor score; a passing audit is evidence for the user's decision, never a substitute for it. While `interview.recommendation_mode` is `pull`, the "recommended" label is withheld and no option claims a preference.

Optionally, before asking, run `moai plan render-html SPEC-<ID>`. It reads the SPEC directory and the latest review report and writes a self-contained page to `.moai/reports/plan-html/<SPEC-ID>-plan.html`; give the user the path alongside the question. Fail-open: if the command is missing or fails, skip it — the question does not depend on it.

When `decision-index.md` exists, list its rows with the question — what is unresolved and why, without a preferred answer. Afterwards write each answer into the row's `Operator verdict:` line using `DECIDE`, `NEED_ANALYSIS`, `NEED_EVIDENCE`, or `DEFER`. An index with zero rows is not an approval; the question is still asked.

<!-- moai:contract-mode-start id="contract-signing-review" -->
Where `workflow.autonomy.mode: contract` — the decision index is presented with the signing summary instead of the Kickoff question: report `moai contract sign <SPEC-ID>` for the operator to run at an interactive terminal, close the turn, and run `moai contract kickoff-check` on the next turn. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "The signing gate".

<!-- moai:contract-mode-end -->

## Optional: branch

Only on `--branch`, or when `git_strategy.<mode>.branch_creation.auto_enabled` is `true` in `.moai/config/sections/git-strategy.yaml` (the shipped default is `false`; `prompt_always: true` means ask before creating one). Otherwise the SPEC stays on the current branch.

1. Choose the base with the Branch Origin Decision Protocol (`.claude/rules/moai/development/branch-origin-protocol.md`): run its three signal checks, recommend a base, and confirm it with the user.
2. Have manager-git create `feature/SPEC-<ID>-<short-name>` from that base (`base=<branch>`), set the upstream if a remote exists, and switch to it.

Plan does not create worktrees. To plan or run in isolation, the user enters a workspace first with `moai cc -w <name>` (or `moai cc -w <name> --spawn` for a separate session window) and runs the command there.

## Optional: GitHub issue

Only on `--issue`, with `gh` available and a remote configured; otherwise skip without comment. manager-git runs:

```bash
gh issue create \
  --title "[SPEC-<ID>] <title>" \
  --label "spec" \
  --body "<one-paragraph summary of the requirements and acceptance criteria; SPEC location .moai/specs/SPEC-<ID>/spec.md>"
```

Then record the number as `issue_number: <N>` in spec.md frontmatter; run and sync use it for `Fixes #<N>` in commits and the PR.
