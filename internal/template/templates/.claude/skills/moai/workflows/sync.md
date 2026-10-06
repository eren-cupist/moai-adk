---
description: >
  Third step of Plan-Run-Sync. After /moai run ships a SPEC, verifies the
  result, runs one independent quality audit (sync-auditor), updates
  documentation and the SPEC status to match what shipped, makes the sync
  commit, and delivers it through manager-git when a push or PR is requested
  or configured.
user-invocable: false
metadata:
  version: "4.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-05-17"
  tags: "sync, documentation, pull-request, quality, verification, pr"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["sync", "docs", "pr", "documentation", "pull request", "changelog", "readme"]
  agents: ["manager-docs", "manager-git", "sync-auditor"]
  phases: ["sync"]
---

# Sync Workflow

`/moai sync` closes a SPEC cycle. Its outcome: the code /moai run shipped has been checked, one independent audit has judged it, the documentation and the SPEC status say what was actually built, a single sync commit records that, and — when the user asked for it or the project is configured for it — the branch is pushed and a PR is open.

## Inputs

`$ARGUMENTS` = `[SPEC-ID] [mode] [path] [flags]`.

- SPEC: the given ID; otherwise the SPEC the current branch or worktree implemented (the in-progress or implemented SPEC whose files the branch changed). With no SPEC named and none inferable, ask whether to sync the current branch's changes; if there are no changes, report "Nothing to sync" and stop.
- Modes: `auto` (default — changed files and the docs they affect), `force` (regenerate all documentation and run full-repository checks), `status` (read-only health report: no auto-fix, no generated tests, no doc or git writes; stop after Phase 10), `project` (project-wide doc update plus full-repository checks).
- Flags: `--pr` (deliver through a PR), `--auto-merge` (merge the PR once checks pass, see `sync/delivery.md`), `--merge` (deprecated alias of `--auto-merge`; log a warning), `--skip-mx` (skip @MX validation and note it in the report).
- `--mode` values are ignored. `--mode pipeline` is rejected with `MODE_PIPELINE_ONLY_UTILITY`: pipeline mode is reserved for utility subcommands, and sync is a multi-agent workflow.

Read `.moai/config/sections/git-strategy.yaml`, `language.yaml`, `quality.yaml`, and `harness.yaml` at the start; the SPEC lives in `.moai/specs/<SPEC-ID>/`.

## Who does what

- The orchestrator (this session) runs the checks, asks the user at the approval points, and keeps the single-writer rule: only one write-capable agent works at a time.
- `sync-auditor` — the independent audit (Phase 7). Read-only.
- `manager-docs` — the only writer of documentation, the SPEC frontmatter status, and the `progress.md §E.4` Sync-phase Audit-Ready Signal section (a heading the SPEC era classifier parses); makes the sync commit. Inject `At start, invoke Skill("moai-workflow-project") for the sync-phase documentation cycle.` into its prompt. It never edits SPEC body content; a needed body change comes back as a blocker and goes to `manager-spec`.
- `manager-git` — branches, pushes, PRs, merges (Phase 13).
- `manager-develop` — fixes, when a check fails and the user chooses to fix.

## Phases

| Phase | What | Detail |
|-------|------|--------|
| 1 | Pre-sync check: tree state, lint, format, type, tests (`gate-sync-1`) | `sync/quality-gates.md` |
| 2 | Database and migration changes | `sync/quality-gates.md` |
| 3 | Full test suite passes | `sync/quality-gates.md` |
| 4 | Configuration and environment changes | `sync/quality-gates.md` |
| 5 | Dependency changes | `sync/quality-gates.md` |
| 6 | Backward compatibility | `sync/quality-gates.md` |
| 7 | Quality verification and the independent audit | `sync/quality-gates.md` and below |
| 8 | Security review of security-sensitive changes | `sync/quality-gates.md` |
| 9 | @MX tag validation | `sync/quality-gates.md` |
| 10 | Coverage | `sync/quality-gates.md` |
| 11 | Divergence analysis and documentation scope (`gate-sync-2`) | `sync/doc-execution.md` |
| 12 | Documentation and SPEC status update | `sync/doc-execution.md` |
| 13 | Sync commit and delivery | `sync/delivery.md` |
| 14 | Completion report | `sync/delivery.md` |

Read each sub-file when you reach its phases. Phases 1-10 are independent of the documentation drafting in Phase 11-12, so nothing in the docs step waits on audit scores, and nothing in the audit reads the docs drafts.

## Approval points

The user decides at: a failed pre-sync check (Phase 1), failing tests (Phases 3 and 7), a breaking change (Phase 6), a critical or high security finding (Phase 8), the documentation scope (`gate-sync-2`, Phase 11), a failed local CI mirror (Phase 13), and any push, PR, merge, or issue comment that the flags did not already request (Phase 13). Ask with AskUserQuestion; put the recommended option first.

