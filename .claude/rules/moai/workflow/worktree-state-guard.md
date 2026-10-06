---
paths: "**/.moai/specs/**,**/internal/worktree/**,**/internal/cli/worktree/**"
---

# Worktree State Guard

`moai worktree snapshot`, `verify` and `restore` let the orchestrator check that an `Agent(isolation: "worktree")` call did not disturb the main checkout. They are optional: use them around a write-capable isolated agent when the main tree's state matters (an in-progress branch, uncommitted work) or when isolation has misbehaved before. Skip them for read-only agents and trivial single-file edits. Layer terms (L1/L2) are defined in `worktree-integration.md`, section "Terminology Glossary".

## Sequence

```bash
# Before the isolated agent call
moai worktree snapshot --agent-name <agent>      # prints JSON including the snapshot path

# After it returns
moai worktree verify --snapshot <snapshot-path> --agent-name <agent> \
  [--agent-response <response.json>]
```

`verify` compares HEAD, the current branch, untracked files under `.moai/specs/`, and `git status --porcelain` (ignoring `.moai/state/`, `.moai/reports/`, `.moai/cache/`, `.moai/logs/`). Any change counts as divergence. Its exit code and JSON report (`divergence`, `suspect_flag`, `report_path`):

| Exit | Meaning | Action |
|---|---|---|
| 0 | clean | continue |
| 1 | divergence | read the report, then ask the user through AskUserQuestion: restore, accept, or abort |
| 2 | suspect — the agent response had an empty `worktreePath`, so isolation may not have happened | tell the user and hold any push until it is understood |
| 3 | both | handle as 1 and 2 |

## Restore

`moai worktree restore --snapshot <snapshot-path>` resets tracked files to the snapshot's HEAD state. It is destructive, so run it only after the user chose restore. Untracked files are listed but not recreated.

The CLI never asks the user anything; it reports exit codes and JSON, and the orchestrator turns them into the question.
