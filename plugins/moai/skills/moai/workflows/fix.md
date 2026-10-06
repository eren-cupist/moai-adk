---
description: >
  Fix workflow for a reported bug, a failing test or CI check, or typecheck
  and lint errors: reproduce, fix minimally, verify, with bounded retries.
  Not for new features (use plan and run).
user-invocable: false
metadata:
  version: "2.5.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-21"
  tags: "fix, bug, ci, lint, type-check, diagnostics"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["fix", "bug", "error", "lint", "type error", "failing test", "ci failure"]
  agents: ["manager-develop"]
  phases: ["fix"]
---

# Fix Workflow

`/moai fix` repairs something that is broken — a reported bug, a failing test or CI check, or typecheck and lint errors — with the smallest change that addresses the root cause, and shows that it worked.

## Input

`$ARGUMENTS` is a description of the problem, an error message or log, a failing test or check name, a path to examine, or nothing (then the target is the errors reported by the project's typecheck and lint). `--dry` reports the diagnosis and the proposed change without editing anything.

## Steps

1. **Reproduce.** Run the failing test or check, or the project's typecheck, lint and tests for the target scope, and keep that output as the baseline. Where a test harness exists, capture the bug as a failing test first; it becomes the proof of the fix. If the problem cannot be reproduced, say so and report what you tried instead of guessing at a fix.
2. **Find the root cause.** Read the code involved, including its `@MX:ANCHOR` and `@MX:WARN` notes (`.claude/rules/moai/workflow/mx-tag-protocol.md`): an anchored function has callers that depend on its contract.
3. **Fix minimally.** Change what the root cause requires and nothing around it. Formatting and import-order errors are fixed by running the project's formatter. Before a change that alters a public API, behavior beyond the bug, or security-relevant code, explain it and get the user's confirmation through AskUserQuestion.
4. **Verify.** The reproduction test passes. Re-run the same commands as step 1 over the same scope and compare with the baseline: the targeted failures are gone and nothing new fails. A regression the fix introduced is reverted or reported by file and line — never left in silently.
5. **Report.** The root cause, the change, the commands with their observed results, and anything still open.

Up to three fix attempts per failure. When the third attempt does not resolve it, stop and ask through AskUserQuestion: continue investigating manually, revert the attempted changes, or stop. Never edit `.env*`, credentials or CI workflow definitions to make a check pass, and never weaken or delete a test to turn it green.

Leave the changes uncommitted unless the user asked for a commit; a CI auto-fix commits as its protocol describes.

## Who does the work

Fix it directly in the main session; most fixes are a few reads and edits. Hand a fix to manager-develop only when it is large enough to deserve its own worktree, briefed per `.claude/rules/moai/development/manager-develop-prompt-template.md`. A failing required CI check on a pull request goes to manager-develop with `cycle_type=autofix`, under `.claude/rules/moai/workflow/ci-autofix-protocol.md`: at most three iterations, each patch a new commit, no force-push, semantic failures (test assertions, panics, races, deadlocks) escalated to the user instead of patched.

A bounded mechanical subtask may use the optional external-model delegation in `.claude/skills/moai/workflows/run/external-delegation.md`, section "External Model Delegation".
