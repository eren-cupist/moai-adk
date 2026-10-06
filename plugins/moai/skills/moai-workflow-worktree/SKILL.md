---
name: moai-workflow-worktree
description: >
  Command reference for MoAI git worktrees: creating, entering, syncing,
  verifying and disposing of worktrees for parallel SPEC or agent work.
  Not for deciding whether a run should be parallelized (see the run workflow).

when_to_use: >
  Use when creating, entering, syncing or cleaning up a git worktree with the
  moai CLI, or when working inside an isolated worktree.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Write, Grep, Glob
user-invocable: false
metadata:
  version: "1.1.0"
  category: "workflow"
  status: "active"
  updated: "2026-07-10"
  modularized: "false"
  tags: "git, worktree, parallel, development, spec, isolation"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000
---

# MoAI Worktrees

A worktree gives a SPEC, a session or an isolated agent its own checkout and branch, so parallel work never shares a working tree. The policy — L1/L2 layers, when agents get isolation, path rules for isolated agents, disposal safety — is in `.claude/rules/moai/workflow/worktree-integration.md`; operations inside a worktree session are in `worktree-integration-ops.md`. This skill is the command card.

## Create and enter

| Command | Effect |
|---|---|
| `moai worktree new <name>` | create an L1 tree at `.moai/worktrees/<name>` (base: `git_strategy.worktree_base_branch` when set) and print its path; does not enter it |
| `moai cc -w <name-or-abs-path>` | start a Claude session in the worktree (short names resolve under `.claude/worktrees/`); `moai glm -w` for GLM |
| `moai cc -w <name> --spawn` | the same in a new tmux window, keeping the current session; refuses outside tmux or without the `tmux`/`moai` binaries; POSIX only |
| `moai cc -w <name> --branch <existing>` | enter an existing branch |
| `claude -w [name]` | a native ephemeral worktree session under `.claude/worktrees/` |
| `EnterWorktree(<path>)` / `ExitWorktree` | move the current session into or out of a worktree |

Inspect with `git worktree list`.

## Maintain and verify

| Command | Effect |
|---|---|
| `moai worktree sync [branch] [--base main] [--strategy merge\|rebase]` | bring the worktree up to date with its base; on a conflict, resolve with plain git in the worktree and rerun |
| `moai worktree snapshot` / `verify` / `restore` | guard the main checkout around an isolated agent call (`worktree-state-guard.md`) |
| `moai worktree recover` | repair the worktree registry |

## Dispose

Hoist evidence before removing a tree: `.moai/reports/` inside it is gitignored and exists nowhere else.

| Command | Effect |
|---|---|
| `moai worktree hoist <tree-path>` | copy the tree's `.moai/reports/` into the root's `.moai/reports/worktrees/<tree-name>/` |
| `moai worktree done [branch] [--delete-branch] [--force] [--no-hoist]` | finish an L2 SPEC tree after both the run and sync PRs merged (hoists first; refuses L1 trees) |
| `moai worktree sweep [--base origin/develop] [--yes] [--json]` | dispose of trees whose branches landed on the remote base, keeping any tree that is dirty, locked, anchored or holding irreplaceable content |
| `moai worktree clean [--stale] [--merged-only] [--yes] [--json] [--base origin/main]` | remove stale or merged clean trees (`--stale` previews unless `--yes`) |
| `moai worktree remove [path] [--force]` | remove one tree |

`moai worktree --help` is authoritative where it differs from these tables.
