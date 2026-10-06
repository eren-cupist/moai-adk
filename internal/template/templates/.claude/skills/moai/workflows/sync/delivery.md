---
description: "Sync Phases 13-14 — the sync commit, local CI mirror, strategy-aware push and PR delivery, auto-merge, the completion report, and edge cases."
user-invocable: false
metadata:
  parent: moai-workflow-sync
  phase: "Phase 13-14: Commit, Delivery, and Completion"
---

# Sync Phases 13-14: Delivery

## Phase 13: Commit and deliver

### Resolve the strategy

From `.moai/config/sections/git-strategy.yaml`: read `git_strategy.mode`, then `git_strategy.{mode}.workflow`. That value alone selects the delivery strategy, and its domain is exactly `github-flow` or `git-flow`. Any other value, or a missing `workflow` key under the active mode, stops delivery: report the value and the domain, and push nothing. Also resolve `main_branch` (`git_strategy.{mode}.main_branch`, default `main`) and `merge_method` (`git_strategy.{mode}.merge_method`: `squash`, `merge`, or `rebase`; default `squash`). Whether a SPEC gets its own branch is a separate setting, `git_strategy.{mode}.automation.auto_branch`, which `manager-git` applies.

### The sync commit

One commit carries all sync output: documentation, reports, the SPEC frontmatter status, and the `progress.md` signal. It is made on the branch the run worked on — never on a branch listed in `workflow.branch_guard.deny_commits_on` or on git-flow's main branch. If the current branch is one of those, stop and ask the user how to proceed (usually: create a branch or switch to the SPEC's worktree).

- Tier S or M without `--pr`: `manager-docs` makes the commit; no PR.
- Tier L, or `--pr`: `manager-git` makes the same commit on the SPEC's feature branch and delivers it through a PR.

Subject: `docs(SPEC-{ID}): sync-phase artifacts + 3-phase close` when the SPEC reaches `completed` (the drift detector reads the `3-phase close` infix as the completed transition, and the subject must name exactly one full SPEC ID); `docs(SPEC-{ID}): sync-phase artifacts` when the SPEC stays `in-progress`. The message language follows `language.git_commit_messages`. The body lists the documents updated and the SPEC status, and carries a short context section:

```
## Context (AI-Developer Memory)
- Decision: <documentation or divergence decision and why>
- Constraint: <constraint discovered>
- Gotcha: <issue found and how it was resolved>
```

After the commit, backfill the real SHA into `sync_commit_sha` in a follow-up commit (a commit cannot contain its own hash).

### Pushing needs approval

`--pr` and `--auto-merge` are the user's request for that delivery. Without them, ask before any push, PR, or issue comment — once, covering everything this phase would send — unless `git_strategy.{mode}.automation.auto_push` is true and the push targets the SPEC's own branch. Brief `manager-git` with what was approved; it does not push or open PRs otherwise.

### Local CI mirror (before a PR)

Before pushing a branch that will get a PR, run the CI checks locally so failures show up in minutes instead of after the push. Skip this for a direct push and for a git-flow `WT-*` integration merge.

