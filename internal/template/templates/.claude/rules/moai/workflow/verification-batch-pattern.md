---
description: "Batching independent read-only verification commands (loads at run/sync-phase completion)"
paths: "**/.moai/specs/**,**/.claude/skills/moai/workflows/run.md,**/.claude/skills/moai/workflows/sync/**,**/verification-batch-pattern.md,**/agent-common-protocol.md"
---

# Verification Batch Pattern

When a run or sync phase finishes, the verification commands — the change-scoped and full test runs, coverage, typecheck, lint, boundary greps, a CLI smoke check — are mostly read-only and independent of each other. Issue those as separate Bash calls in one turn so they run in parallel and each returns its own output block; the wait is then the slowest command rather than the sum.

Run commands in sequence when one depends on another or they share mutable state:

- a build whose output a later test invokes;
- two runs that write the same file (two coverage runs both writing `cover.out`);
- anything that changes the tree or the checkout (`git checkout`, fixture generation) next to commands that read it.

Chaining with `&&` is not batching — it runs in order and stops at the first failure — and backgrounding with `&` interleaves the outputs. Separate tool calls in the same turn are the batch.

Each command's result is evidence only for the tree it ran on; when a claim cites a result recorded earlier (for example a verification snapshot), the snapshot key must still match the current tree, or the command is run again (`.claude/rules/moai/core/verification-claim-integrity.md`).