A FAIL from the audit or a blocking quality gate stops the sync until it is fixed or the user decides otherwise. When sync was entered automatically at the end of `/moai run`, these approval points still apply.

The Stop hook `.claude/hooks/moai/sync-phase-quality-gate.sh` runs a fast structural check (lint, test, coverage delta, dependency-manifest observation); with `MOAI_SYNC_GATE_BLOCKING=1` it blocks on lint or test failure or a coverage drop over 5 points.

<!-- moai:contract-mode-start id="contract-sync-gates" -->
Where `workflow.autonomy.mode: contract` — the sync phase carries the last three lifecycle stages: Closure (status transition and verdict), Integration (re-measure on the merged tree), Push (inactive until the stop-before-push on a missing second review exists). No documentation-scope approval (`gate-sync-2`), next-step, or current-branch question is asked, and the failure decision points — Phase 1 gate failure, Phase 3 and Phase 7 test failure, Phase 6 compatibility break, Phase 8 critical security finding, Phase 13 local CI mirror failure — are routed to an escalation report instead of a question. The CI auto-fix loop after a pull request is unchanged. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Gate disposition".

<!-- moai:contract-mode-end -->
## The independent audit (Phase 7)

Spawn `sync-auditor` once when the harness level has `evaluator: true` (`harness.yaml` `levels.<level>.evaluator`; standard and thorough). At `minimal` (`evaluator: false`), the Phase 7 verification results stand in for it — unless the audit plan below requires a cross-model backend, in which case spawn it anyway.

Brief the auditor with: the SPEC ID and directory, the audited commit (full `HEAD` SHA) and the base, the acceptance-criteria source (`spec.md` for tier S, `acceptance.md` for tier M/L), the evaluator profile, the Phase 1-10 results it may cite, and — on a re-audit — the defect list it is checking. It returns findings (severity, confidence, blocking or optional, file:line), a verdict, and an `AUDIT-VERDICT:` last line.

The auditor cannot write files, so save its report yourself to `.moai/reports/<SPEC-ID>/sync-audit.md`, noting that you relayed it. On FAIL, send the blocking findings to `manager-develop`, then re-audit scoped to those findings. Optional findings are reported, not auto-fixed. After three fix rounds without a PASS, stop and ask the user.

### Cross-model audit plan

The tree's audit plan decides whether a non-Claude backend must also pass. Learn your own toplevel with `git rev-parse --show-toplevel` (in a worktree session `CLAUDE_PROJECT_DIR` names the primary checkout) and run `moai verify audit-plan --project-root <toplevel>`. It is read-only and prints one JSON object:

- A plan — `config_status` is `ok` or `absent`. With `cross_model_required: false`, the sync-auditor verdict alone is binding. With `cross_model_required: true`, the auditor calls `audit_multi` without a `gates` argument and returns the result's digest members.
- The `verify` group's help text — output containing `Shared diagnostic snapshot contract` and neither a `config_status` member nor an `audit-plan:` line — means the binary predates the verb. That is the legacy path: continue with the auditor's own legacy handling and record "plan surface unreachable, legacy path used" as a named Gap.
- Anything else — `config_status: unreadable`, an `audit-plan:` error line, a refused, crashed, or malformed run — is not the legacy path. It is a PASS-blocking Gap, and a configured required backend that does not answer stays fail-closed by name.

When the auditor called `audit_multi`, its verdict is not final until the orchestrator's check passes. Write the digest it returned (`overall_verdict`, `gate_unmet`, `plan_source`, and per `per_backend_verdicts` entry `backend`, `gate`, `verdict` — never summary or finding text) with your own Write tool to `<toplevel>/.moai/state/audit-plan-result.json`, fresh — overwrite any earlier file, and never reuse one from an earlier audit. Then run `moai verify audit-plan --result-file <that path> --project-root <toplevel>`, passing only the path — never the JSON on the command line, which the worktree guard refuses.

The sync verdict is binding only when both must pass: the sync-auditor PASS, and the cross-model half — `audit_multi` verdict `pass`, empty `gate_unmet`, and `convergence_check.ok: true`. When `audit_multi` cannot be called while the plan requires it, the reason is `audit_multi unreachable`. State the outcome with these lines: `cross_model: required (<backends>)`, `audit_multi: pass|fail|unavailable`, `gate_unmet: <backends|none>`, `audit_receipt: <ids|none>`, `plan_check: ok|unmet (<backends>)`; when the verdict is not binding, `binding: no — cross-model gate unmet: <backend>` (or `audit_multi unreachable`), never a bare PASS.

## Done

Sync is done when: Phases 1-10 have results (or recorded skips with reasons); the audit verdict is PASS or PASS-WITH-DEBT and binding; the documentation the divergence report named is updated; the SPEC status and the Sync-phase Audit-Ready Signal in `progress.md` reflect the outcome; the sync commit exists; any requested delivery happened; and the completion report is shown. In `status` mode, done is the read-only report after Phase 10.