Read `.github/workflows/` (or the project's CI config); if there is none, log "No CI config detected" and continue. Map each job to its local equivalent — tests, lint, type check, build, and the build matrix's cross-compile targets (for Go, each `GOOS`/`GOARCH` pair with `CGO_ENABLED=0`) — and run them in parallel. Jobs that cannot run on this machine (for example Windows-only tests on macOS) are listed as skipped with the reason. A tool that is not installed is skipped and noted in the PR body; an installed tool that fails is a failure.

On failure, ask: fix now (recommended — `manager-develop`, then re-run the mirror), push anyway with the failure stated in the PR body, or abort and keep the commit locally. Put a results table (check, status, notes) in the PR body.

### github-flow

- Feature branch or worktree branch: `git push -u origin <branch>`. If `gh pr list --head <branch> --json number` finds a PR, comment on it with the sync summary; otherwise `gh pr create --base <main_branch>` with a title from the SPEC and a body holding the summary, quality and audit results, and the deployment notes from Phases 2-6 (migrations, environment changes, breaking changes). When the SPEC frontmatter has a non-zero `issue_number`, end the body with `Fixes #<issue_number>`. Show the PR URL.
- The main branch: push directly only when it is not protected and the user approved; otherwise stop and ask.
- Any other state (for example a detached HEAD): stop, report the state and the routes above, and push nothing.

### git-flow

Evaluate in this order; the first match wins.

1. `WT-*` branch — integration-worktree merge, no PR. The integration branch is checked out in exactly one designated worktree; if that worktree does not exist, stop and report, since creating it is the coordinating session's job. Otherwise:
   1. Agree the merge window with the coordinating session, so one session integrates at a time.
   2. Enter the integration worktree; never check out the integration branch in the current one.
   3. `git merge --no-ff <branch>`; resolve conflicts, or report an unresolvable one back to the coordinating session instead of forcing it.
   4. `git push origin <integration-branch>` — never force. If the push is rejected, fetch, integrate, and push again.
   5. Leave the integration worktree and report the branch, merge commit, and evidence path.
2. `feature/*` — push, then create or update a PR to `develop`.
3. `release/*` — push, then create or update a PR to `main`.
4. `hotfix/*` — push, create or update a PR to `main`; after it merges, open a follow-up PR back-merging into `develop`.
5. `develop` — push directly (with approval).
6. The main branch — error: git-flow does not allow direct commits there; suggest a hotfix or release branch.
7. Anything else — stop, report the branch and this route list, and push nothing. Do not improvise a route or borrow another strategy's behavior.

### After the PR

- Team mode with a draft PR (`git_strategy.team.draft_pr`): `gh pr ready`. Otherwise leave readiness to the user.
- From the primary checkout, return to the base branch and pull (`develop` for git-flow feature branches, `main` for release and hotfix). In a worktree, stay where you are.
- When the SPEC has a non-zero `issue_number` and the user approved it, comment on the issue: completed, or partial with the current status. The PR's `Fixes #N` closes it on merge.

### Auto-merge

Only with `--auto-merge` (or the deprecated `--merge`, with a warning; `--no-merge` is a deprecated no-op, also with a warning). Worktree context alone never triggers a merge. In team mode, merge only after the required approvals; in personal and manual modes there is no approval condition. CI must pass and the PR must have no conflicts. `manager-git` runs: `gh pr ready`, `gh pr checks --watch`, `gh pr merge --<merge_method> --delete-branch`.

After a successful merge with `workflow.worktree.auto_cleanup: true`, clean up the SPEC's worktree with `moai worktree done <branch> --auto --delete-branch`. A cleanup failure does not affect the merge: log a warning with the manual command.

## Phase 14: Completion

Report: the strategy used, mode and scope, files updated, the SPEC status, the audit verdict (and the cross-model lines when the plan required them), deployment notes, and the PR URL or push target.

When sync ran as the tail of a full `/moai` pipeline and no decision is pending, end with the report — no next-step question. After an explicit `/moai sync`, offer next steps with AskUserQuestion, fitted to the result: after a PR — review it (recommended), auto-merge it (`/moai sync --auto-merge`), start the next SPEC (`/moai plan`); after a direct push or integration merge — start the next SPEC (recommended) or continue; in a worktree — review the PR, return to the main checkout, or remove the worktree.

<!-- moai:contract-mode-start id="contract-next-steps" -->
Where `workflow.autonomy.mode: contract` — no next-step question is asked; close with the completion report. A Phase 13 local CI mirror failure is an escalation report, not a question. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Escalation routing".

<!-- moai:contract-mode-end -->
## Edge cases

- Abort at any decision point: report exactly what was already changed — files written, commits made, anything pushed — and the retry command (`/moai sync [mode]`). Say "no changes" only when `git status` and the commit log confirm it.
- No SPEC named and uncommitted changes on the current branch: ask whether to sync those changes; with nothing to sync, report "Nothing to sync".
- Partial implementation: the SPEC stays `in-progress`, the documentation covers only what was built, and the PR body names the remaining work.

<!-- moai:contract-mode-start id="contract-error-flow" -->
Where `workflow.autonomy.mode: contract` — the current-branch confirmation is not asked: sync only the signed SPEC's changes, and escalate when uncommitted changes fall outside the contract's ownership. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Escalation routing".

<!-- moai:contract-mode-end -->
