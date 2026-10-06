---
name: moai-workflow-spec
description: >
  SPEC format reference: GEARS requirement patterns, acceptance-criteria shape,
  SPEC directory layout, and the line formats the SPEC linter reads. Use when
  writing, revising, or auditing SPEC documents under .moai/specs/; not for
  implementing a SPEC.

when_to_use: >
  Writing or reviewing spec.md, plan.md, acceptance.md, design.md, research.md,
  or progress.md for a MoAI SPEC.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Write, Edit, Bash(git:*), Bash(ls:*), Bash(wc:*), Bash(mkdir:*), Grep, Glob
user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-01-08"
  modularized: "true"
  tags: "workflow, spec, gears, ears, requirements, acceptance-criteria"
  author: "MoAI-ADK Team"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000
---

# SPEC Format Reference

A SPEC is the contract between planning and implementation: it says what will be built and why, and how to check that it was built. It states observable behavior and constraints, not function names, class layouts, or API schemas — those are decided in run. `moai spec lint` checks the shapes below, so follow them exactly where they say the linter reads them. A worked example is in `references/examples.md`.

## Directory layout

One directory per SPEC: `.moai/specs/SPEC-<DOMAIN>-<NNN>/` — never a flat `.md` file in `.moai/specs/`. The artifact set depends on the tier; read it from `.claude/rules/moai/workflow/spec-workflow.md`, section "SPEC Complexity Tier (S/M/L)". A directory is complete when it holds exactly what its tier names: a Tier S directory with only `spec.md` and `plan.md` is complete, and adding `acceptance.md` there would duplicate the criteria already inline in `spec.md`.

| File | Tier | Holds |
|------|------|-------|
| `spec.md` | all | frontmatter, HISTORY, context, scope, requirements, Out of Scope; at Tier S also the acceptance criteria |
| `plan.md` | all | approach, milestones, risks, assumptions, open questions |
| `acceptance.md` | M, L | acceptance criteria, edge cases, Definition of Done |
| `design.md` | L | architecture and design decisions with their trade-offs |
| `research.md` | L | codebase findings with file paths, existing patterns, constraints |
| `progress.md` | all | the `§E` phase signals (not counted in the tier) |
| `decision-index.md` | when `interview.decision_gate: on` | unresolved decisions for the user (format in manager-spec) |

Frontmatter fields, the status enum, and status ownership are in `.claude/rules/moai/development/spec-frontmatter-schema.md`. Only `spec.md` carries `status:`.

## GEARS requirement patterns

Each requirement is one sentence in one of these patterns. `<subject>` can be any noun — the system, a component, a service, an agent, a command.

| Pattern | Form |
|---------|------|
| Ubiquitous | `The <subject> shall <behavior>.` |
| Event-driven | `When <event>, the <subject> shall <behavior>.` |
| State-driven | `While <state>, the <subject> shall <behavior>.` |
| Capability gate | `Where <capability, feature flag, or config>, the <subject> shall <behavior>.` |
| Unwanted event | `When <undesired condition is detected>, the <subject> shall <response>.` |

The modifiers chain into one compound clause, in this order, with any subset: `Where <precondition>, while <state>, when <event>, the <subject> shall <behavior>.`

A prohibition is written `The <subject> shall not <action>.`

Legacy EARS `IF <condition> THEN <action>` is replaced by the unwanted-event `When` form; the linter warns on it with `LegacyEARSKeyword` (an error under `--strict`). Older SPECs written in EARS remain valid.

Write `shall` for every obligation. "should", "may", "appropriate", "reasonable", and "as needed" leave the requirement open to interpretation; if the behavior is truly optional, make it a capability gate.

## Requirement lines (what the linter reads)

Write each requirement as a list item that starts with its ID and a colon:

```markdown
- REQ-<DOMAIN>-<NNN>: When <event>, the <subject> shall <behavior>.
```

- ID pattern: `REQ-<DOMAIN>[-<DOMAIN>]-<NNN>[-<NNN>]` — each domain segment starts with an uppercase letter, each number is exactly three digits. Anything else is `InvalidREQID`; a repeated ID is `DuplicateREQID`.
- The sentence must contain `SHALL` / `shall` as a word. A requirement that opens with When/While/Where/If/The but has no `shall` is `ModalityMalformed`; one with neither is reported `ModalityUnjudged`. A Korean requirement is judged when it carries `(SHALL)` after `해야 한다`.
- `### REQ-<DOMAIN>-<NNN>: <text>` headings are also read. Requirements in a table are not collected unless each row carries a modality marker (`REQTableRowsRejected`), so prefer list items.
- Number requirements consecutively without gaps.

