---
description: 보고된 버그, 실패한 테스트·CI 체크, 타입·린트 오류를 재현하고 최소 수정 후 검증
argument-hint: "[problem description | error | failing check] [--dry]"
allowed-tools: Skill
---

Dispatch the moai workflow `fix` with: $ARGUMENTS
- Harness with the Skill tool (Claude Code): invoke `Skill("moai")` with arguments: `fix` $ARGUMENTS
- Harness without a skill loader (Codex CLI): read `.agents/skills/moai/SKILL.md` (the mirrored dispatcher body) and follow its routing for the `fix` subcommand with the same arguments
