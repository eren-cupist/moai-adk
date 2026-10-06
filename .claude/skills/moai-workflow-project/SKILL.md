---
name: moai-workflow-project
description: >
  Keeping a project's own documentation current during /moai sync: the
  .moai/project documents (product.md, structure.md, tech.md), the Project
  Navigator files, README, CHANGELOG, docs sites, and API docs. Not for SPEC
  documents (moai-workflow-spec) or code.

when_to_use: >
  Use when updating .moai/project/product.md, structure.md, or tech.md,
  regenerating or auditing the Project Navigator, or updating README,
  CHANGELOG, a docs site, or generated API docs after a change ships.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Write, Edit, Bash(git:*), Bash(npm:*), Bash(npx:*), Bash(uv:*), Bash(pip:*), Bash(ls:*), Bash(mkdir:*), Grep, Glob, WebFetch, WebSearch
user-invocable: false
metadata:
  version: "4.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-04-25"
  modularized: "true"
  tags: "workflow, project, documentation, navigator, readme, changelog, docs-generation"
  aliases: "moai-workflow-project"
  related-skills: "moai-workflow-spec"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000
---

# Project Documentation

This skill covers the documentation `manager-docs` keeps in line with the code during `/moai sync`. The rule throughout: documents describe what the project is now, and a sync updates only what the change affected.

## Project documents (`.moai/project/`)

Three documents give every agent the project context it reads before acting:

| File | Describes | Update when the change |
|------|-----------|------------------------|
| `product.md` | what the product does, for whom, its features | adds, removes, or reshapes a user-facing feature |
| `structure.md` | directory layout, modules, architecture, how layers interact | adds or moves directories or modules, or changes the architecture |
| `tech.md` | languages, frameworks and their versions, tooling, quality gates, deployment | adds, removes, or upgrades a dependency, tool, or framework |

Edit only the affected sections and keep the rest verbatim; these files are partly hand-written, and regenerating them wholesale loses the user's edits. Check claims against the code: framework versions against the lockfile or manifest (`package.json`, `pnpm-lock.yaml`, `go.mod`, `pyproject.toml`), directories against the tree. Leave no placeholder text (`TODO`, `TBD`, `{{...}}`) behind. When the directory does not exist, sync does not create it.

When a document is missing but the directory exists, start it from the matching template in `templates/doc-templates/` (`product-template.md`, `structure-template.md`, `tech-template.md`), filling the `{{PLACEHOLDER}}` fields from the code and dropping sections that do not apply.

## Project Navigator (`.moai/project/navigator/`)

The navigator is a generated, present-tense map of the SPEC registry: `navigator.md` (current frontier and next task), `capability-map.md` (one row per SPEC), and `progress-map.md` (per-SPEC progress). A SessionStart hook reads `navigator.md` to brief each new session. Regenerate it before the sync commit so it is never more than one sync behind:

```bash
bash .claude/skills/moai-workflow-project/scripts/navigator-regen.sh
```

The script reads only SPEC frontmatter, `progress.md`, and `git log`; it always exits 0, skips malformed SPECs with a warning in `.moai/logs/navigator-warnings.log`, and writes atomically. Details: [references/navigator.md](references/navigator.md).

`scripts/navigator-audit.sh` compares `product.md` / `structure.md` / `tech.md` with `capability-map.md` and writes `.moai/project/navigator/audit-report.{md,json}` (features described but without a SPEC, SPECs not reflected in the docs). It is read-only over its inputs; run it when you want evidence that the project documents and the SPEC registry agree. `scripts/navigator-enrich.sh` wraps `moai navigator-enrich`, which adds per-capability source symbols under `.moai/project/codemaps/`; sync does not run it.

## README, CHANGELOG, docs site, API docs

- **CHANGELOG** — Keep a Changelog style: entries under `[Unreleased]` in Added / Changed / Fixed / Removed / Security, one line per user-visible change, naming the SPEC. `manager-docs` runs its B12 self-tests before writing the entry.
- **README** — update the parts a reader acts on: features, install and usage commands, configuration, version references, badges. Verify every command you document against the code or scripts.
- **Docs site** — use the project's existing generator and structure (Nextra, Docusaurus, MkDocs, Sphinx, TypeDoc, and so on; detect it from the config files). Update every locale the site ships, keep navigation and frontmatter consistent with neighboring pages, and build it to confirm.
- **API docs** — regenerate from the source of truth the project already uses (OpenAPI spec, docstrings, TSDoc/JSDoc) rather than hand-writing reference pages.

Write in the user's conversation language unless `.moai/config/sections/language.yaml` sets a different `documentation` language. When a library's current documentation matters (a framework's recommended pattern, a changed API), check the official docs rather than relying on memory.

## Done

The documents the divergence report named are updated; links resolve; Markdown and Mermaid render; the docs build passes where there is one; no credentials appear; the navigator was regenerated; and nothing outside the affected sections changed.