## Acceptance criteria

Requirements say what must be true; acceptance criteria are the checks that show it is. Criteria use Given/When/Then — that is the correct format for this layer, not a GEARS violation.

```markdown
- AC-<DOMAIN>-<NNN>: Given <precondition>, When <action>, Then <observable result>. (maps REQ-<DOMAIN>-<NNN>)
```

- Every requirement is covered by at least one criterion. Coverage is read from the `(maps REQ-…)` text on criterion lines in `spec.md` and the sibling `acceptance.md`; an uncovered requirement is reported as `CoverageIncomplete`. Several requirements can be listed, comma-separated after `maps`.
- Criterion IDs are unique across the SPEC (`DuplicateAcceptanceID`).
- Each criterion is binary: a tester can decide pass or fail without judgment. Name the command, output, file, status code, or measured threshold that decides it.
- Cover failure and edge cases, not only the happy path.
- At Tier S the criteria go in a section of `spec.md`; at Tier M and L they go in `acceptance.md`, together with the Definition of Done.

## spec.md body

In order: HISTORY (a version table, the first row being the initial draft), context and motivation, scope, requirements, the acceptance criteria at Tier S, and the exclusions section.

The exclusions section is required (`MissingExclusions`): it must contain the words "Out of Scope", at least one H3 sub-heading of the form `### Out of Scope — <topic>`, and at least one `-` bullet under it. A bare H2 with no `### Out of Scope —` sub-heading fails.

```markdown
## Exclusions (What NOT to Build)

### Out of Scope — Social login
- OAuth providers other than the existing email/password flow.
```

For changes to existing code, requirements may be grouped with delta markers so run knows what needs characterization tests:

```markdown
### [DELTA] <Module>
- [EXISTING] unchanged context
- [MODIFY] existing behavior that changes — characterize before changing
- [NEW] new behavior
- [REMOVE] behavior to delete — check dependents first
```

## plan.md body

The technical approach; milestones `M1`, `M2`, … in priority order (no time estimates), each naming the acceptance criteria it completes with an `Exit:` line (e.g. `Exit: AC-<DOMAIN>-001, AC-<DOMAIN>-002 green`); risks and how they are handled; assumptions made instead of asking; and references to existing code that the implementation should follow (`path/to/file:line`). Lead with the decisions most likely to change in review — data model, interfaces, user-facing flow — and put mechanical steps last. When the SPEC belongs to a multi-SPEC group, name the Epic (`.claude/rules/moai/development/sprint-round-naming.md`).

## Open questions

`[NEEDS CLARIFICATION: <topic>]` marks a question that has to be answered before implementation starts. Put it only in `plan.md` or `research.md`, never in `spec.md` or `acceptance.md`, followed by the question and the constraint that makes it matter:

```markdown
[NEEDS CLARIFICATION: session store]
PostgreSQL or Redis for sessions? The read path must stay under 100 ms at 10k concurrent users.
```

plan-auditor fails a SPEC that still contains one, and the Implementation Kickoff Approval waits until each is resolved. A `TODO` or `@MX:TODO` in code is something else: implementation debt that needs no answer from the user.

## progress.md

Create it at plan time with the four section headings, in this order, each with a one-line placeholder:

```markdown
## §E.1 Plan-phase Audit-Ready Signal
## §E.2 Run-phase Evidence
## §E.3 Run-phase Audit-Ready Signal
## §E.4 Sync-phase Audit-Ready Signal
```

The plan workflow fills only §E.1 (`plan_complete_at`, `plan_status: audit-ready`); run fills §E.2 and §E.3, sync fills §E.4. The headings are matched literally by the era classifier — see "progress.md Section Map" in the schema file.

## What is not a SPEC

`.moai/specs/` holds only work to be built. An analysis of existing code (security audit, performance report, dependency review) goes to `.moai/reports/<type>-<date>/`; documentation of how something works or how to use it goes to `.moai/docs/`.
