---
description: "Constitution zone registry — CONST-* clause records consumed by moai constitution CLI and zone audits"
paths: "**/zone-registry.md,**/.moai/config/sections/constitution.yaml"
---

# Zone Registry

Single source of truth enumerating every HARD clause in the MoAI-ADK rules tree.
Each entry carries a unique ID, Zone classification, source file, anchor, verbatim clause, and canary_gate field.

## HISTORY

| Version | Date | Author | Description |
|---------|------|--------|-------------|
| 1.0.0 | (initial) | maintainer | Initial creation — annotation pass over 4 load-bearing source files |
| 1.1.0 | (later)   | maintainer | Coverage gap closure — CONST-V3R5-001..041 added (parallel namespace), zone_class 4-classification introduced (retroactive on all 115 entries) |

## ID Allocation Policy

ID format: `CONST-V3R2-NNN` (initial namespace) or `CONST-V3R5-NNN` (parallel namespace)

Allocation rules:
- Fixed file order: `CLAUDE.md` → `.claude/rules/moai/core/moai-constitution.md` → `.claude/rules/moai/core/agent-common-protocol.md`
- Within each file, assign IDs in ascending `(anchor_line_number)` order
- 001-050: pre-existing clauses (HARD clauses found in the 3 files above)
- 150+: reserved for future additions (V3R2 namespace)

V3R5 namespace policy:
- New entries use the parallel namespace starting at `CONST-V3R5-001`
- The 3 internal V3R2 gaps (047/048/050) are NOT filled — preserved as historical record
- `zone_class` field (4-enum): `frozen-canonical` | `frozen-safety` | `evolvable-tuning` | `evolvable-experimental`

CanaryGate defaults (plan.md §7 OQ6 decision):
- Frozen → `canary_gate: true`
- Evolvable → `canary_gate: false`

## Retiring an Entry

A clause that is no longer in force is **retired, not deleted** — deleting it destroys the record that the clause once existed and was withdrawn.

To retire an entry, prefix its `clause` with a `[SUPERSEDED …]` marker naming what replaced it:

```text
clause: "[SUPERSEDED by <replacement>] <the original clause text>"
```

(The fence above is deliberately tagged `text`, not `yaml`. The registry loader reads the **first** yaml-tagged fence in this file as the entry list, so a yaml-tagged example placed before `## Entries` would be parsed as the registry and fail to load — and so would that fence marker written out literally in prose.)

`moai constitution validate` then counts the entry as retired and skips its drift, canary-gate, and source-file checks — the source text is gone by definition, so those checks could only ever fail. The entry still appears in `moai constitution list`, and the retired total is reported (`retired_count` in JSON output).

Two boundaries:

- The marker is a **prefix**, not a substring. A live clause that merely mentions `[SUPERSEDED …]` in its own text stays fully checked.
- `canary_gate: false` is **not** a retirement marker — it is the documented default for every Evolvable entry. A retired Frozen entry may carry `canary_gate: false` because it is retired, not the other way round.

`moai constitution validate --strict` ignores the marker and checks retired entries verbatim, for auditing what they still hold.

## Usage Guide

```bash
# List the entire registry
moai constitution list

# Filter by Frozen zone
moai constitution list --zone frozen

# List clauses from a specific file only
moai constitution list --file .claude/rules/moai/core/moai-constitution.md

# JSON-format output
moai constitution list --format json
```

## Entries

