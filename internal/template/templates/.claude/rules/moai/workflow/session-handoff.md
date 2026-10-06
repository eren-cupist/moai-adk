# Session Handoff

When work will continue in a fresh session, give the user a paste-ready resume message so nothing is lost at the boundary.

## When to emit one

- Context usage reaches the threshold in `context-window-management.md`.
- A plan, run or sync phase finishes and more SPECs in the same series remain.
- The user says the session is ending.
- A multi-milestone task reaches a stable checkpoint with milestones still to go.

For single-turn, trivial or read-only work, a short completion note is enough.

## The message

Render it in the response body as a fenced `text` block between cut-line markers, so it can be copied cleanly out of long terminal scrollback. Writing it to memory or saving it with the CLI is not the same as showing it; the user can only paste what was rendered.

```text
✂──── 여기부터 복사 ────✂

<SPEC-ID> <phase> 진입.
applied lessons: <memory-file>, <memory-file>
source_session_id: <uuid from `moai session current`>

전제 검증:
1) <check> → <expected result>
2) <check> → <expected result>

실행: <the one command that starts the work>

머지 후: <the next SPEC or /moai sync>

✂──── 여기까지 복사 ────✂
```

The example is the Korean rendering. Marker text and labels follow `conversation_language`; the `✂` and `─` characters stay as they are in every language.

| Element | English | Korean |
|---|---|---|
| Top / bottom marker text | `Copy from here` / `Copy to here` | `여기부터 복사` / `여기까지 복사` |
| Opener verb | `entering` | `진입` |
| Preconditions header | `Preconditions:` | `전제 검증:` |
| Run header | `Run:` | `실행:` |
| Closing header | `After merge:` (PR flow) or `Follow-up:` (no PR) | `머지 후:` or `후속:` |

`<phase>` is `plan`, `run` or `sync`. List at most four preconditions, each checkable by a command or a file check. `Run:` carries a single action. Omit the closing block when nothing is queued, and omit `source_session_id` when `moai session current` is unavailable. When the next session should start in a non-default orchestration mode, add `mode: fanout` or `mode: agent-team` after the opener (for `agent-team`, also append `--team` to the run command); a mode line seeds the next session and never grants the kickoff approval. When the work happened in a worktree, put the entry command first (`moai cc -w <name>` or `moai cc -w <absolute-path>`, not a bare `cd`) and make the first precondition `git rev-parse --show-toplevel` → the worktree path.

Keep it to what the next session needs to start; it is not a history of this one.

## Saving it

Alongside rendering, save the block as the pending handoff: write it to a file, then run `moai handoff save --stdin --spec <SPEC-ID> --phase <phase> --lang <conversation_language> [--session <uuid>] < <file>`. Pass the body through a file rather than an inline heredoc, because a heredoc that mentions a git command is refused by the branch and worktree guards. If `moai` is missing or the save fails, still render the message unchanged and do not retry. The record lands in `.moai/state/handoff/pending.json`; with `handoff.mode: auto` in `.moai/config/sections/handoff.yaml` it is injected automatically at the next `/clear`, and with the default `manual` it is left alone. Either way the plan-to-run kickoff approval still applies in the resumed session.

Also save the message in auto-memory as a project topic file under a `## Next Session Entry Point` heading (`## 다음 세션 시작점` in Korean) with its `MEMORY.md` index line, so it is discoverable after `/clear`; mark a replaced entry `[SUPERSEDED by <new-file>]`. When the SPEC closes, the verbatim block in that file can shrink to a one-line summary.
