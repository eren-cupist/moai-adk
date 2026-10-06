---
name: moai
description: >
  MoAI entry point for SPEC-driven development: routes /moai plan, run, sync,
  fix, and review — or a plain-language development request — to the matching
  workflow.
allowed-tools: Agent, AskUserQuestion, Skill, TaskCreate, TaskUpdate, TaskList, TaskGet, Bash, Read, Write, Edit, Glob, Grep
argument-hint: "[subcommand] [args] | \"natural language task\""
---

## Pre-execution Context

!`git status --porcelain 2>/dev/null || true`
!`git branch --show-current 2>/dev/null || true`

## Raw User Input

$ARGUMENTS

## Routing

Pick one workflow for the input above, read its file, and follow it. Everything after the subcommand word is the workflow's input, not a routing signal.

1. **Subcommand.** If the first word is one of these (or an alias), route there:

   | Subcommand | Aliases | Workflow file | Purpose |
   |------------|---------|---------------|---------|
   | `plan` | `spec` | `.claude/skills/moai/workflows/plan.md` | write and audit a SPEC, ending with the Implementation Kickoff Approval |
   | `run` | `impl` | `.claude/skills/moai/workflows/run.md` | implement a SPEC with DDD or TDD (`constitution.development_mode` in quality.yaml) |
   | `sync` | `docs`, `pr` | `.claude/skills/moai/workflows/sync.md` | sync documentation, close the SPEC, prepare the PR |
   | `fix` | | `.claude/skills/moai/workflows/fix.md` | find and fix lint, type, and build errors in one pass |
   | `review` | `code-review` | `.claude/skills/moai/workflows/review.md` | review code for correctness, security, and @MX tag compliance |

   Match on the first word only when the input is in Latin script or starts with the literal `/moai `. In Korean, Japanese, or Chinese text, an English word at the start ("plan", "fix", "run") is often an ordinary loanword, so classify the whole message instead (step 3).

2. **SPEC ID.** Otherwise, if the input contains a SPEC ID such as `SPEC-AUTH-001`, route to `run` for that SPEC.

3. **Intent.** Otherwise classify the whole message by meaning, in any language: planning or requirements → `plan`; reviewing code or changes, or a security audit → `review` (security scope for the latter); errors, failing builds, lint → `fix`; documentation, changelog, PR → `sync`; building or changing something with a clear scope → the default route.

4. **Default route.** A development request with no closer match goes to `.claude/skills/moai/workflows/moai.md`: plan → Implementation Kickoff Approval → run → sync. When the intent is genuinely ambiguous between two or three workflows, ask with AskUserQuestion, recommended choice first.

## Agents by workflow

- plan: manager-spec writes the SPEC, plan-auditor audits it, Explore for wide codebase investigation, manager-git for an optional branch or issue.
- run: manager-develop implements; manager-git for commits and branches.
- sync: manager-docs updates documentation, sync-auditor audits the result, manager-git handles pushes and PRs when they are requested or approved.
- fix: the orchestrator repairs directly, or briefs manager-develop for a larger repair.
- review: the orchestrator reviews in this session; only a large diff that splits into independent areas gets one read-only subagent per area.

Inject the domain skills listed in `.moai/config/sections/delegation.yaml` into a subagent's prompt when the work matches them (`At start, invoke Skill("<name>")`).

<!-- moai:contract-mode-start id="contract-signing-router" -->
Where `workflow.autonomy.mode: contract` — the Kickoff approval named here is the contract signature checked by `moai contract kickoff-check`. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "The signing gate".

<!-- moai:contract-mode-end -->

## Flags

| Flag | Applies to | Effect |
|------|-----------|--------|
| `--resume SPEC-<ID>` | plan, run, default | continue an existing SPEC from its first unfinished step |
| `--branch` | plan, default | create a branch for the SPEC at plan time |
| `--issue` | plan, default | create a GitHub issue for the SPEC |
| `--pr` | run, sync, default | use the PR route; with sync it is also the request to push and open the PR |
| `--auto-merge` | sync | merge the PR once its checks pass |

Workflow-specific flags are listed in each workflow file — review, for example, takes `--staged`, `--branch <base>`, `--commit <SHA>`, `--file <path>`, `--security`, `--lean`, `--deep` and `--patch`.

`--branch` (without a base, as the plan flag) with `run` or `sync` is an error: those phases use the branch plan created, and making a new one mid-lifecycle splits the SPEC's history. Say so in the user's language, show the right usage (`/moai plan "<description>" --branch`, then `/moai run SPEC-<ID>`), and stop.

`--worktree` is not supported: plan does not create a workspace. Tell the user to enter one first — `moai cc -w <name>`, or `moai cc -w <name> --spawn` to open it in a new session window — and run the command there.

## Running a workflow

Read the workflow file, load the `.moai/config/sections/*.yaml` files it names, and follow it. Track multi-step work with TaskCreate and TaskUpdate. Ask the user only where the workflow says to, or when a decision is genuinely theirs; when you do, use AskUserQuestion with the recommended option first and a short description of what each option implies (while `interview.recommendation_mode` is `pull`, give no option a recommended label). Implementation starts only after the user's Implementation Kickoff Approval: plan ends by asking it, and run asks it unless the user already approved in this conversation. Only contract mode replaces the question, with the signed contract. Answer the user in their conversation language.
