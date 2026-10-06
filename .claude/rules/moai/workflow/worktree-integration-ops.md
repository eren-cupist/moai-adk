---
paths: "**/.claude/worktrees/**"
---

# Worktree Integration — Session Operations

Operations for a session working inside a worktree: commands the worktree guard refuses, recovering a misresolved session anchor, the base-branch and copy settings of native worktrees, and the worktree hooks. Creation, isolation and disposal are in `worktree-integration.md`.

## Refused Commands in a Worktree-Isolated Session

Two different guards refuse commands in a worktree-isolated session, and the message is the only way to tell them apart:

| The refusal says | Whose guard | Changeable here? |
|---|---|---|
| `Dangerous command blocked:` … | MoAI's own pre-tool hook | yes |
| `This session is isolated in the worktree <path>, but this command …` | the Claude Code runtime | no — nothing in this repository implements or configures it |

The runtime guard's clause after "but this command" varies — `redirects git to the shared checkout via -C` (or `via --git-dir`), `is too complex to verify that it stays inside the worktree`, `names git in a form too complex to verify`, `points git at a directory computed at runtime`, `runs <command> with <argument> in a plain command, so what it runs cannot be shown not to be git`, and others. MoAI's failure observer matches only the opening `isolated in the worktree`. Treat an unfamiliar clause as one more shape of the same guard, record its exact wording and the `claude --version`, and do not assume the guard is consistent across versions — observed behavior has changed between releases.

What has been seen to trigger the runtime guard: a git call it cannot bind statically to the worktree (git inside `$()` whose result is expanded later, git nested in another git command's arguments, git piped into a loop, git inside a `( … )` subshell, a git-bearing assignment used in a later statement), several mutation steps bundled into one compound command, and a quoted-delimiter heredoc whose body holds braces around quoted key/value pairs (a JSON line) even with no git in the command. Plain, separately invocable git verbs, `;`/`&&` chains of them, and redirecting output to a file pass. Anything not observed is unknown, not permitted.

Workarounds, by situation:

- **Writing file content** — use the `Write` tool instead of a heredoc; it issues no shell command, so the guard never parses the content.
- **A compound command** — split it into separate plain commands. Keep together only pairings that must share one process, such as an `unset … && <command>` environment scrub.
- **Do not** move a refused command into a script file to get past the guard: the guard cannot see inside it, so the risk is hidden, not removed.

Acceptance-criterion verification commands that must run in a worktree follow the same shape: one plain git verb per line run from the worktree root; capture exit codes on their own line (`git diff --quiet …; echo "exit=$?"`) rather than in a `$()` variable; derive ranges with the three-dot form (`git diff --name-only develop...HEAD`) and record `git merge-base develop HEAD` on its own line when the base is evidence; keep scratch writes under `/tmp` and cleanup as a separate step; pin the tree SHA the measurement was taken on (`verification-completeness.md`, sections 2.1 and 4).

```bash
# Refused: assignment, then later expansion of a git-bearing substitution
B=$(git merge-base develop HEAD)
git diff --name-only "$B"..HEAD -- src/pkg/ | wc -l

# Runs: plain verbs, base recorded on its own line, three-dot range
git merge-base develop HEAD
git diff --name-only develop...HEAD -- src/pkg/ | wc -l
git diff --quiet develop...HEAD -- src/pkg/; echo "diff_exit=$?"
```

## Background agents and the session anchor

A background subagent is not pinned to the tree it was spawned in: its working directory resolves against the session's current anchor at each call. If the parent session moves to another tree, a background agent's path-less commands keep succeeding — silently — in the new tree. Do not move the session's tree while a background agent is running; if it happened, treat any verification it ran during the move as unaccounted for until you have checked which tree it ran in.

## Session anchor misresolution

The runtime's session anchor can resolve to the wrong tree (another lane's worktree), or refusals can appear for commands in your own tree. Recover by calling `ExitWorktree` (answer **keep**) and then `EnterWorktree(<your worktree path>)`; re-entry registers a fresh anchor. Waiting, retrying the command, or restarting a subagent does not clear it. The anchor itself is runtime state; MoAI only observes it.

To make the next occurrence attributable: `MOAI_ANCHOR_TRACE=1` writes one JSONL row per anchor decision (session id, pid, cwd, time) to `.moai/logs/anchor-trace.jsonl`; every registry `cwd` rewrite is appended to `.moai/logs/anchor-relocation-audit.jsonl`; and `workflow.anchor_relocation_guard.enabled` refuses an ownership-flagged relocation instead of proceeding.

## Native worktree settings

**Base branch.** Native worktrees (`claude -w` and `isolation: worktree`) branch from `origin/HEAD`, falling back to local `HEAD` without a remote. The Claude Code setting `"worktree": {"baseRef": "head"}` (only `"fresh"` or `"head"` are accepted) makes them branch from the HEAD of the tree the launcher ran in — note that this is the launching tree's HEAD, not the branch you may have in mind. `claude --worktree "#1234"` creates `.claude/worktrees/pr-1234` from that pull request.

`git_strategy.worktree_base_branch` in `.moai/config/sections/git-strategy.yaml` names the integration branch reproducibly: at session start in the primary checkout, MoAI points `refs/remotes/origin/HEAD` at it, and `moai worktree new` passes it as the base. Empty (the default) changes nothing; a value without a remote-tracking branch is refused. `moai doctor --check 'Worktree Base Branch'` reports the current alignment without changing it. For a new, still-empty tree on the wrong base, `git -C <tree> reset --hard <ref>` and then confirm with `git -C <tree> merge-base --is-ancestor <ref> HEAD`; never do this once the tree has commits of its own.

**Copying ignored files.** A native worktree is a fresh checkout without untracked files such as `.env`. A `.worktreeinclude` file at the project root (gitignore syntax) lists gitignored files to copy into new native worktrees — for example `.env`, `.env.local`, `.moai/config/sections/*.local.yaml`. It is not processed when a custom `WorktreeCreate` hook replaces the default behavior.

## WorktreeCreate and WorktreeRemove hooks

MoAI does not register these hooks. In Claude Code, `WorktreeCreate` replaces the default git worktree creation: the hook must create the directory itself and print only its absolute path on stdout; empty output or a non-zero exit aborts creation. `WorktreeRemove` is an observer. Registering a stub would break worktree creation. A custom creator, if ever needed, reads the stdin JSON (`worktree_path`, `name`, `cwd`, `session_id`), runs the creation with its own output sent to stderr or `/dev/null`, prints the path, and exits 0. The handlers `moai hook worktree-create` / `worktree-remove` and the `.claude/hooks/moai/handle-worktree-{create,remove}.sh` wrappers exist for that case but are not wired in `settings.json`.
