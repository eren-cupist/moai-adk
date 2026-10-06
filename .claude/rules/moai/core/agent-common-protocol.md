---
description: Shared protocol loaded in every MoAI session and subagent — the user-interaction boundary, blocker reports, hook signals, and language of reports. Intentionally always-loaded (no paths restriction).
---

# Agent Common Protocol

This rule loads in the main session and in every MoAI agent, so each agent definition can rely on it instead of repeating it.

## User Interaction Boundary

Subagents cannot ask the user questions: they run in an isolated context with no question channel, and the orchestrator is the user's single point of contact. Everything a subagent needs from the user arrives in its spawn prompt. When required input is missing, the subagent stops and returns a blocker report instead of asking in prose or writing out an AskUserQuestion call:

```markdown
## Missing Inputs

| Parameter | Expected values | Why it is needed |
|-----------|-----------------|------------------|
| [name]    | [values]        | [reason]         |

**Blocker**: cannot proceed without these inputs. Re-delegate with them in the prompt.
```

The orchestrator answers a blocker report by asking the user and re-delegating with the answer. A blocker report is a return; when a delegation aborts with no return at all, the orchestrator states what was delegated, that it did not return and why, before moving on, and never presents it as a success.

## Hook signals

Hooks never ask the user either. A blocking hook exits 0 and prints JSON on stdout with `"decision": "block"` and a `reason` (optionally `systemMessage` or `details`); on exit 2 that JSON is discarded and only stderr surfaces. Read the reason and fix the gate that failed; if overriding the block looks warranted, the orchestrator brings that decision to the user. `sync-phase-quality-gate.sh` blocks a sync-phase commit only when `MOAI_SYNC_GATE_BLOCKING=1`; otherwise it reports. The Stop hook fires at every turn end and stays silent unless the turn is a real completion point, and during recovery turns (compaction, `prompt_too_long`, `max_output_tokens`) hooks exit 0 so recovery does not loop; see `runtime-recovery-doctrine.md`.

## Language of reports

Analysis, reports and documentation for the user are written in the user's `conversation_language`. Code, identifiers, skill names, paths, commands and flags stay in their original form. Code comments follow `code_comments` and commit messages follow `git_commit_messages` in `.moai/config/sections/language.yaml`. Prompts passed between agents stay in English. User-facing output is Markdown; XML is only for structured data passed between agents.
