---
paths: "**/.claude/agents/**,**/.claude/worktrees/**,**/.moai/worktrees/**,**/.claude/teams/**"
---

# Worktree Integration

How MoAI work uses git worktrees: which kind to use, how to create and enter one, how agents are isolated, and how a tree is disposed of without losing work. Operations inside a worktree session (refused commands, anchor recovery, base-branch settings, worktree hooks) are in `worktree-integration-ops.md`.

## Terminology Glossary

| Layer | What it is | Path | Lifetime |
|---|---|---|---|
| **L1** session worktree | Isolation for one session or one isolated subagent. Created by `moai worktree new <name>`, by `claude -w` / `EnterWorktree`, or automatically for an `Agent(isolation: "worktree")` call. | `.moai/worktrees/<name>/` (MoAI-created), `.claude/worktrees/<name>/` (Claude-native) | Until its session ends and its branch is integrated |
| **L2** persistent SPEC worktree | A long-lived working directory for one SPEC, reused by run and sync. | `~/.moai/worktrees/<project>/<SPEC>/` | Until both the run and sync PRs merge; disposed with `moai worktree done <SPEC>` |

L1 trees are outside the L2 lifecycle registry: `moai worktree done` refuses both L1 roots, with or without `--force`.

## Creating and entering

- `moai worktree new <name>` creates an L1 tree at `.moai/worktrees/<name>`, cut from `git_strategy.worktree_base_branch` when that is set, and prints its absolute path. It does not enter it.
- `moai cc -w <name-or-absolute-path>` (likewise `moai glm -w`) starts a new Claude session inside a worktree; a short name resolves under `.claude/worktrees/`, so use the absolute path for `.moai/worktrees/` and L2 trees. `--spawn` opens the session in a new tmux window and keeps the current one (it refuses, rather than falling back, outside tmux or without the `tmux` and `moai` binaries). `moai cc -w <name> --branch <existing>` enters an existing branch.
- `claude --worktree [name]` (`-w`) starts a fresh native session in `.claude/worktrees/<name>/`; at session end it offers to keep or remove the tree.
- Inside a running session, `EnterWorktree(<path>)` moves the session into an existing tree and `ExitWorktree` returns to the original checkout. Use these rather than `cd`, `git -C` or a subshell `cd` in instructions the orchestrator emits for the current session. Codex uses `codex -C <absolute-path>` for a new session and `git -C <path>` within one.
- Never create trees with a bare `git worktree add`, and never copy a worktree directory by hand; `moai update` moves registered legacy `.claude/worktrees/` trees to `.moai/worktrees/` with `git worktree move` (`--dry-run` lists the moves).

## Agent isolation

`Agent(isolation: "worktree")` — or `isolation: worktree` in an agent's frontmatter, as manager-develop declares — gives each spawn a new L1 tree under `.claude/worktrees/<auto-name>/`, branched from the default remote branch (see `worktree.baseRef` in the ops file). The agent's working directory is that tree's root. It is not a way to re-enter an existing worktree; use `EnterWorktree` or `moai cc -w` for that.

Use isolation for agents that change files in parallel with other writers or that make cross-file changes, so two writers never share a tree. Read-only agents need no worktree, but an agent is read-only only when no tool in its list can change files: omitting `Write` and `Edit` is not enough while it has `Bash`, a file-changing MCP tool, or `Agent`. `moai agent lint` reports `ORC_WORKTREE_MISSING` (LR-05) for a write-heavy agent without isolation and `ORC_WORKTREE_ON_READONLY` (LR-09) for a `permissionMode: plan` agent with it. Worktree working-directory isolation needs Claude Code 2.1.97 or later.

### Prompt paths for isolated agents

In a prompt for an isolated agent, refer to files by project-root-relative paths (`src/auth/handler.ts`, `.moai/specs/SPEC-XXX/spec.md`) and give commands without a `cd` prefix (`pnpm test`, not `cd /path/to/repo && pnpm test`). An absolute path into the main checkout, or a `cd` into it, makes the agent read and edit the main tree instead of its own. `$CLAUDE_PROJECT_DIR` in hook commands is fine; Claude Code resolves it per agent.

## Parallel-Session Branch Conflict Auto-Isolation

When `.moai/state/active-sessions.json` shows another active session on the same checkout while you are about to change files — through a spawned agent or directly — move your work into a worktree instead of sharing the branch state: one tree per foreign session entry, named `auto-<first 8 chars of that session's id>-<SPEC-ID>`, under `.claude/worktrees/` or `~/.moai/worktrees/<project>/`. Any foreign entry triggers it; a stale entry costs only an extra tree the user can delete. Mention it in a one-line note (including the entry's age) rather than asking.

## Branch names

`EnterWorktree(<name>)` names the branch `worktree-<name>`. A branch renamed to `WT-<slug>` (lowercase `a-z0-9-`, at most three tokens and 24 characters, describing the change) becomes a candidate for the PR-merge auto-cleanup sweep when `workflow.worktree.auto_cleanup` is enabled (off by default): the sweep removes a `WT-` tree whose branch reads merged and whose tree is clean and not anchored by a live session. Other branch names are never swept; they are disposed of manually.

## Disposing of a worktree

- An unpushed worktree branch is the only copy of its work. Dispose of no tree until its branch is integrated and the merge has landed on the remote.
- `.moai/reports/` inside a tree is gitignored and machine-local, so removing the tree destroys that evidence even when `git status` is clean. Run `moai worktree hoist <tree-path>` before any removal; it copies the evidence into the project root's `.moai/reports/worktrees/<tree-name>/` and never overwrites. `moai worktree done` hoists automatically (`--no-hoist` skips it).
- `moai worktree sweep` previews, and `--yes` performs, disposal of trees whose branches are confirmed landed on the remote base (`origin/develop` by default, `--base` to change), hoisting first and keeping any tree that is protected, locked, anchored by a live session, occupied by a process, dirty, or holding irreplaceable ignored content.
- `moai worktree clean --stale` (preview; `--yes` to act; `--json` to report only) removes abandoned clean trees with no unique commits; `--merged-only` removes trees whose branches are merged into the base (`origin/main` by default).
- For a tree none of these reach: `git worktree unlock <path>` if a dead session still holds the lock, then `git worktree remove <path>`. Branches are never deleted by these commands.

## SPEC-to-Worktree Mapping

The canonical per-step rules are in `.claude/rules/moai/workflow/spec-workflow.md` (SPEC Phase Discipline); on conflict that file wins.

| Step | Phase | Worktree | Notes |
|---|---|---|---|
| 1 | Plan | no — main checkout | plan artifacts are markdown and need cross-SPEC visibility |
| 2 | Run | optional (`moai cc -w <name>`) | the entered tree |
| 3 | Sync | same tree as run | never a new tree; sync works on the run-modified state |
| 4 | Cleanup | — | `moai worktree done SPEC-XXX`, only after both the run and the sync PR have merged |

Worktree use is the user's choice; by default every phase runs on a feature branch in the main checkout.
