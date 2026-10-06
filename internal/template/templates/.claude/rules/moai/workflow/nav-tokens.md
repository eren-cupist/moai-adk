---
description: "NAV:DEC / NAV:SYM binding-token grammar for design docs, read by the navigator graph build"
paths: "**/.moai/project/*.md,**/.moai/docs/**/*.md,**/nav-tokens.md"
---

# Navigator Binding Tokens

The navigator joins design decisions, SPECs and code symbols into `.moai/project/navigator/nav-graph.json`. Two tokens are authored for it:

| Token | Grammar | Where | Links |
|---|---|---|---|
| `@NAV:DEC-<id>` | `<id>` matches `[A-Z][A-Z0-9-]*` | design docs (`.moai/project/{product,structure,tech}.md`, `.moai/docs/**/*.md`) and Go comments | a design decision to a SPEC or symbol |
| `@NAV:SYM:<symbol>` | `<symbol>` matches `[A-Za-z_][A-Za-z0-9_.]*`; `pkg.Name` is conventional, a bare name resolves by suffix | design docs and Go code comments | a doc location to a named code symbol |

Code links to SPECs use the existing `@MX:SPEC:<SPEC-ID>` sub-line (`mx-tag-protocol.md`); don't add a NAV token for that.

```markdown
Decision @NAV:DEC-AUTH-STRATEGY: OAuth2 over client-credentials.
The header parser (see @NAV:SYM:pkg.ParseHeader) extracts the bearer token.
```

The build scans the design docs above and non-test Go files outside `vendor/`; it does not scan `.moai/specs/`, `.moai/reports/` or `.moai/state/`. A token with an empty id or symbol is skipped with a warning in `.moai/logs/navigator-sync.log`, and the build always exits 0. The graph is byte-stable for a given `HEAD`; its schema only grows by addition.
