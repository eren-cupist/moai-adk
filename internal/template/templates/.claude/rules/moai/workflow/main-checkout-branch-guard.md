# Main-Checkout Branch Guard

The git rules for the shared primary checkout are in `AGENTS.md`, section "Git, branches, and the shared checkout": no branch switching, `git stash`, `git reset --hard` or rebase/merge there; stage by explicit pathspec; re-read `HEAD` and the current branch right before any commit or push. This file describes the hook that enforces them.

A PreToolUse hook (`internal/hook/branch_guard.go`) applies those rules only in the primary checkout, so worktree flows are unaffected. It is off by default because the hazard does not arise in single-session repositories; enable it with `workflow.branch_guard.enabled: true` in `.moai/config/sections/workflow.yaml`. When enabled:

- Branch-changing commands in the primary checkout are denied with a reason starting `BRANCH_GUARD_VIOLATION:`. Read-only forms (`git branch --list`, `-v`, `--show-current`, `--contains`, `git stash list`, `git merge-base`) pass.
- `git commit`, `git revert` and `git cherry-pick` are also denied in the primary checkout when the current branch is listed in `workflow.branch_guard.deny_commits_on`. A detached HEAD is allowed.
- The guard fails open: it denies only when it has positively matched a branch-state command in the primary checkout, and logs anything uncertain to `.moai/logs/branch-guard-audit.log`.
- An agent spawned under the name `manager-git` is exempt, and `MOAI_BRANCH_GUARD_EXEMPT=1` is read from the hook's own environment, so exporting it inside the guarded command does nothing.

When you hit `BRANCH_GUARD_VIOLATION`, the fix is a worktree, not a workaround: `moai worktree new <name>`, enter it with `moai cc -w <name>`, and drive it with `git -C <absolute-worktree-path>`.
