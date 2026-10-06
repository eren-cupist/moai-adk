---
description: "The brief format for delegating SPEC implementation to manager-develop. Load when writing a manager-develop spawn prompt."
paths: ".moai/specs/**,.claude/agents/moai/manager-develop.md,.claude/skills/moai/workflows/run.md"
---

# manager-develop Brief

A manager-develop agent starts with none of the conversation's context. Everything it needs to do the job without guessing goes into one brief, written once; a vague brief ("implement the SPEC") produces silent assumptions and a second round trip. A small Tier S track needs only a few lines per section; a Tier L track spells each section out.

## Sections

1. **Task and location.** The SPEC path (`.moai/specs/SPEC-XXX/`) and which files to read there; the branch the work builds on; the `cycle_type` — `tdd` for new behavior, `ddd` for changing code that lacks tests (from `quality.yaml` `constitution.development_mode`, or the track's nature), or `autofix` for a failing required CI check; and the skill to load first (`moai-workflow-tdd` or `moai-workflow-ddd`, plus a domain skill such as `moai-ref-react-patterns` or `moai-ref-api-patterns` when it helps).
2. **Acceptance criteria in scope.** The AC IDs this agent owns, quoted or pointed at in `acceptance.md`/`spec.md`. When the run is split across agents, say which criteria belong to other tracks.
3. **Files and areas owned.** The files or directories this agent may change, as project-root-relative paths. Name what it must not touch: other tracks' files, unrelated untracked files, runtime state under `.moai/state/` and `.moai/cache/`, SPEC body content.
4. **Constraints.** Anything the codebase or the user requires that the SPEC does not say: existing contracts to preserve (an API shape, an `@MX:ANCHOR` function), conventions to follow, dependencies not to add, the project's package manager and runner (e.g. pnpm with turbo). Include the facts the orchestrator already established so the agent does not rediscover them, and any decision the user made at Kickoff.
5. **How to verify.** The exact commands that prove the work: change-scoped tests, the typecheck and lint tasks, the build if relevant, and what "passing" looks like for each acceptance criterion.
6. **What to return.** Commits on its worktree branch and the branch name; for each owned AC, PASS or FAIL with the command and its observed output; the failing output each new test showed before its implementation (TDD); files touched versus planned, with any scope change and the reason; and a blocker report for anything it could not decide — never a question to the user, since a subagent cannot ask one.

## cycle_type=autofix

For a failing required CI check, the brief names the pull request and branch, the failing check and its log, and the iteration number. The agent diagnoses the failure, applies the smallest patch that fixes the root cause, re-runs the check locally, and commits the patch as a new commit. The limits — three iterations, no auto-patch for semantic failures, no edits to secrets or CI definitions, no force-push — are in `.claude/rules/moai/workflow/ci-autofix-protocol.md`.
