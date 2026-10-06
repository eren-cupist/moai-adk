---
description: "Four Karpathy coding principles as short checkpoints for code changes"
paths: "**/*.go,**/*.py,**/*.ts,**/*.js,**/*.java,**/*.rs,**/*.c,**/*.cpp,**/*.rb,**/*.php,**/*.kt,**/*.swift,**/*.dart,**/*.ex,**/*.scala,**/*.hs,**/*.zig"
---

# Karpathy Coding Principles

Four checkpoints for any code change, adapted from Andrej Karpathy's coding principles (attribution in `.claude/rules/moai/NOTICE.md`):

- **Think before coding** — state the assumptions the change rests on; if two readings of the request lead to different code, resolve which one is meant.
- **Simplicity first** — prefer existing code, the standard library, a platform feature or an installed dependency over new code, in that order (the reuse ladder in `.claude/rules/moai/core/moai-constitution.md`, Agent Core Behaviors, Enforce Simplicity). Add an abstraction only when it removes real duplication or complexity now.
- **Surgical changes** — touch what the task needs, match the surrounding style, and leave adjacent refactors for their own change.
- **Goal-driven execution** — know the check that proves the change works, run it, and report what it showed.

Wrong/right examples: `.claude/skills/moai/references/anti-patterns.md`.
