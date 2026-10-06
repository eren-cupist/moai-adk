# Project Navigator — Regeneration Reference

The Project Navigator is a small set of generated files that tell a returning session what the project looks like now: which SPECs exist, where each stands, and what comes next. `/moai sync` regenerates it before the sync commit.

## Files

Under `.moai/project/navigator/` (generated, never hand-edited; the template ships no content there):

| File | Content |
|------|---------|
| `navigator.md` | entry point — current frontier and next task, with a link to the maps |
| `capability-map.md` | one row per SPEC: spec-id, title, status, implementation path, commit SHA, captured-at |
| `progress-map.md` | one row per SPEC: spec-id, status, phase, last-commit SHA and time, frontier milestone |

Every map row carries a 40-character commit SHA and an ISO-8601 time taken from `git log` for the owning file's last commit. Rows are references, never copies of SPEC body content.

Runtime state: `.moai/state/navigator/last-regen-commit.txt` holds the HEAD SHA at the last regeneration (written last, so it means "all three files reflect this commit"). Warnings go to `.moai/logs/navigator-warnings.log`.

```markdown
# Project Navigator

> Living project reorientation brief. Regenerated from the SPEC registry + git log.
> Last regeneration: commit `<sha>` at <iso8601>

## Current frontier
- **SPEC-X-001** (in-progress) — <title> — module `<path>` — [<sha-short> @ <iso8601>]

## Next task
Advance **SPEC-X-001** toward its next milestone. See `progress-map.md`.
```

## Regeneration

```bash
bash .claude/skills/moai-workflow-project/scripts/navigator-regen.sh
```

The script needs only `git`, `awk`, `sed`, and `grep` (no `jq`, no `moai` binary). It resolves the project root from `CLAUDE_PROJECT_DIR` (falling back to `$PWD`), reads `.moai/specs/SPEC-*/spec.md` frontmatter (`id`, `title`, `status`, `phase`, `module`), derives the frontier milestone from `progress.md` where it can, and takes provenance from `git log`. It does not read the `.moai/project/` documents or the @MX corpus.

Behavior a caller can rely on:

- **Always exits 0.** A SPEC with unparseable frontmatter (no closing `---`, no `id`) is skipped with a warning; the other rows are still written.
- **Idempotent.** Timestamps come from the HEAD commit's committer date, not the clock, so two runs on the same commit with no SPEC change produce identical files.
- **Atomic.** Each file is written to `<file>.tmp` and moved into place, so a reader sees the old or the new version, never a partial one. (`NAVIGATOR_PRE_RENAME_BARRIER` exists for concurrency tests.)
- **Empty project.** With no SPECs or no commits, `navigator.md` carries the placeholder `no features tracked yet` and the maps carry only their headers.

## Who reads it

- The SessionStart hook `handle-session-start-navigator.sh` reads `navigator.md` and adds a short brief (under 500 tokens: frontier, next task, link) to each new session. It fails open on a missing or unreadable file. When `last-regen-commit.txt` is more than 3 commits behind HEAD (override with `navigator.staleness_cycles` in `.moai/config/sections/navigator.yaml`), the brief carries a staleness note.
- `/moai plan` can check `capability-map.md` before drafting a SPEC, so an overlapping feature amends the existing SPEC instead of duplicating it.
- `/moai run` can read the owning SPEC's `progress-map.md` row to start oriented.
