---
name: moai-ref-git-workflow
description: >
  Reference for MoAI git delivery: branch naming from the git-strategy
  config, the two delivery strategies, merge-method selection, the PR body
  template, the commit context section, and the git safety rules. Loaded by
  manager-git. Not for writing code, tests, or documentation content.

when_to_use: >
  Use when choosing a branch name, a merge method, or a PR body for a MoAI
  SPEC, or when checking whether a git operation is safe to run.

user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-03-30"
  tags: "git, branch, commit, pr, workflow, reference"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 3000
---

# Git Workflow Reference

Settings come from `.moai/config/sections/git-strategy.yaml` (resolve `git_strategy.mode` first, then read `git_strategy.{mode}.*`), the commit convention from `git-convention.yaml`, and the commit language from `language.yaml` `git_commit_messages`.

## Strategies

`git_strategy.{mode}.workflow` is exactly one of:

- **github-flow** — feature branches off `main_branch`; every change reaches it through a PR; branches are deleted after merge.
- **git-flow** — `feature/*` → PR to `develop`; `release/*` and `hotfix/*` → PR to `main` (hotfix then back-merged to `develop`); `WT-*` worktree branches → merged into the integration branch inside the integration worktree, holding `moai integration acquire` / `release`; no direct commits on `main`.

The full routing, in order, is in `.claude/skills/moai/workflows/sync/delivery.md`.

## Branch names

With `automation.auto_branch: true`, a SPEC branch is `<branch_prefix><SPEC-ID>` (the shipped prefix is `feature/SPEC-`, giving `feature/SPEC-AUTH-001`). Work in a Claude Code worktree runs on `worktree-<name>`. Branches outside a SPEC follow the change type:

| Pattern | Example |
|---------|---------|
| `fix/<SPEC-ID>-<slug>` | `fix/SPEC-BUG-042-null-check` |
| `refactor/<slug>` | `refactor/extract-auth-middleware` |
| `docs/<slug>` | `docs/api-reference-update` |
| `chore/<slug>` | `chore/upgrade-dependencies` |

## Merge method

Every PR merge uses `git_strategy.{mode}.merge_method` (`squash`, `merge`, or `rebase`; default `squash`) as `gh pr merge --<merge_method> --delete-branch` — never a hardcoded method. Squash keeps feature history to one commit per PR; merge preserves the individual commits (useful for release branches); rebase suits a few clean commits.

## PR body

```markdown
## Summary
- <what this PR does, 1-3 bullets>

## Changes
- <file or area>: <change>

## Test plan
- <tests added or updated, and what was run>

## Local CI mirror
| Check | Status | Notes |
|-------|--------|-------|

## Deployment notes
- <migrations, new environment variables, breaking changes — or "none">

## SPEC
- <SPEC-ID>: <title>

Fixes #<issue_number>
```

Include `Fixes #<issue_number>` only when the SPEC frontmatter has a non-zero `issue_number`. Create it as a draft when `git_strategy.team.draft_pr` is true.

## Commit context section

Implementation and sync commits carry the decisions a later session needs:

```
feat(SPEC-AUTH-001): M2 rotate refresh tokens

## Context (AI-Developer Memory)
- Decision: rotation instead of a sliding window (a stolen token stops working after one use)
- Constraint: the token store must be shared across instances
- Gotcha: the in-memory blacklist does not survive restarts; moved to Redis
- Pattern: middleware chain RateLimit -> Auth -> Authz -> Handler
```

## Safety

| Operation | Rule |
|-----------|------|
| `git push --force` | never on `main`, `develop`, or any shared branch; on your own branch use `--force-with-lease`, with approval |
| `git reset --hard`, `git checkout .`, `git clean -fd` | discard work: only with the user's approval |
| `git branch -D` | only after the merge is confirmed |
| `--no-verify` | never, unless the user explicitly asks |
| `git rebase -i`, `git add -i` | not available (interactive) |
| commit on a protected branch | never — see `workflow.branch_guard.deny_commits_on`; use the SPEC branch or worktree |
| push, PR, merge, issue comment | only with the user's request or approval |
