---
description: Reproduce and minimally fix a reported bug, failing test or CI check, or type and lint errors, then verify
argument-hint: "[problem description | error | failing check] [--dry]"
allowed-tools: Skill
---

Dispatch the moai workflow `fix` with: $ARGUMENTS
- Harness with the Skill tool (Claude Code): invoke `Skill("moai")` with arguments: `fix` $ARGUMENTS
- Harness without a skill loader (Codex CLI): read `.agents/skills/moai/SKILL.md` (the mirrored dispatcher body) and follow its routing for the `fix` subcommand with the same arguments
