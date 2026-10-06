---
description: "How the orchestrator decides who implements a run: directly, one manager-develop, or parallel manager-develop agents in separate worktrees"
paths: ".moai/specs/**,.claude/skills/moai/workflows/run.md,.claude/skills/moai/workflows/plan.md,.claude/rules/moai/workflow/spec-workflow.md"
metadata:
  version: "2.0.0"
  status: "active"
  tags: "orchestration, delegation, run, worktree"
---

# Orchestration Choice for a Run

After the Implementation Kickoff Approval, the orchestrator decides how the SPEC gets implemented. The decision is the orchestrator's own judgment call; it is not a question for the user. Nothing about this choice relaxes the Kickoff gate: every option below starts only after that approval.

<!-- moai:contract-mode-start id="contract-signing" -->
Where `workflow.autonomy.mode: contract` — a human signature of the SPEC's contract (`signer_kind: human`, `method: interactive-tty`) for which `moai contract kickoff-check` exits 0 is the Implementation Kickoff Approval this paragraph requires; no other signature is equivalent by this block. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Equivalence clause (human signature only)".

<!-- moai:contract-mode-end -->
## The options

Applied to a run:

- **Directly in the main session** — the default for small or tightly coupled SPECs: a handful of files, one area of the code, or changes that depend on each other so closely that splitting them would mean constant coordination. Most Tier S and many Tier M SPECs fit here.
- **One manager-develop** — when the SPEC is substantial but forms a single track, and isolation is worth having: the agent works in its own worktree, so the main checkout stays clean until its branch is merged, and the main session's context stays free for verification and review.
- **Parallel manager-develop agents** — only when the SPEC splits into sizeable tracks that are genuinely independent: disjoint files, no shared types or schema being changed mid-flight, and each track verifiable on its own (for example a backend endpoint and an unrelated frontend screen that meet only at an already-agreed contract). Each agent gets its own worktree and a full brief; the orchestrator merges the branches one at a time and verifies after each merge.

When unsure between two options, take the one with fewer agents: coding work rarely parallelizes as well as it seems, and every extra agent adds a merge and a reconciliation step.

Read-only research before implementation (surveying unfamiliar code, checking several independent areas) can fan out to a few concurrent read-only agents when the areas really are independent; they write nothing and need no worktree.

## Constraints

- One writer per working tree: two write-capable agents never share a tree. Parallel writers each get their own worktree (manager-develop declares `isolation: worktree`).
- Subagents cannot ask the user anything. Every decision they would need is settled before they start, and a subagent that hits a decision returns a blocker report; the orchestrator asks the user and re-briefs.
- Agents spawned by the orchestrator do not spawn further agents for implementation; the orchestrator owns the split and the merge.

## Record

Note the choice and a one-sentence reason in the Mode Selection section of the SPEC's `progress.md` (the section letter is assigned in `.claude/rules/moai/development/spec-frontmatter-schema.md`, section "progress.md Section Map"), for example `Mode Selection: direct — 4 files in one package`. If the scope changes mid-run (a blocker reveals more work), revisit the choice.
