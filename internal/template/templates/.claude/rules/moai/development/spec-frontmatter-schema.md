---
description: "SPEC frontmatter schema, status lifecycle, and progress.md section map — the reference the SPEC linter enforces"
paths: "**/.moai/specs/**,internal/spec/**"
---

# SPEC Frontmatter Schema

This is the reference for SPEC frontmatter, status ownership, and the `progress.md` layout. `internal/spec/lint.go` (`FrontmatterSchemaRule` and the rules around it) enforces it through `moai spec lint`; when this file and the linter disagree, the linter is what runs.

## Canonical 12 Required Fields

Every `spec.md` carries exactly these 12 fields. A missing field, an empty value, or a snake_case alias produces a `FrontmatterInvalid` finding.

```yaml
---
id: SPEC-{DOMAIN}-{NUM}
title: "Human-readable title"
version: "X.Y.Z"
status: draft
created: YYYY-MM-DD
updated: YYYY-MM-DD
author: Author Name
priority: P1
phase: "vX.Y.Z target"
module: "path/to/module"
lifecycle: spec-anchored
tags: "tag1, tag2, tag3"
---
```

### Artifact Statelessness

The 12-field obligation binds `spec.md` only. The sibling artifacts `plan.md`, `acceptance.md`, `design.md`, and `research.md` are stateless on the status axis: they carry no `status:` field (`ArtifactStatusFieldForbidden`), because a SPEC's lifecycle state lives in exactly one place — `spec.md` — and every lint, audit, and close path reads it from there. Other frontmatter fields in those files are allowed, and so is no frontmatter at all. This holds at every tier. `progress.md` records progress in body sections, not frontmatter.

## Field Reference

| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| `id` | string | `^SPEC(-[A-Z][A-Z0-9]*)+-[0-9]{3}$` | e.g. `SPEC-AUTH-001`; no letter suffix after the number |
| `title` | string | non-empty, quoted | |
| `version` | string | semver `X.Y.Z`, quoted | start at `"0.1.0"` |
| `status` | enum | see Status Enum | lifecycle state |
| `created` | date | `YYYY-MM-DD` | |
| `updated` | date | `YYYY-MM-DD` | |
| `author` | string | non-empty | |
| `priority` | enum | `P0`, `P1`, `P2`, `P3`, or `High`, `Medium`, `Low`, `Critical` | default `P1` |
| `phase` | string | a release or milestone target label | see Prohibited phase values |
| `module` | string | non-empty, path-like | affected module or directory |
| `lifecycle` | enum | `spec-anchored`, `spec-lite`, `exploratory` | default `spec-anchored` |
| `tags` | string | comma-separated, non-empty | not a YAML list |

### Prohibited phase values

`phase` names the release or milestone the SPEC targets, and it does not change as the SPEC moves through the workflow — `status` carries the stage. The lifecycle words `plan`, `run`, `sync`, and `mx` are rejected as the whole value (trimmed, case-insensitive) with `FrontmatterPhaseInvalid` at error severity, on every SPEC including older ones. A longer label that merely contains one of these words, such as `"v3.0.0 — Runtime Hardening"`, is valid. When the target release is unknown, use the next unreleased version. `internal/spec/era.go` also reads `phase` for era classification, so a stage word there quietly weakens that check.

## Status Enum (8 values)

Valid values: `draft`, `planned`, `in-progress`, `implemented`, `completed`, `superseded`, `archived`, `rejected`.

The active flow is `draft → in-progress → implemented → completed`, with `superseded`, `archived`, or `rejected` as exits. `planned` stays in the enum so older SPECs that recorded it still parse, but no agent writes it.

