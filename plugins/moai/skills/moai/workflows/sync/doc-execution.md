---
description: "Sync Phases 11-12 — SPEC-implementation divergence analysis, the documentation scope approval (gate-sync-2), and the documentation and SPEC status update."
user-invocable: false
metadata:
  parent: moai-workflow-sync
  phase: "Phase 11-12: Divergence Analysis and Documentation Update"
---

# Sync Phases 11-12: Documentation

The goal is documentation that describes what was actually built. Phase 11 works out what changed relative to the SPEC and what needs updating; the user approves that scope; Phase 12 has `manager-docs` write it.

## Phase 11: Divergence analysis

Requirements: `.moai/` and `.claude/` exist and the project is a git repository.

Collect the change set: the commits since the base plus any uncommitted changes, grouped by area (backend, frontend, tests, config, docs). Note whether this session is in a MoAI worktree (the git directory path contains `/.moai/worktrees/`, or `.moai/worktrees/registry.json` has an active entry for the SPEC); Phase 13 uses that to pick the delivery route and next-step options. Worktree context never decides auto-merge on its own.

Compare the SPEC with reality. Read `spec.md`, `plan.md`, and `acceptance.md`, then list:

- `scope_expansion` — features or endpoints built beyond the SPEC
- `unplanned_additions` — files, directories, or dependencies the plan did not name
- `deferred_items` — planned items or acceptance criteria not implemented
- `structural_changes` — unplanned refactors or architecture changes

The SPEC's `lifecycle` frontmatter decides what to do about divergence (`.claude/rules/moai/development/spec-frontmatter-schema.md`):

- `spec-anchored` (default) — the SPEC should match the implementation. SPEC body updates are `manager-spec`'s to make: after approval, the orchestrator delegates them to `manager-spec` before `manager-docs` runs.
- `spec-lite` — the SPEC is closed with a short implementation-notes section (also written by `manager-spec`), not a full rewrite.
- `exploratory` — the SPEC is not modified; divergence is reported as warnings for the user.

From this, plan the documentation updates: which of CHANGELOG, README, docs site, API docs, and `.moai/project/` documents need changes, and what the PR description should say (deployment notes from Phases 2-6 included).

### Documentation scope (`gate-sync-2`)

Show the user the divergence report, the planned documentation updates, the SPEC changes it implies, and the PR description draft when a PR will be made. Ask with AskUserQuestion: proceed (recommended), change the plan (then redo this phase), show the full details, or abort. Nothing has been written before this point, so aborting here leaves the tree as it was.

<!-- moai:contract-mode-start id="contract-doc-scope" -->
Where `workflow.autonomy.mode: contract` — the documentation scope approval (`gate-sync-2`) is not asked: regenerate the documents the divergence report names within the contract's ownership, and escalate when the scope would reach outside it. See `.claude/rules/moai/workflow/contract-autonomy.md`, section "Gate disposition".

<!-- moai:contract-mode-end -->
## Phase 12: Documentation and SPEC status

Brief `manager-docs` with the approved plan, the divergence report, the change set, the Phase 1-10 results, and the audit verdict. It is the only writer in this phase. It:

- adds the CHANGELOG `[Unreleased]` entry after its duplicate, acceptance-criteria count, and file-path self-tests;
- updates README, the docs site (every locale the project ships), and API docs where user-facing behavior changed;
- updates `.moai/project/product.md`, `structure.md`, and `tech.md` only in the sections the change affects — new directories, technologies, or features — when that directory exists;
- regenerates the navigator with `bash .claude/skills/moai-workflow-project/scripts/navigator-regen.sh` (it always exits 0 and writes `.moai/project/navigator/`);
- updates the SPEC status: `in-progress → implemented → completed` in `spec.md` frontmatter when every acceptance criterion is met; when the implementation is partial, the status stays `in-progress` and the report names what remains;
- writes the Sync-phase Audit-Ready Signal section of `progress.md`, with `sync_commit_sha` as the `pending-backfill` placeholder until the commit exists;
- checks its output: links resolve, Markdown and Mermaid render, the docs build passes where the project has one, and no credentials appear in any document.

Documents are written in the user's conversation language unless `language.yaml` sets a different `documentation` language. If `manager-docs` reports that a SPEC body change is needed, delegate that to `manager-spec`, then resume `manager-docs`.
