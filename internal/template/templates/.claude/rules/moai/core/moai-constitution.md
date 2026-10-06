# MoAI Constitution

Principles that apply across every MoAI workflow. Orchestration, questions and safety rails are in `CLAUDE.md`; the cross-harness contract (evidence, git, code changes, output) is in `AGENTS.md`. This file holds what those two do not.

## Quality gates

Code changes are held to TRUST 5: Tested (coverage target from `quality.yaml`; characterization tests before changing untested code), Readable (clear names; comments in the `code_comments` language and the surrounding density), Unified (the project's formatter), Secured (validated external input, no secrets in code) and Trackable (conventional commits). Thresholds live in `.moai/config/sections/quality.yaml`, and `sync-auditor` checks the result against them.

## @MX annotations

Agents add, update and remove `@MX` tags as part of their work without asking: `@MX:ANCHOR` on functions with three or more callers, `@MX:WARN` with `@MX:REASON` on dangerous patterns, `@MX:NOTE` for context a reader needs, `@MX:TODO` for untested public functions, and `@MX:DEBT` with `@MX:CEILING` / `@MX:UPGRADE` for deliberate simplifications. Report tag changes in the completion report. The tag grammar is in `mx-tag-protocol.md`, which loads with source files.

## Agents in an isolated worktree

When spawning an agent with `isolation: "worktree"`, write project-root-relative paths in its prompt and no `cd /absolute/path &&`. The agent starts in its own worktree root, and an absolute path into the main checkout would make it write there instead. Worktree lifecycle rules are in `worktree-integration.md`.

## Web sources

Cite only URLs that came from WebSearch or that you opened with WebFetch, mark unverified information as uncertain, and add a Sources section when you used web search.

## Memory

Lessons from user corrections go into Claude Code's auto-memory as topic files indexed by `MEMORY.md`. Write memory files in English. Every topic file needs exactly one index line, because topic files are found only through the index; when the index grows, make lines shorter rather than removing them, and mark a replaced entry `[SUPERSEDED by <new-file>]` instead of deleting it. `moai memory doctor` reports the store paths and index size. Agent memory rules are in `moai-memory.md`.
