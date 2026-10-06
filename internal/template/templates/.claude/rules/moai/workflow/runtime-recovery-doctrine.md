---
paths: "**/.moai/specs/**,**/runtime-recovery-doctrine.md"
---

# Runtime Recovery Doctrine

What to do when the session itself fails partway through a SPEC — the context fills up, an output is cut off, a compaction fails. These are policy notes for agents and for hook authors; MoAI runs on top of Claude Code and does not change its query loop or its own compaction.

## Recoverable errors

`prompt_too_long`, `max_output_tokens`, `media_size` and a failed compaction are recoverable. Hitting one mid-turn is not a failed turn: recover with the cheapest step that works before reporting failure.

## Recovery ladder

Try the cheapest rung first and move up only when it has been tried for the current failure:

1. **Continue in the turn** — correct course and keep working, without recap or apology.
2. **Resume message + `/clear`** — write a paste-ready resume message (format in `.claude/rules/moai/workflow/session-handoff.md`) and clear the context.
3. **Restart in a fresh worktree** — when the problem is structural (a stale worktree base, a diverged tree, accumulated compaction residue), resume in a clean worktree using the worktree-anchored resume in `session-handoff.md`.
4. **Abort and preserve** — persist the in-flight state to the SPEC's `progress.md` and end the session, so the next one resumes instead of restarting.

Three limits keep recovery from looping: after three consecutive failures at one rung, move to the next; never repeat a recovery step that already failed in the same turn; and when a compaction itself overflows the context, go straight to rung 4. Across a compaction or restart, keep the account of what was tried, why it failed and what remains, in the evidence format of `.claude/rules/moai/core/verification-claim-integrity.md`, so the next turn does not retry a known failure. An abort that leaves `progress.md` stale forces the next session to rediscover everything.

## Anti-death-spiral hook carve-out

### Recovery-Signal Carve-Out

Stop and PostToolUse hooks should exit 0 on a recovery turn (letting the turn end or the tool call proceed) rather than exit 2, so a recovery turn is not caught in an error → hook blocks → retry → error loop. A recovery turn is one recovering from a sync failure, a compaction, `prompt_too_long`, `max_output_tokens`, `media_size` or a failed compaction.

This is guidance for hook authors, not enforced behavior: the current hooks (`sync-phase-quality-gate.sh` on Stop, `status-transition-ownership.sh` on PostToolUse) do not read a recovery signal from `stopReason`, so they cannot tell a recovery turn from a normal one. On normal turns they keep blocking genuine gate failures; the carve-out never weakens that. Do not describe the carve-out as mechanically enforced.
