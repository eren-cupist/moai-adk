---
description: DDD/TDD 방법론으로 SPEC 요구사항 구현
argument-hint: "SPEC-XXX [--pr] [--skip-audit] [--resume]"
allowed-tools: Skill
---

Dispatch the moai workflow `run` with: $ARGUMENTS
- Harness with the Skill tool (Claude Code): invoke `Skill("moai")` with arguments: `run` $ARGUMENTS
- Harness without a skill loader (Codex CLI): read `.agents/skills/moai/SKILL.md` (the mirrored dispatcher body) and follow its routing for the `run` subcommand with the same arguments
