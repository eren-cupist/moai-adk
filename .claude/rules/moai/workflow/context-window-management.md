# Context Window Management

Plan a clean handoff before the context window fills, because a stream that runs out of room mid-task can stall and lose work.

| Session window | Hand off at |
|---|---|
| 1M tokens (Opus 5.5 and other 1M models on the Anthropic API) | 50% |
| 200K tokens (any session running a 200K window, including 1M models under `CLAUDE_CODE_DISABLE_1M_CONTEXT=1`, and Haiku) | 90% |

The window the session actually runs with decides the row, not the model name. Above 95% on any window, the next action is `/clear`.

Read usage from the statusline snapshot `.moai/state/context-usage/<session-id>.json`: `raw_pct` is the raw usage and `stage` is `none`, `soft` or `hard`. It relays Claude Code's own percentage, so it is accurate to about a point. Without it, estimate from the volume of output, large tool results and returned agent reports, and err early: a premature `/clear` costs one paste, a missed one costs a stalled stream.

When usage reaches the threshold, stop starting large tool calls or delegations, save in-flight progress to `.moai/specs/<SPEC-ID>/progress.md`, emit the resume message described in `session-handoff.md`, and tell the user that `/clear` followed by pasting it will continue the work. That is a status note, not a question. While `interview.recommendation_mode` is `pull`, report the usage and name `/clear` as the action that resets it without recommending it.

Cheaper moves come before `/clear` when the task must continue: `/compact <instructions>` summarizes in place, and `/rewind` (Esc Esc) can summarize up to a point or restore a checkpoint. `/btw` answers a side question without adding it to the transcript.
