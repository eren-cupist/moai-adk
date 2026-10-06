---
description: CI auto-fix loop protocol — limits for the manager-develop autofix cycle (cycle_type=autofix). Loaded when working on the autofix cycle or on CI workflow definitions.
paths: ".claude/agents/moai/manager-develop.md,.claude/rules/moai/development/manager-develop-prompt-template.md,.github/workflows/**"
---

# CI Auto-Fix Protocol

The CI auto-fix loop repairs a pull request whose required check is failing, by handing the failure to manager-develop with `cycle_type=autofix`. Its limits exist because an unattended loop on CI can otherwise hide a real failure, rewrite shared history, or leak secrets.

## Entry Condition

[ZONE:Frozen] Enter the auto-fix loop only when the orchestrator hands off a failing required check: the orchestrator observes the failure, then briefs manager-develop with the pull request and branch, the failing check and its log, and the iteration number. No failing required check, no loop.

## Iteration Cap

[ZONE:Frozen] Attempt at most three iterations per pull request; the counter lives in `.moai/state/ci-autofix-<PR>.json` (created at entry, updated after each iteration, deleted once the checks pass; a state file older than 24 hours may be reclaimed).

[ZONE:Frozen] After the third failed iteration, ask the user through a blocking AskUserQuestion with no timeout: fix it manually (recommended), revise the SPEC and restart implementation, or close the pull request. Nothing further happens until the user answers.

## Patch Commit Rule — No Force-Push

[ZONE:Frozen] Apply every auto-fix patch as a new commit on the PR branch, staging only the patched files, with a message like `fix(ci): auto-fix <classification> failure (iter <N>)`. Never `git push --force`, `-f`, `--force-with-lease` or `git commit --amend`: the branch is shared with reviewers and CI history. After pushing, wait for the same required check to re-run before deciding on another iteration.

## Semantic Failure — No Auto-Patch

Classify the failure from the check's own output. Lint, formatting, build, type and missing-dependency failures are mechanical and may be patched. Data races, deadlocks, panics and test assertion failures are semantic, and so is anything that cannot be classified.

[ZONE:Frozen] For a semantic failure, do not patch: escalate immediately through AskUserQuestion with the diagnosis — what failed, the likely root cause, and the options. A read-only investigation may produce the diagnosis; it proposes no patch.

For mechanical failures, confirm the first patch with the user; trivial mechanical patches in iterations two and three may be applied without asking.

## Secrets and Credentials Protection

[ZONE:Frozen] The auto-fix loop never modifies `.env`, `.env.*`, credentials files, API key files, `**/*secret*`, `**/*_key.json`, `.claude/settings.json` or `.claude/settings.local.json`. A patch that would touch one is rejected and escalated to the user.

## CI Infrastructure Preservation

[ZONE:Frozen] The auto-fix loop never modifies CI watch infrastructure scripts or workflow definitions, required-check configuration, or any script that observes or reports CI status: a patch to the reporting layer can turn a real failure into a false green.

## Audit Log

Append one entry per iteration to `.moai/logs/ci-autofix/<PR-NNN>-<YYYY-MM-DD>.md` (local, gitignored): iteration number, classification, action (applied / escalated / aborted), the patch commit SHA when applied, and the escalation reason when escalated.
