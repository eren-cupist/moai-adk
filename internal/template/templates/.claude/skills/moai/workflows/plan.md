---
description: >
  /moai plan — turn a feature or change request into an approved SPEC
  directory (spec.md, plan.md, acceptance.md and the tier's other artifacts)
  that /moai run can execute, with an independent plan-auditor review and the
  Implementation Kickoff Approval.
user-invocable: false
metadata:
  version: "3.0.0"
  category: "workflow"
  status: "active"
  tags: "plan, spec, gears, requirements, acceptance-criteria"
---

# /moai plan

## Outcome

A SPEC directory at `.moai/specs/SPEC-<ID>/` that the user can read and approve and that `/moai run SPEC-<ID>` can execute: GEARS requirements, acceptance criteria covering every requirement, the artifact set for the SPEC's tier, a passing plan-auditor report, the audit-ready signal in `progress.md` §E.1, and the user's Implementation Kickoff Approval. Planning writes no implementation code.

Deliver what was asked, at the intended scope. Make routine judgment calls yourself and check in only when different readings would lead to materially different work. If the ask looks mistaken, say so in a sentence and continue as asked. Finish the whole task; if something genuinely can't be completed, do the rest and state plainly what is missing and why.

## Arguments

`/moai plan "<description>" [--branch] [--issue] [--resume SPEC-<ID>]`

- `--resume SPEC-<ID>` continues an existing draft: read its directory and pick up at the first missing or failing piece.
- `--branch` creates a branch for the SPEC (see Optional steps). `--issue` creates a GitHub issue. Without these flags plan creates neither.
- No description and no `--resume`: ask the user what to plan.
- `--mode pipeline` is rejected with `MODE_PIPELINE_ONLY_UTILITY`; plan is a multi-agent workflow.

## Reference

- SPEC format (GEARS, requirement and criterion lines, spec.md and plan.md contents): the `moai-workflow-spec` skill.
- Tiers and their artifact sets, the plan audit gate, report paths: `.claude/rules/moai/workflow/spec-workflow.md`.
- Frontmatter and `progress.md` sections: `.claude/rules/moai/development/spec-frontmatter-schema.md`.
- Artifact checklist, audit mechanics, and the optional branch and issue steps: `workflows/plan/spec-assembly.md`.

## How to get there

**Understand the request and the code it touches.** Explore as much as the request needs. For a change confined to a few files, read them yourself. For wide work — several modules, unfamiliar territory, or a question about how things fit together — launch one or two Explore subagents in one message so they run in parallel, each with a distinct question; don't send a subagent to read what you could read in a few calls. Check `.moai/specs/` for an existing SPEC covering the same ground (extend or supersede it rather than duplicating it) and for SPECs this one depends on. Read `.moai/project/product.md`, `structure.md`, and `tech.md` when they exist.

Subagents multiply cost and time: each re-establishes context, re-explores and reports back. Do small work (a few reads, a handful of edits, simple verification) directly. Delegate large, genuinely independent tracks, and run independent agents in one message so they run in parallel. Brief a subagent fully the first time, and don't redo its work once it reports. Verification belongs in the main loop, not in an extra subagent.

**Ask only what changes the SPEC.** When different readings of the request would produce materially different SPECs — what is in or out of scope, which behavior is wanted, a constraint only the user knows — ask, all in one AskUserQuestion round (up to 4 questions, recommended option first, each option with a short description of what it implies). Everything else is a routine call: decide it and have manager-spec record it under Assumptions in plan.md. A question that cannot be settled now becomes a `[NEEDS CLARIFICATION: <topic>]` marker in plan.md (or research.md) — never in spec.md or acceptance.md — and must be answered before the Kickoff Approval.

**Pick the tier and the ID.** Choose S, M, or L from the scope you found, using the table in `spec-workflow.md`; a user-stated tier wins. Ask only when the evidence genuinely sits between two tiers and the difference matters. Choose a SPEC ID `SPEC-<DOMAIN>-<NNN>` that does not exist yet.

**Have manager-spec write the SPEC.** Brief it with the request in the user's words, what exploration found (files, patterns, constraints), the user's answers, the assumptions you made, the tier, and the ID. It writes the tier's artifact set plus `progress.md`, runs `moai spec lint`, and reports back. Where `interview.decision_gate` is `on`, it also writes `decision-index.md`.

**Get an independent audit.** Invoke plan-auditor with only the SPEC directory, the user's original request verbatim, and the round number — not your reasoning or the conversation, so the review is independent. It writes `.moai/reports/plan-audit/<SPEC-ID>-review-<N>.md` and returns the verdict. Fix every blocking finding (send manager-spec the report path and the findings to fix; a one-line mechanical fix you may make yourself). Apply optional findings that are clearly right, and mention the rest to the user. Re-audit only when the fixes were substantive — requirements, criteria, or scope changed — and stay within the round ceiling for the tier in `harness.plan_audit_tier_ceilings` (`.moai/config/sections/harness.yaml`). If the SPEC still does not pass at the ceiling, or the auditor reports the score dropping between rounds, stop iterating and ask the user whether to revise the request, accept the SPEC with the listed defects, or stop.

**Resolve open questions.** If any `[NEEDS CLARIFICATION` marker remains, ask the user (one round), have manager-spec fold the answers in, and re-audit if the answers changed requirements or criteria.

<!-- moai:contract-mode-start id="contract-clarification" -->
Where `workflow.autonomy.mode: contract` — markers still block the run: the contract is signed only after they are resolved, and `moai contract decide` treats any marker left in plan.md or research.md as a failed precondition that routes the Kickoff to a human signature. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "The signing gate".

<!-- moai:contract-mode-end -->

**Record the audit-ready signal.** Once the audit passes and no marker remains, write under `## §E.1 Plan-phase Audit-Ready Signal` in `progress.md`:

```
- plan_complete_at: <ISO-8601 timestamp>
- plan_status: audit-ready
```

`/moai run` reads this signal before its Phase 1 Plan Audit Gate (`progress.md §E.1`; run and sync own §E.2–§E.4).

**Ask for the Implementation Kickoff Approval.** Summarize the SPEC for the user — scope, the requirements in a sentence each, the assumptions made, the audit verdict and anything you chose not to fix — with the paths to the SPEC directory and the audit report. Then ask with AskUserQuestion: start `/moai run SPEC-<ID>` now (recommended), revise the SPEC first, or stop here. This approval is required whatever the audit score; nothing in plan starts implementation without it. Details are in `spec-assembly.md`.

## Optional steps

These happen only when the user asks or the project config turns them on; the facts the code needs are in `spec-assembly.md`.

- **Branch** — on `--branch`, or when `git_strategy.<mode>.branch_creation.auto_enabled` is `true` in `.moai/config/sections/git-strategy.yaml`: pick the base with the Branch Origin Decision Protocol (`.claude/rules/moai/development/branch-origin-protocol.md`) and have manager-git create it.
- **GitHub issue** — only on `--issue`, through manager-git, then record `issue_number` in spec.md.
- **Commit** — plan does not commit or push on its own. When the user asks, commit the SPEC directory with the subject `feat(SPEC-<ID>): plan-phase artifacts (<tier>, <N> artifacts)` and an `Authored-By-Agent: manager-spec` trailer; never commit straight to a protected or shared branch.

## Done

The SPEC directory holds the tier's artifacts and `progress.md`; `moai spec lint SPEC-<ID>` shows no errors; the latest plan-audit report says `Verdict: PASS`; no `[NEEDS CLARIFICATION` marker remains; `progress.md` §E.1 carries `plan_complete_at` and `plan_status: audit-ready`; and the user answered the Implementation Kickoff Approval.

Your text between tool calls is what the user reads. Before the first tool call say in a sentence what you're about to do; give brief updates when you find something load-bearing or change direction. Lead the final report with the outcome, in complete sentences, then the detail a reader needs to act. Keep reports to the length the work needs.