A `completed` SPEC may go back to `in-progress` for an in-place amendment. Declare it with the optional `amendment_of:` field (the SPEC's own ID for an in-place amendment, the parent ID for a successor) and an `## Amendments` section in the HISTORY that records the prior completed version, `prior_completed_sha`, rationale, and scope. `moai spec audit` (`SyncStatusDrift`) exempts the amendment from drift only when `status: in-progress`, `amendment_of:` is set, the `## Amendments` (or `### Amendments`) section exists, and its `prior_completed_sha` equals the prior close's `sync_commit_sha`.

## Status Transition Ownership Matrix

| Transition | Owning agent | Commit subject pattern |
|------------|--------------|------------------------|
| `(none) → draft` | manager-spec | `feat(SPEC-{ID}): plan-phase artifacts ({tier}, {N} artifacts)` — N is the tier's artifact count |
| `draft → in-progress` | manager-develop, on the first run-phase commit | `feat(SPEC-{ID}): M1 ...` or `fix(SPEC-{ID}): M1 ...` |
| `in-progress → implemented → completed` | manager-docs, in the single sync commit | `docs(SPEC-{ID}): sync-phase artifacts + 3-phase close` (or `chore(SPEC-{ID}): …`); the `3-phase close` infix is what marks the completed transition |
| `* → superseded` | manager-spec, when writing the superseding SPEC | `feat(SPEC-{NEW-ID}): supersedes SPEC-{OLD-ID}` |
| `* → archived` | manager-docs | `chore(specs): archive SPEC-{ID}` |
| `* → rejected` | orchestrator decision, recorded by manager-docs | `chore(SPEC-{ID}): rejected per <rationale>` |
| `completed → in-progress (amendment)` | manager-spec, on orchestrator re-delegation | `feat(SPEC-{ID}): in-place amendment <rationale-summary>` |

The lifecycle has three phases; MX tag validation is part of sync, and the `completed` transition rides the sync commit rather than a separate commit. The matrix covers `status:` transitions only — repairing a mis-written field is covered under Non-transition frontmatter corrections.

## progress.md Section Map

`internal/spec/era.go` (`ClassifyEra`) matches the literal headings `§E.2`, `§E.3`, `§E.4`, `§E.5` and the field names `sync_commit_sha` and `mx_commit_sha` in `progress.md`. Renaming any of them misclassifies the SPEC's era, so keep them exactly.

| Section heading | Purpose | Parsed by era.go | Written by |
|-----------------|---------|------------------|------------|
| `## §E.1 Plan-phase Audit-Ready Signal` | `plan_complete_at`, `plan_status: audit-ready` | no | manager-spec / plan workflow |
| `## §E.2 Run-phase Evidence` | run evidence; the heading marks its start | yes | manager-develop |
| `## §E.3 Run-phase Audit-Ready Signal` | run completion signal | yes | manager-develop |
| `## §E.4 Sync-phase Audit-Ready Signal` | sync close, including `sync_commit_sha:` | yes, heading and field | manager-docs |
| `## §E.5 Mx-phase` | retired; still matched so legacy SPECs classify | yes, legacy only | nobody — do not write new ones |
| `## §F Phase 4 Mode Selection` | orchestrator's mode-selection log | no | orchestrator, before the first run-phase spawn |
| `## §H Recursive Self-Diagnosis Log` | run-phase diagnose/patch/verify record | no | manager-develop / orchestrator |
| `## §I Token Accounting` | per-SPEC token spend at sync close | no | manager-docs |

A new progress.md concern takes a fresh top-level letter; it never reuses the `§E` namespace, which belongs to the parsed lifecycle structure. (The `§F` here is a progress.md section, unrelated to any section letter inside spec.md or plan.md.)

### Close-subject full-ID mandate

Every close commit — the sync commit carrying the `implemented → completed` transition — names exactly one full SPEC-ID in its subject scope, e.g. `chore(SPEC-{DOMAIN}-{SUB}-001): … 3-phase close`. A combined or abbreviated scope that names only a shared prefix is prohibited, because the drift detector extracts exact SPEC-ID tokens and cannot map a prefix to its sibling SPECs. To close N sibling SPECs, make N close commits. The close-infix matcher (`internal/spec/transitions.go`) accepts both `3-phase close` and the legacy `4-phase close` found in older history.

### Forbidden ownership crossings

- manager-docs does not edit the body of `spec.md`, `plan.md`, or `acceptance.md`. It may update `status:` and `updated:` in `spec.md` across the sync transition, and refresh `updated:` in the stateless `plan.md` / `acceptance.md`. When sync shows the SPEC body needs a change, it returns a blocker report and the orchestrator re-delegates to manager-spec.
- manager-develop does not edit SPEC bodies either. It may set `status:` and `updated:` in `spec.md` on `draft → in-progress`, and refresh `updated:` in the stateless `plan.md` / `acceptance.md`. A needed body change goes back through the orchestrator to manager-spec.

SHA placeholder backfill exemption: the `sync_commit_sha` / `mx_commit_sha` fields in `progress.md` §E.3 and §E.4 are filled after the fact, because a commit cannot reference its own hash. The phase's own commit writes a `pending-backfill` placeholder and the phase owner backfills the real SHA in a later commit; that is not an ownership crossing.

### Non-transition frontmatter corrections

A frontmatter value written wrong — a stage word in `phase:`, for example — is repaired in place by manager-spec through orchestrator re-delegation. The repair leaves `status:` alone, so it is not an amendment: no `amendment_of:`, no `## Amendments` entry. It touches only the wrong field and `updated:`.

## Optional Fields

These are not required by `FrontmatterSchemaRule`:

| Field | Type | Notes |
|-------|------|-------|
| `issue_number` | integer or null | GitHub issue number; omit when not tracking one |
| `depends_on` | list | SPEC IDs this SPEC needs completed first; checked at run Phase 1 |
| `lint.skip` | list | lint rule codes to skip; for documented debt only |
| `bc_id` | string | breaking-change tracking ID (required when `breaking: true`) |
| `amendment_of` | string | marks an in-place or successor amendment; see Status Enum |
| `tier` | enum `S`, `M`, `L` | complexity tier; absent means Tier L. Artifact sets are in `.claude/rules/moai/workflow/spec-workflow.md`, section "SPEC Complexity Tier (S/M/L)" |

manager-spec also uses `related_specs` (list), `superseded_by` (with `status: superseded`), `partially_superseded_by` (list), `merged_pr`, and `merged_commit`.

## Rejected Snake_Case Aliases

The YAML decoder binds the canonical names only, so an alias is silently dropped and the field reads as empty.

| Do not use | Use |
|------------|-----|
| `created_at:` | `created:` |
| `updated_at:` | `updated:` |
| `labels:` | `tags:` |
| `spec_id:` | `id:` |

## OwnershipTransitionRule

`internal/spec/lint_ownership.go` checks the most recent `status:` transition in the SPEC's git history against the matrix above. The actor is read from the `Authored-By-Agent:` trailer in the commit body, never from the subject.

- `OwnershipTransitionInvalid` (warning): the trailer names an agent that does not own that transition.
- `OwnershipTransitionUnreachable` (info): git history could not be read.
- `OwnershipTransitionUnmeasured` (info): the transition commit has no `Authored-By-Agent:` trailer.

No transition in the lookback window, a transition the matrix does not map, or a trailer naming an agent that owns no transition produces nothing. `--strict` escalates only non-advisory warnings. `lint.skip: [OwnershipTransitionInvalid]` opts a SPEC out of the warning.

## Examples

Correct:

```yaml
---
id: SPEC-AUTH-001
title: "OAuth2 Authentication"
version: "0.1.0"
status: draft
created: YYYY-MM-DD
updated: YYYY-MM-DD
author: Author Name
priority: P1
phase: "v3.0.0"
module: "internal/auth"
lifecycle: spec-anchored
tags: "auth, oauth2, security"
---
```

Wrong — `created_at`, `updated_at`, and `labels` are dropped by the decoder, producing three `FrontmatterInvalid` findings:

```yaml
---
id: SPEC-AUTH-001
title: "OAuth2 Authentication"
version: "0.1.0"
status: draft
created_at: YYYY-MM-DD
updated_at: YYYY-MM-DD
author: Author Name
priority: P1
phase: "v3.0.0"
module: "internal/auth"
lifecycle: spec-anchored
labels: [auth, oauth2]
---
```
