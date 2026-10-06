---
paths: "**/.moai/specs/**,**/.claude/agent-memory/**"
---

# MoAI Memory

Two memory stores are involved, with different rules:

| Store | Scope | Format |
|---|---|---|
| `~/.claude/projects/<project>/memory/` | Claude Code auto-memory: the `MEMORY.md` index plus topic files. All worktrees of one repository share it | Plain markdown; Claude Code loads `MEMORY.md` every session and truncates it at 200 lines or 25KB |
| `.claude/agent-memory/<agent-name>/` | Per-agent memory, gitignored and per project root | Frontmatter required (below), audited by a PostToolUse hook |

Write all memory files in English, whatever the conversation language. Store stable facts and lessons, not session state: code conventions already visible in the code, git history, fix recipes captured in a commit, and anything already in `CLAUDE.md` or `.claude/rules/` do not belong in memory. SPEC documents under `.moai/specs/` are the main cross-session record for SPEC work.

## The MEMORY.md index

Topic files are loaded on demand and found only through the index, so every topic file keeps exactly one index line, written as `- [title](file) — short summary`. When the index grows, shorten lines (take the summary from the topic file's `description`) rather than removing or merging them; a removed line makes its file unreachable. A completed one-off record can move to `memory/_archive/` with its line removed; never delete it outright. Keep each topic file's `description` to one line of about 150 characters, since it is what recall matches against.

`moai memory doctor` lists every candidate store with whether it exists and its file and index-line counts. Run it once bare to see the resolved paths, then pass the one you mean with `--dir`, and only trust a zero finding count when the store reports `exists: true`.

## Agent memory files

Every file in `.claude/agent-memory/<agent-name>/` starts with:

```markdown
---
name: <short name>
description: <one-line summary used for relevance matching>
type: <user | feedback | project | reference>
---
```

`feedback` and `project` files lead with the rule or fact and add `**Why:**` and `**How to apply:**` lines. A filename prefix (`project_`, `feedback_`, `reference_`) should match the type. The hook's warnings are `MEMORY_MISSING_TYPE`, `MEMORY_MISSING_FRONTMATTER`, `MEMORY_BODY_STRUCTURE_MISSING` and `MEMORY_EXCLUDED_CATEGORY`; they never block, and `MOAI_MEMORY_AUDIT=0` turns them off during bulk migrations. At session start, files older than 24 hours are flagged as possibly stale; verify them before acting on them.

A worktree is its own project root, so agent memory written there would die with the tree. Writes through Write/Edit are mirrored to the primary checkout's store automatically; for anything written another way, `moai memory drain` previews and, with `--yes`, copies worktree agent memory into the primary store without overwriting existing files.