```yaml
# ============================================================
# 001-010: CLAUDE.md HARD clauses (§1 Hard Rules)
# ============================================================
- id: CONST-V3R2-001
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/spec-workflow.md
  anchor: "#plan-phase"
  clause: "Create comprehensive specification using EARS format."
  canary_gate: true

- id: CONST-V3R2-002
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/moai-constitution.md
  anchor: "#quality-gates"
  clause: "All code changes must pass TRUST 5 validation"
  canary_gate: true

- id: CONST-V3R2-006
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/agent-common-protocol.md
  anchor: "#user-interaction-boundary"
  clause: "`AskUserQuestion` is the **only** user-facing question channel"
  canary_gate: true

# ============================================================
# 008-020: CLAUDE.md HARD clauses (§1 Hard Rules — orchestrator behavior)
# ============================================================
- id: CONST-V3R2-014
  zone: Evolvable
  zone_class: evolvable-tuning
  file: CLAUDE.md
  anchor: "#7-safe-development-protocol"
  clause: "Before non-trivial code, explain the approach + which files change + why; get user approval"
  canary_gate: false

- id: CONST-V3R2-018
  zone: Frozen
  zone_class: frozen-canonical
  file: CLAUDE.md
  anchor: "#8-user-interaction-architecture"
  clause: "Every question directed at the user MUST be asked via AskUserQuestion. Free-form prose questions in response text are prohibited."
  canary_gate: true

# ============================================================
# 020-030: CLAUDE.md §14 Worktree Isolation Rules + §11 Background Agent
# ============================================================
- id: CONST-V3R2-020
  zone: Evolvable
  zone_class: frozen-safety
  file: CLAUDE.md
  anchor: "#14-parallel-execution-safeguards"
  clause: "subagents run in the background by default (the runtime chooses foreground only when it needs the result; every permission prompt still surfaces in the main session); MoAI does not set `background:` — the retained safeguard is concurrency, not backgrounding"
  canary_gate: false

- id: CONST-V3R2-021
  zone: Evolvable
  zone_class: evolvable-experimental
  file: CLAUDE.md
  anchor: "#14-parallel-execution-safeguards"
  clause: "[SUPERSEDED by worktree-opt-in policy — see CLAUDE.md §14 + worktree-integration.md § Terminology Glossary] Implementation teammates in team mode (role_profiles: implementer, tester, designer) MUST use isolation: worktree when spawned via Agent()"
  canary_gate: false

- id: CONST-V3R2-022
  zone: Evolvable
  zone_class: evolvable-experimental
  file: CLAUDE.md
  anchor: "#14-parallel-execution-safeguards"
  clause: "[SUPERSEDED by worktree-opt-in policy — see CLAUDE.md §14 + worktree-integration.md § Terminology Glossary] Read-only teammates (role_profiles: researcher, analyst, reviewer) MUST NOT use isolation: worktree"
  canary_gate: false

- id: CONST-V3R2-023
  zone: Evolvable
  zone_class: evolvable-experimental
  file: CLAUDE.md
  anchor: "#14-parallel-execution-safeguards"
  clause: "[SUPERSEDED by worktree-opt-in policy — see CLAUDE.md §14 + worktree-integration.md § Terminology Glossary] One-shot sub-agents making cross-file changes SHOULD use isolation: worktree"
  canary_gate: false

- id: CONST-V3R2-024
  zone: Evolvable
  zone_class: evolvable-experimental
  file: CLAUDE.md
  anchor: "#14-parallel-execution-safeguards"
  clause: "[SUPERSEDED by worktree-opt-in policy — see CLAUDE.md §14 + worktree-integration.md § Terminology Glossary] GitHub workflow fixer agents MUST use isolation: worktree for branch isolation"
  canary_gate: false

# ============================================================
# 025-035: moai-constitution.md HARD clauses
# ============================================================
- id: CONST-V3R2-025
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/moai-constitution.md
  anchor: "#moai-orchestrator"
  clause: "AskUserQuestion is the sole user-facing question channel"
  canary_gate: true

- id: CONST-V3R2-026
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/moai-constitution.md
  anchor: "#moai-orchestrator"
  clause: "used ONLY by the MoAI orchestrator (subagents must never prompt users)"
  canary_gate: true

# ============================================================
# 036-045: agent-common-protocol.md HARD clauses
# ============================================================
- id: CONST-V3R2-036
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/agent-common-protocol.md
  anchor: "#user-interaction-boundary"
  clause: "Subagents MUST NOT prompt the user. AskUserQuestion is reserved exclusively for the MoAI orchestrator."
  canary_gate: true

- id: CONST-V3R2-038
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/agent-common-protocol.md
  anchor: "#user-interaction-boundary"
  clause: "AskUserQuestion is reserved exclusively for the MoAI orchestrator"
  canary_gate: true

- id: CONST-V3R2-044
  zone: Evolvable
  zone_class: frozen-safety
  file: .claude/rules/moai/core/agent-common-protocol.md
  anchor: "#background-agent-execution"
  clause: "The retained safeguard is **concurrency, not backgrounding**"
  canary_gate: false

# ============================================================
# CONST-V3R5-001..041: new parallel namespace
# Completes coverage of unmapped [HARD] rules — 11 source files newly registered
# ============================================================
- id: CONST-V3R5-001
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/askuser-protocol.md
  anchor: "#orchestratorsubagent-boundary"
  clause: "Subagents MUST NOT invoke `AskUserQuestion`"
  canary_gate: true

- id: CONST-V3R5-002
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/askuser-protocol.md
  anchor: "#orchestratorsubagent-boundary"
  clause: "Subagents MUST NOT output free-form prose questions directed at the user"
  canary_gate: true

- id: CONST-V3R5-003
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/askuser-protocol.md
  anchor: "#orchestratorsubagent-boundary"
  clause: "Subagents MUST NOT embed AskUserQuestion call syntax in their response body"
  canary_gate: true

# --- ci-autofix-protocol.md (10 entries: V3R5-004..013) ---
- id: CONST-V3R5-004
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#entry-condition"
  clause: "Enter the auto-fix loop only when the orchestrator hands off a failing required check"
  canary_gate: true

- id: CONST-V3R5-005
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#iteration-cap"
  clause: "Attempt at most three iterations per pull request"
  canary_gate: true

- id: CONST-V3R5-006
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#iteration-cap"
  clause: "After the third failed iteration, ask the user through a blocking AskUserQuestion with no timeout"
  canary_gate: true

- id: CONST-V3R5-007
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#patch-commit-rule-no-force-push"
  clause: "Apply every auto-fix patch as a new commit on the PR branch"
  canary_gate: true

- id: CONST-V3R5-010
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#semantic-failure-no-auto-patch"
  clause: "For a semantic failure, do not patch: escalate immediately through AskUserQuestion with the diagnosis"
  canary_gate: true

- id: CONST-V3R5-011
  zone: Frozen
  zone_class: frozen-safety
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#secrets-and-credentials-protection"
  clause: "The auto-fix loop never modifies `.env`, `.env.*`, credentials files"
  canary_gate: true

- id: CONST-V3R5-013
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#ci-infrastructure-preservation"
  clause: "The auto-fix loop never modifies CI watch infrastructure scripts or workflow definitions"
  canary_gate: true

# --- spec-workflow.md (2 new entries: V3R5-027..028; CONST-V3R2-001 covers the third) ---
- id: CONST-V3R5-027
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/spec-workflow.md
  anchor: "#spec-phase-discipline"
  clause: "Step 1 (plan) MUST execute in main checkout on BOTH routes. NO L2/L3 worktree at this step"
  canary_gate: true

- id: CONST-V3R5-028
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/spec-workflow.md
  anchor: "#spec-phase-discipline"
  clause: "Step 4 (cleanup) applies to **Route B only**. It MUST happen ONLY after BOTH run AND sync PRs are merged"
  canary_gate: true

# --- branch-origin-protocol.md (1 entry: V3R5-035) ---

- id: CONST-V3R5-035
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/development/branch-origin-protocol.md
  anchor: "#hard-rules"
  clause: "Skill body BODP gate MUST follow the askuser-protocol Socratic structure: `(권장)` first, ≤4 options, conversation_language match"
  canary_gate: true

# ============================================================
# CONST-V3R6-NNN: V3R6 modern-era parallel namespace
# (first V3R6 entry: a runtime-recovery predecessor SPEC, M3)
# ============================================================
# --- runtime-recovery-doctrine.md (1 entry: V3R6-001 anti-death-spiral) ---
- id: CONST-V3R6-001
  zone: Evolvable
  zone_class: frozen-safety
  file: .claude/rules/moai/workflow/runtime-recovery-doctrine.md
  anchor: "#anti-death-spiral-hook-carve-out"
  clause: "Stop and PostToolUse hooks should exit 0 on a recovery turn (letting the turn end or the tool call proceed) rather than exit 2"
  canary_gate: true
```
