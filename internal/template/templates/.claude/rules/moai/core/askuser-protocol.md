---
description: How to ask the user a question with AskUserQuestion — loading the deferred tool, payload limits, option format, recommendation mode, and the report-before-ask rule. Intentionally always-loaded (no paths restriction).
---

# AskUserQuestion Protocol

When a question is genuinely the user's to answer (see `CLAUDE.md`, section "Questions"), the orchestrator asks it with `AskUserQuestion` rather than in prose. The tool appends an "Other" choice to every question, so a free-form answer is always possible without a prose question. Subagents never call it; they return a blocker report (`agent-common-protocol.md`).

## Loading the tool

`AskUserQuestion` is a deferred tool: its schema is not loaded at session start, and calling it unloaded fails with `InputValidationError`. Load it once per session with `ToolSearch(query: "select:AskUserQuestion")` before the first call. `TaskCreate`, `TaskUpdate`, `TaskList` and `TaskGet` are deferred the same way.

## Composing a question

- At most 4 questions per call and 4 options per question (Claude Code limits).
- Question text, headers, labels and descriptions are in the user's `conversation_language`, written as native UTF-8 rather than `\uXXXX` escapes; a malformed escape turns the payload into a validation error.
- Each option's `description` says what happens if it is chosen, including side effects, risk and irreversibility, in neutral language.
- Put the recommended option first and mark its label with `(권장)` for Korean or `(Recommended)` for English. Recommend the choice the user would most likely make, and say in its description when the recommendation would not hold.
- `preview` adds a monospace block beside the options for comparing concrete alternatives (code, layouts, configs). It works only with single-select questions and does not scroll, so keep it to about 12 lines and keep it neutral.

When `interview.recommendation_mode` in `.moai/config/sections/interview.yaml` is `pull`, leave the recommendation label off every option and don't signal a preference through ordering, wording or preview; still present the evidence and the options, and give the recommendation when the user asks for one. `push` (the default) behaves as described above.

## Report before asking

When the options of a decision come from investigation (agent results, audits, verification runs), put the findings in the response body before the question, so the user can judge the options; the option text alone is not a report. When the user asked for a report, analysis or explanation, that is the deliverable: give it and end the turn without appending a decision question. A completion report ends either with an `AskUserQuestion` for a real decision or with no question at all — not with a prose "what next?".

<!-- moai:contract-mode-start id="contract-ambiguity" -->
Where `workflow.autonomy.mode: contract` — after the contract is signed, an ambiguity the contract already answers is recorded, not asked; an ambiguity that contradicts the contract escalates instead. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Gate disposition".

<!-- moai:contract-mode-end -->
