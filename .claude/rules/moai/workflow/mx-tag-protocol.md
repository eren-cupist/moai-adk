---
paths: "**/*.go,**/*.py,**/*.ts,**/*.js,**/*.java,**/*.rs,**/*.c,**/*.cpp,**/*.rb,**/*.php,**/*.kt,**/*.swift,**/*.dart,**/*.ex,**/*.scala,**/*.hs,**/*.zig"
---

# @MX Tag Protocol

`@MX` tags are code comments that carry context between sessions: why code is shaped the way it is, what must not break, and where the danger is. Read the tags in a file before changing it; add or update them for code you write. Full reference: `.claude/skills/moai/references/mx-tag.md`.

## Syntax

```
// @MX:TYPE: [AUTO] description
// @MX:SUB_KEY: value
```

Use the file's comment prefix (`//` for Go, TS/JS, Java, Rust, C/C++, Swift, Kotlin, Dart, Zig, Scala; `#` for Python, Ruby, Elixir; `--` for Haskell). Tags an agent adds carry the `[AUTO]` prefix. Write descriptions in the `code_comments` language from `.moai/config/sections/language.yaml` (English when unset).

| Type | Meaning | Add when |
|---|---|---|
| `@MX:ANCHOR` | invariant contract other code depends on | fan-in reaches `thresholds.fan_in_anchor` (3), a public API boundary, an external integration point |
| `@MX:WARN` | danger zone | goroutines/async work without cancellation, cyclomatic complexity ≥ 15, if-branches ≥ 8, global state mutation |
| `@MX:NOTE` | context or intent | a magic constant, an unexplained business rule, a long exported function without docs |
| `@MX:TODO` | incomplete work | a public function without tests, an unimplemented SPEC requirement, an unhandled returned error |
| `@MX:DEBT` | a deliberate, working simplification | the code is correct within a known limit and has a known revisit trigger |

Sub-lines: `@MX:REASON` (required on every ANCHOR and WARN), `@MX:SPEC` (only when a SPEC exists), `@MX:LEGACY`, `@MX:TEST`, `@MX:PRIORITY`, and for DEBT `@MX:CEILING` (the limit it is valid up to) and `@MX:UPGRADE` (the trigger to revisit it).

```
// @MX:DEBT: in-memory map cache, no eviction
// @MX:CEILING: < 10k entries
// @MX:UPGRADE: switch to LRU when entry count exceeds 10k
```

## Lifecycle

- **TODO** is added in RED/ANALYZE and removed when the test passes or the work is done; one left unresolved for more than three iterations becomes a WARN.
- **ANCHOR** is updated when callers or the SPEC change, and demoted to NOTE (with a note in the report) when fan-in drops below the threshold — never deleted automatically.
- **WARN** is removed when the dangerous structure is gone; structural ones (a goroutine's lifecycle) persist.
- **NOTE** is re-checked when the signature changes and removed with its code.
- **DEBT** persists, by design, until its `@MX:UPGRADE` trigger fires and the simplification is replaced; it does not escalate to WARN. A DEBT without `@MX:UPGRADE` has no exit condition — `moai mx query --kind debt` marks it `"rotRisk": "no-trigger"`.

## Limits and exclusions

Per-file limits come from `.moai/config/sections/mx.yaml` `limits` (defaults: 3 ANCHOR, 5 WARN, 10 NOTE, 5 TODO). Over the limit, demote the ANCHOR with the lowest fan-in and keep the highest-priority WARNs. Files matching the per-language `exclude` patterns there (generated code, vendor, mocks) are not tagged.

## Report

After a phase that changed tags, list tags added, removed and updated, and any new WARN that needs a reviewer's attention.
