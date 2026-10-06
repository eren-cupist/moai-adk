---
description: >
  Run phase of Plan-Run-Sync: implements an approved SPEC so its acceptance
  criteria are met, the project's tests pass, and progress.md records the
  evidence. Methodology follows quality.yaml constitution.development_mode.
user-invocable: false
metadata:
  version: "2.6.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-23"
  tags: "run, implementation, ddd, tdd, spec"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["run", "implement", "build", "create", "develop", "code"]
  agents: ["manager-develop", "manager-git"]
  phases: ["run"]
---

# Run Workflow

`/moai run SPEC-XXX` turns an approved SPEC into working code. The run is done when every acceptance criterion in scope is met in the codebase, the project's tests, typecheck and lint pass, and `progress.md` records the evidence. `/moai sync` follows and owns documentation, the independent sync audit and the PR.

## Inputs

- `$ARGUMENTS`: the SPEC-ID (e.g. `SPEC-AUTH-001`), plus optional flags: `--skip-audit` (bypass the Plan Audit Gate, recorded as BYPASSED), `--pr` (take the PR route regardless of tier), `--resume` (continue from what `progress.md` shows as done; re-running `/moai run` on a SPEC with recorded run evidence resumes the same way).
- `.moai/specs/SPEC-{ID}/`: `spec.md` (requirements; frontmatter `tier:` and `status:`), `plan.md`, `acceptance.md` (Tier M/L), and when present `tasks.md`, `research.md`, `design.md`, `spec-compact.md` (a shorter spec.md; read it instead when present), `progress.md` (scaffolded by plan, filled by run).
- `.moai/config/sections/quality.yaml`: `constitution.development_mode` (`tdd` or `ddd`), `test_coverage_target`, `lsp_quality_gates.run` (max errors / type errors / lint errors, `allow_regression`), `memory_guard` (when enabled, run tests in module-sized batches).
- `.moai/config/sections/git-strategy.yaml` (mode, `automation.*`, branch prefix) and `language.yaml` (`git_commit_messages`, `code_comments`).
- `.moai/project/structure.md` and `tech.md` when present, for architecture conventions.

If the SPEC directory does not exist, say so and offer `/moai plan` through AskUserQuestion.

## Flow

1. **Load context.** Read the SPEC files and config above, and `git status` / current branch. If `progress.md` already carries run evidence, resume from the first unmet criterion instead of starting over.
2. **Plan Audit Gate.** Every run passes it before code is written — procedure in `workflows/run/phase-execution.md`.
3. **Implementation Kickoff Approval** (below).
4. **Decide who implements** and implement with the configured methodology — `workflows/run/task-decomposition.md`, with the decision rule in `.claude/rules/moai/workflow/orchestration-mode-selection.md`.
5. **Verify in the main loop**, record `progress.md`, commit per the git strategy, and offer the next step — also in `workflows/run/task-decomposition.md`.

| Sub-file | Read when |
|---|---|
| `workflows/run/phase-execution.md` | Plan Audit Gate, methodology routing, `[MODIFY]`/`[NEW]` delta markers |
| `workflows/run/task-decomposition.md` | implementing, briefing manager-develop, verification, progress.md, git, completion |
| `workflows/run/external-delegation.md` | optional hand-off of a bounded mechanical subtask to an external model |

## Implementation Kickoff Approval

Code for a SPEC is written only after the user has approved starting implementation — the Implementation Kickoff Approval, the one human gate between plan and run. If this conversation already holds the user's yes to the plan phase's Kickoff question for this SPEC and the plan artifacts have not changed since, do not ask again. Otherwise, once the Plan Audit Gate has passed, ask a single AskUserQuestion round: start implementation (recommended), revise the plan, or stop. Put any other preference the run needs from the user (for example the PR route when the tier leaves it open) in the same round, so nothing has to be asked mid-run.

The approval is required regardless of the plan-auditor score: a high score only lets the gate reuse a prior PASS instead of re-auditing; it never stands in for the user's approval. The question is always the orchestrator's — subagents cannot ask the user, so it is never delegated.

<!-- moai:contract-mode-start id="contract-signing-run" -->
Where `workflow.autonomy.mode: contract` — the gate above is `moai contract kickoff-check <SPEC-ID> --card <card>` exiting 0, and no Kickoff question is emitted. A kickoff decision outcome of reject or human refuses the receipt signature and needs a human decision, so follow the human signing procedure for both. A non-zero exit is escalated, never downgraded to guided mode. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "The signing gate".

<!-- moai:contract-mode-end -->
<!-- moai:contract-mode-start id="contract-lifecycle-run" -->
Where `workflow.autonomy.mode: contract` — the run phase carries the first four stages of the one-pass lifecycle: Discovery (re-observe the contract's `reobserve` list), RED (commit the failing tests first), GREEN (the same tests pass), Qualification (scoped tests, lint, coverage, mutation check, second-model review). Record each stage's evidence and run `moai contract kickoff-check` at every stage boundary; an open escalation record or a revocation stops the run there. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "One-pass lifecycle".

<!-- moai:contract-mode-end -->
## Scope and judgment

Deliver what the SPEC asks, at the scope it states. Make routine judgment calls yourself and check in only when different readings would lead to materially different work. If the SPEC looks mistaken, say so in a sentence and continue as written. When implementation shows that a requirement or acceptance criterion itself has to change, stop and route that change back through manager-spec rather than editing the SPEC body during run. Finish the whole SPEC; if something genuinely can't be completed, do the rest and state plainly what is missing and why.

## Bounded fix loops

A mechanical failure (lint, format, type, build, import) is fixed and re-checked; after three attempts on the same failure without progress, stop and ask through AskUserQuestion whether to continue manually, revert and re-plan, or stop. A failing test, panic, data race or deadlock means the code is wrong: fix the code, not the test, and never weaken, skip or delete an existing test to get a green run unless the SPEC requires that change. Never edit `.env*`, credentials, or CI workflow definitions as part of a fix.

## Communication

Your text between tool calls is what the user reads. Before the first tool call say in a sentence what you're about to do; give brief updates when you find something load-bearing or change direction. Lead the final report with the outcome, in complete sentences, then the detail a reader needs to act. Keep reports to the length the work needs.
