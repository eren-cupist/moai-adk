---
description: "Constitution zone registry — CONST-* clause records consumed by moai constitution CLI and zone audits"
paths: "**/zone-registry.md,**/.moai/config/sections/constitution.yaml"
---

# Zone Registry

The registry of constitution clauses that `moai constitution`, `moai doctor` and `moai spec lint` check. Each entry names a clause that must appear verbatim, on exactly one line, in its source file, under the heading its `anchor` points at. Entries pin the rules that must not drift silently: the SPEC format, approval before large changes, the question channel, concurrent writers, the CI auto-fix limits and secrets protection, and the plan and cleanup ordering of the SPEC workflow.

Entry fields: `id` (`CONST-V3R2-NNN`, `CONST-V3R5-NNN` or `CONST-V3R6-NNN`; IDs are never reused), `zone` (`Frozen` or `Evolvable`), `zone_class` (`frozen-canonical`, `frozen-safety`, `evolvable-tuning` or `evolvable-experimental`), `file`, `anchor`, `clause`, and `canary_gate` (`true` for Frozen entries).

## Retiring an Entry

A clause that is no longer in force is retired by prefixing its `clause` with a marker naming what replaced it, so the record of the withdrawn clause survives:

```text
clause: "[SUPERSEDED by <replacement>] <the original clause text>"
```

(This example fence is tagged `text` on purpose: the loader reads the first yaml-tagged fence in this file as the entry list.) `moai constitution validate` counts retired entries and skips their drift, canary and source-file checks; `--strict` checks them verbatim. The marker must be a prefix; `canary_gate: false` alone does not retire an entry.

## Usage

```bash
moai constitution list
moai constitution list --zone frozen
moai constitution list --file .claude/rules/moai/core/agent-common-protocol.md
moai constitution list --format json
```

## Entries

```yaml
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
  clause: "Code changes are held to TRUST 5"
  canary_gate: true

# --- CLAUDE.md ---
- id: CONST-V3R2-014
  zone: Evolvable
  zone_class: evolvable-tuning
  file: CLAUDE.md
  anchor: "#approach-and-approval"
  clause: "Before a large or ambiguous change, explain the approach and the files it touches, and get the user's approval."
  canary_gate: false

- id: CONST-V3R2-018
  zone: Frozen
  zone_class: frozen-canonical
  file: CLAUDE.md
  anchor: "#questions"
  clause: "use AskUserQuestion for questions that are genuinely the user's to answer"
  canary_gate: true

- id: CONST-V3R2-020
  zone: Evolvable
  zone_class: frozen-safety
  file: CLAUDE.md
  anchor: "#safety-rails"
  clause: "One writer per working tree: run write-capable agents in parallel only when each writes a different worktree"
  canary_gate: false

- id: CONST-V3R2-021
  zone: Evolvable
  zone_class: evolvable-experimental
  file: CLAUDE.md
  anchor: "#safety-rails"
  clause: "[SUPERSEDED by worktree-opt-in policy — see worktree-integration.md] Implementation teammates in team mode (role_profiles: implementer, tester, designer) MUST use isolation: worktree when spawned via Agent()"
  canary_gate: false

- id: CONST-V3R2-022
  zone: Evolvable
  zone_class: evolvable-experimental
  file: CLAUDE.md
  anchor: "#safety-rails"
  clause: "[SUPERSEDED by worktree-opt-in policy — see worktree-integration.md] Read-only teammates (role_profiles: researcher, analyst, reviewer) MUST NOT use isolation: worktree"
  canary_gate: false

- id: CONST-V3R2-023
  zone: Evolvable
  zone_class: evolvable-experimental
  file: CLAUDE.md
  anchor: "#safety-rails"
  clause: "[SUPERSEDED by worktree-opt-in policy — see worktree-integration.md] One-shot sub-agents making cross-file changes SHOULD use isolation: worktree"
  canary_gate: false

- id: CONST-V3R2-024
  zone: Evolvable
  zone_class: evolvable-experimental
  file: CLAUDE.md
  anchor: "#safety-rails"
  clause: "[SUPERSEDED by worktree-opt-in policy — see worktree-integration.md] GitHub workflow fixer agents MUST use isolation: worktree for branch isolation"
  canary_gate: false

# --- agent-common-protocol.md ---
- id: CONST-V3R2-036
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/core/agent-common-protocol.md
  anchor: "#user-interaction-boundary"
  clause: "Subagents cannot ask the user questions"
  canary_gate: true

# --- ci-autofix-protocol.md (10 entries: V3R5-004..013) ---
- id: CONST-V3R5-004
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#entry-condition"
  clause: "The CI auto-fix loop MUST be entered ONLY when the orchestrator hands off"
  canary_gate: true

- id: CONST-V3R5-005
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#iteration-cap"
  clause: "The auto-fix loop MUST attempt at most **3 iterations**"
  canary_gate: true

- id: CONST-V3R5-006
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#iteration-cap"
  clause: "The AskUserQuestion at iteration > 3 MUST be a blocking call"
  canary_gate: true

- id: CONST-V3R5-007
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#patch-commit-rule-no-force-push"
  clause: "Every auto-fix patch MUST be applied as a **new commit** on the PR branch"
  canary_gate: true

- id: CONST-V3R5-010
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#semantic-failure-no-auto-patch"
  clause: "The orchestrator MUST immediately escalate via AskUserQuestion with the diagnosis report."
  canary_gate: true

- id: CONST-V3R5-011
  zone: Frozen
  zone_class: frozen-safety
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#secrets-and-credentials-protection"
  clause: "The auto-fix loop MUST NOT modify `.env`, `.env.*`, credentials files"
  canary_gate: true

- id: CONST-V3R5-013
  zone: Frozen
  zone_class: frozen-canonical
  file: .claude/rules/moai/workflow/ci-autofix-protocol.md
  anchor: "#ci-infrastructure-preservation"
  clause: "The auto-fix loop MUST NOT modify CI watch infrastructure scripts or"
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
  anchor: "#4-anti-death-spiral-hook-carve-out-documentation-only-policy"
  clause: "Stop/PostToolUse hooks SHOULD exit 0 (allow the turn to end / the tool call to proceed) rather than exit 2 (block), so that recovery turns are NOT placed into the `error → stop-hook-blocks → retry → error` loop"
  canary_gate: true
```
