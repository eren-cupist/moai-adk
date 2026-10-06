---
paths: ".claude/skills/moai/workflows/plan.md,.claude/skills/moai/workflows/plan/**/*.md"
---

# Branch Origin Decision Protocol (BODP)

BODP decides one thing: which branch a new SPEC's branch starts from. It applies only when plan creates a branch — on `/moai plan --branch` or when the git-strategy config enables branch creation — and it has no command of its own.

The choice matters because a wrong base fails silently. Branching from a stale local `main` leaves out teammates' merged work; branching from `origin/main` when the user has unpushed local commits leaves out their own. Either surfaces later as a confusing merge.

## Hard Rules

- The BODP base question goes to the user through AskUserQuestion: recommended option first, at most 4 options, in the conversation language.
- Derive the recommendation from the signals below, not from habit. When no signal fires, recommend `origin/main`, which reflects the latest merged state rather than whatever the local checkout holds.
- Recommend local `main` only when `git log origin/main..main` shows unpushed local commits; otherwise it reintroduces the stale-base problem.

## Signals

| Signal | Detection |
|--------|-----------|
| A — Code dependency | the SPEC's `depends_on` names work on the current branch, or `git diff` overlaps the new SPEC's paths |
| B — Working tree co-location | `git status --porcelain` already lists `.moai/specs/<NewSpecID>/` |
| C — Open PR head | `gh pr list --head <currentBranch> --state open --json number` returns at least one entry (skip when `gh` is absent) |

## Decision Matrix

```
¬a ¬b ¬c → main      @ origin/main
 a ¬b ¬c → stacked   @ currentBranch
¬a  b ¬c → continue  @ ""
¬a ¬b  c → stacked   @ currentBranch
 a  b ¬c → continue  @ "" (b dominates)
 a ¬b  c → stacked   @ currentBranch
¬a  b  c → continue  @ "" (b dominates)
 a  b  c → continue  @ "" (b dominates)
```

Signal B dominates: when the SPEC's files are already in the working tree, the work has started on the current branch, and moving it would strand them. When Signal C fires, say in the option description that a stacked branch needs a rebase once its parent PR merges.

`main` creates the branch from the recommended base, `stacked` from the current branch, and `continue` stays on the current branch. The chosen base is passed to manager-git as `base=<branch>`. An "Other" answer is read as a branch name; an invalid one falls back to `origin/main` with a warning.
