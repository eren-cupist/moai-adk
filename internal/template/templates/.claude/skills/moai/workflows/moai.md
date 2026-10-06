---
description: >
  Default /moai route for a development request without a subcommand: plan the
  work as a SPEC, get the Implementation Kickoff Approval, implement it with
  /moai run, and finish with /moai sync.
user-invocable: false
metadata:
  version: "4.0.0"
  category: "workflow"
  status: "active"
  tags: "moai, pipeline, plan-run-sync, default"
---

# /moai (default route)

## Outcome

The requested change implemented, verified against its SPEC, and documented: plan → Implementation Kickoff Approval → run → sync, each phase following its own workflow file. The user approves once, at the plan→run boundary; after that the pipeline continues on its own unless a gate fails.

Deliver what was asked, at the intended scope. Make routine judgment calls yourself and check in only when different readings would lead to materially different work. If the ask looks mistaken, say so in a sentence and continue as asked. Finish the whole task; if something genuinely can't be completed, do the rest and state plainly what is missing and why.

A change too small to be worth a SPEC — a typo, a one-line fix with an obvious test — does not need the pipeline. Make it directly, verify it, and say that you skipped the SPEC and why.

## Flags

`--branch` and `--issue` pass to plan; `--pr` puts the SPEC on the PR route (`.claude/rules/moai/workflow/spec-workflow.md`, section "SPEC Phase Discipline"); `--resume SPEC-<ID>` continues an existing SPEC from its first unfinished phase; `--team` (experimental Agent Teams) and `--solo` pass to run.

## Pipeline

1. **Plan** — follow `workflows/plan.md` to its end: SPEC written, plan-auditor PASS, audit-ready signal in `progress.md` §E.1.
2. **Implementation Kickoff Approval** — the plan workflow's AskUserQuestion: start implementation now (recommended), revise the SPEC, or stop. It is asked whatever the audit score, and nothing below starts without it. If where the work runs matters, ask it in the same AskUserQuestion call as a second question: this session (default), or an isolated workspace the user opens with `moai cc -w SPEC-<ID>` and continues in.
3. **Run** — follow `workflows/run.md` for `SPEC-<ID>`. Its Plan Audit Gate reuses the PASS from step 1 when the plan artifacts are unchanged.
4. **Sync** — when run finishes with every acceptance criterion met, continue into `workflows/sync.md` without another approval round; say in the transcript that you are moving on. Sync's own gates and its approval before outward-facing actions (push, PR) still apply.

The pipeline stops and reports when a gate fails: a plan audit that will not pass, a run that cannot meet a criterion or keeps failing the same way, or a sync audit FAIL. Bring the failure to the user with the evidence and the options; never continue past a failing gate.

<!-- moai:contract-mode-start id="contract-pipeline-gates" -->
Where `workflow.autonomy.mode: contract` — gate 2 is the contract signature: `moai contract kickoff-check <SPEC-ID> --card <card>` must exit 0 and no Kickoff question is emitted; a plan-audit FAIL is repaired automatically up to `budget.audit_retries` and then escalated. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "The signing gate".

<!-- moai:contract-mode-end -->
<!-- moai:contract-mode-start id="contract-merged-round" -->
Where `workflow.autonomy.mode: contract` — the step 2 round carries no Kickoff question: it asks only the execution-shape question, and run-phase entry waits for `moai contract kickoff-check` to exit 0. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "The signing gate".

<!-- moai:contract-mode-end -->

## Finishing

When sync completes, report the outcome: what changed, how it was verified, the SPEC ID, and anything left open. If there is no real decision left for the user, end there rather than inventing a next-step question.

Your text between tool calls is what the user reads. Before the first tool call say in a sentence what you're about to do; give brief updates when you find something load-bearing or change direction. Lead the final report with the outcome, in complete sentences, then the detail a reader needs to act. Keep reports to the length the work needs.
