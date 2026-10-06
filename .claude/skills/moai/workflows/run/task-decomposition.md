---
description: "Run phase — implementing (directly or through manager-develop), verification, progress.md, git, and completion"
user-invocable: false
metadata:
  parent: moai-workflow-run
  phase: "Implementation, verification, and completion"
---

# Implementation

## Who implements

Decide once, after the Kickoff approval, using `.claude/rules/moai/workflow/orchestration-mode-selection.md`: implement directly in the main session, hand the SPEC to one manager-develop, or split genuinely independent tracks across parallel manager-develop agents, each in its own worktree. Note the choice in one line in the Mode Selection section of `progress.md`.

## Doing the work

Work through the acceptance criteria with the methodology from `workflows/run/phase-execution.md`. When `tasks.md` exists it lists the planned tasks and files; update each task's status as it completes. Respect the `@MX` annotations in files you touch and add or update tags for new code per `.claude/rules/moai/workflow/mx-tag-protocol.md`.

The run phase owns exactly one status transition: on the first run-phase commit, set `status: in-progress` in `spec.md` (and the status line in `progress.md`, where it has one) and refresh `updated:`. `plan.md` and `acceptance.md` carry no `status:` field; refresh only their `updated:`. The later `implemented → completed` close belongs to manager-docs in `/moai sync`. Run does not rewrite SPEC body content (requirements, scope, acceptance criteria); a needed change goes back through manager-spec.

Keep a short divergence note as you go — files planned versus files actually touched, scope adjustments, new dependencies or directories. `/moai sync` uses it to update the SPEC and project documents.

## Delegating to manager-develop

Brief the agent once, completely, using `.claude/rules/moai/development/manager-develop-prompt-template.md`: SPEC path, the acceptance criteria and files it owns, `cycle_type` (`tdd` or `ddd`), constraints, how to verify, and what to return. Reference files by project-root-relative paths — a spawned agent's working directory is its own worktree root, and an absolute path into the main checkout or a `cd` prefix would make it edit the wrong tree.

manager-develop runs with `isolation: worktree`: it works in its own tree under `.claude/worktrees/`, commits there, and reports the branch. When it returns, read its report, merge the branch into the SPEC branch (`git merge --no-ff <branch>`, or cherry-pick its commits), resolve conflicts, and then run the verification below on the merged tree. Remove the worktree once its commits are merged (`git worktree remove <path>`). Don't repeat the agent's work; check it by running the project's commands. If it returned a blocker report, resolve the blocker (asking the user if the decision is theirs) and re-brief with the answer.

Parallel agents must own disjoint files. Merge them one at a time and verify after each merge, so a conflict or regression is attributable to one track.

# Verification

Verification happens in the main loop, by running the project's own commands — there is no separate verification agent; the independent audit is sync-auditor's job in `/moai sync`. Run what the repository uses: in a pnpm/turbo monorepo `pnpm turbo run test --affected` (or `pnpm --filter <pkg> test`) plus the typecheck and lint tasks, never a bare root runner that skips turbo's build ordering; in Go `go test ./...`, `go vet ./...` and the configured linter; in Python `pytest` and `ruff`. Independent checks can be issued as parallel Bash calls in one turn.

The run is verified when:

- each acceptance criterion in scope has been observed passing (a test, a command and its output, or a reproduced behavior);
- the full relevant test suite passes, with no previously passing test now failing;
- typecheck and lint meet `lsp_quality_gates.run` in `quality.yaml` (zero errors, zero type errors, zero lint errors, no regression from the starting baseline);
- coverage for changed code meets `test_coverage_target` where the project measures coverage.

A run that selected zero tests is not a pass — check that the runner actually executed the tests you meant (`.claude/rules/moai/development/verification-completeness.md`). For TDD work, keep the failing output of each new test from before its implementation; it is the evidence that the test was written first.

# progress.md

Fill the run-phase sections of `.moai/specs/SPEC-{ID}/progress.md` that plan scaffolded. Keep their headings exactly as spelled in `.claude/rules/moai/development/spec-frontmatter-schema.md`, section "progress.md Section Map" (including the lettered prefix): the `moai` CLI's SPEC era and closure checks match those literal headings.

- **Run-phase Evidence** — one table row per acceptance criterion, `| AC-ID | status | evidence |`, status `PASS`, `FAIL` or `PASS-WITH-DEBT`, evidence being the command and its observed output.
- **Run-phase Audit-Ready Signal** — a fenced `yaml` block with `run_complete_at`, `run_commit_sha` (a placeholder until the commit exists, then backfilled), `run_status`, `ac_pass_count`, `ac_fail_count`.
- The one-line Mode Selection note, and any bounded-fix-loop attempts in the Recursive Self-Diagnosis Log section.

# Git

Follow `.moai/config/sections/git-strategy.yaml` for the active mode.

- **Tier S/M without `--pr`:** commit on the current branch. If that branch is a shared or protected one the project integrates through review (`develop`, `release/*`, or `main` in a PR-based project), create a feature branch with the configured prefix first, or ask.
- **Tier L, or `--pr`:** manager-git creates the feature branch and opens the PR with the configured `merge_method`.

Use Conventional Commits in the `git_commit_messages` language, and add `Fixes #<n>` when `spec.md` has an `issue_number`. Stage specific paths, never unrelated files. Never use `--no-verify`, `--amend` on pushed commits, or force-push. Push only when `automation.auto_push` is true for the mode or the user approves it.

# Completion

Report the outcome first — criteria met, tests passing, coverage, commits — then anything deferred and why. Then ask the next step through AskUserQuestion: run `/moai sync` (recommended), review the changes, or stop. An explicit `/moai run` stops there; the default `/moai` pipeline continues into sync as `workflows/moai.md` describes.
