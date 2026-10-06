---
paths: "**/.moai/specs/**,**/sprint-round-naming.md"
---

# Epic Naming Taxonomy

Four terms name units of work. Use them with these meanings so a reader can tell a group of SPECs from a step inside one SPEC.

| Term | Meaning | Korean |
|------|---------|--------|
| **Epic** | A group of one or more SPECs bundled by schedule, release, or theme. Numbered (`Epic 1`) or named (`Epic Docs-v3`); may be split into parallel lanes (`Epic 1 Lane A`). | 에픽 |
| **SPEC** | One unit of work (feature, refactor, bugfix) with a `SPEC-{DOMAIN}-{NUM}` ID, carried through plan → run → sync. | (English) |
| **Milestone** | An ordered step inside one SPEC (`M1`, `M2`, …), and the unit manager-develop is delegated. | 마일스톤 |
| **Constitution** | The project's governing technical choices — stack, naming, forbidden libraries, architecture, security and logging standards — recorded in `.moai/project/tech.md` and `.moai/config/sections/constitution.yaml`. | 컨스티튜션 |

Identifiers, file names, frontmatter, and rule references use the English terms. User-facing text in a Korean conversation may use the Korean term.

When a plan.md belongs to a multi-SPEC group, its context section names the Epic. A standalone SPEC needs no Epic.

## Legacy aliases

Older SPECs, commits, and memory entries use retired terms. Read them with this mapping, and do not use them in new content:

| Legacy term | Current term |
|-------------|--------------|
| `Sprint`, `Sprint N Lane A` | `Epic`, `Epic N Lane A` |
| `cohort` | an Epic, or a lane within one |
| `Round` (a step inside one SPEC) | `Milestone` |
| `Wave` | `Epic` (CI pipeline "wave" numbering is unrelated) |
