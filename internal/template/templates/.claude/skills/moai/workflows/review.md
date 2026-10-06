---
description: >
  Code review of the change the user means (branch vs base, a PR, staged
  changes, a commit, or given paths): finds correctness, security, test, and
  consistency defects, checks @MX tag compliance, and returns a ranked findings
  list with severity, confidence, and file:line. Optional lean (over-engineering)
  and deep security modes.
user-invocable: false
metadata:
  version: "3.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-21"
  tags: "review, code-review, security, quality, mx, audit"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["review", "code review", "security audit", "quality check", "code analysis"]
  agents: []
  phases: ["review"]
---

# Workflow: Review

`/moai review` finds the defects in a change and hands the user a findings list they can act on. The review runs in two separate steps: first find everything, then filter and rank what to present. Keeping them apart matters — a reviewer told to report only important issues drops real bugs it was unsure about, while a reviewer that reports everything and ranks afterwards loses nothing.

## Scope

Review the diff the user means:

| Input | Diff |
|-------|------|
| no flag, on a feature branch | `git diff <base>...HEAD` plus uncommitted changes, where `<base>` is `git_strategy.{mode}.main_branch` from `.moai/config/sections/git-strategy.yaml` (default `main`) |
| no flag, on the base branch | uncommitted changes; if there are none, the last commit (`git diff HEAD~1`) |
| a PR number or URL | `gh pr diff <n>` (and `gh pr view <n>` for the description) |
| `--staged` | `git diff --staged` |
| `--branch <base>` | `git diff <base>...HEAD` |
| `--commit <SHA>` | `git show <SHA>` |
| `--file <path>` or bare paths | those files as they are now |
| `--repo` | the whole tree (only with `--lean` or `--deep`) |

Read surrounding code wherever the diff alone does not show whether something is wrong — callers of a changed function, the type a field belongs to, the test that covers a branch. A review that judges only the changed lines misses the bugs that live at their edges.

If the scope is empty, say so and stop.

## Step 1: Find

Go through the change and record every issue you find, including ones you are unsure about or consider minor. For each finding record:

- `path:line`
- category: `correctness`, `security`, `tests`, `performance`, `consistency`, `api` (breaking or surprising interface change), or `mx`
- severity: `critical` (data loss, security breach, crash in a main path), `high` (wrong behavior users will hit), `medium` (wrong in edge cases, missing error handling, weak tests), `low` (minor or cosmetic)
- confidence: `high` (verified against the code), `medium` (likely, not fully traced), `low` (suspicion worth a look)
- what is wrong, what it causes, and a concrete fix

Do not filter for importance or confidence at this step.

What to look for, in rough order of value: logic errors and unhandled cases; error paths that swallow or mis-report failures; concurrency and async hazards (races, missing cancellation, unawaited promises, React effects with stale closures or missing cleanup); resource leaks; security issues (below); tests that are missing for changed behavior or that would pass even if the code were wrong; performance problems that matter at realistic sizes (N+1 queries, unbounded loops or allocations, needless re-renders on hot paths); divergence from the project's existing conventions; public API, schema, or config changes that break callers.

### Security

Check the change for OWASP-relevant issues: injection (SQL, command, template), broken authentication or authorization (missing checks, IDOR), unvalidated input at trust boundaries, XSS and CSRF in web code, SSRF, unsafe deserialization, sensitive data in logs or responses, and secrets in code. When the change touches authentication, sessions, permissions, data access, request handling, or user input — or `--security` is set — load `Skill("moai-ref-owasp-checklist")` (and `Skill("moai-ref-secops")` for CI/CD, container, or infrastructure changes).

Secrets: scan the diff for private keys (`-----BEGIN [A-Z]+ PRIVATE KEY-----`), AWS access key IDs (`AKIA[0-9A-Z]{16}`), GitHub tokens (`ghp_[A-Za-z0-9]{36}`), and other provider tokens or passwords. With `--security` or `--deep`, also scan reachable history:

```bash
git log -p --all -G '(-----BEGIN [A-Z]+ PRIVATE KEY-----|AKIA[0-9A-Z]{16}|ghp_[A-Za-z0-9]{36})'
```

Ignore only the example access key ID that AWS publishes in its documentation (it ends in `EXAMPLE`). Report a secret found in history separately from one in the working tree — a historical leak needs rotation even after the file is fixed.

With `--security`, also run the ecosystem's dependency audit when the tool is installed (`pnpm audit` / `npm audit`, `govulncheck ./...`, `pip-audit`), and report a missing tool as skipped rather than clean.

### @MX tag compliance

For changed code, check the @MX tag protocol (`.claude/skills/moai/references/mx-tag.md`): a function with three or more callers needs `@MX:ANCHOR`; a dangerous pattern (goroutine without cancellation, async work without error handling, shared mutable state) needs `@MX:WARN` with `@MX:REASON`; a non-obvious exported function benefits from `@MX:NOTE`; an untested public function should carry `@MX:TODO`. Report missing or outdated tags as `mx` findings.

## Step 2: Filter and rank

Now decide what to present. Re-check each high-severity claim against the code once; drop a finding only when the code clearly shows it is wrong, and merge duplicates. Then group:

1. **Fix first** — `critical`/`high` correctness and security findings with `high` or `medium` confidence.
2. **Worth fixing** — other correctness and security findings, and `medium` or higher findings in tests, performance, api, and consistency.
3. **Lower priority or uncertain** — `low` severity, and `low`-confidence findings of any severity (say why you are unsure).
4. **@MX tags** — `mx` findings.

Nothing found in Step 1 is silently discarded: groups 3 and 4 may be terse, but they are listed.

## Report

```markdown
## Review: <scope, e.g. feature/x vs main, 12 files>

<One or two sentences: overall assessment and the most important problem.>
Findings: <n> total — <n> fix first, <n> worth fixing, <n> lower priority, <n> MX.

### Fix first
1. [high · confidence high · correctness] src/cart/total.ts:42 — Discount is applied after tax, so totals are overstated for every discounted order. Fix: apply `discount()` before `addTax()`.

### Worth fixing
...

### Lower priority or uncertain
...

### @MX tags
...
```

Write the report in the user's conversation language; keep code, paths, and identifiers as they are.

Then ask with AskUserQuestion what to do next: fix the "Fix first" findings now (recommended when there are any), fix a chosen subset, save the report to `.moai/reports/review-<timestamp>.md`, or stop. When fixing, make the changes in this session, run the affected tests, and report what changed. `/moai fix` is the separate route for lint, type, and LSP errors.

## Large diffs

Review in this session. Subagents multiply cost and time: each re-establishes context, re-explores and reports back. Only when a diff is large and splits into genuinely independent areas (for example a frontend app and a Go service in the same branch) is it worth running one read-only subagent per area in a single message so they run in parallel. Give each the area's file list, the base, the Step 1 finding format, and the instruction to report everything it finds. Do Step 2 yourself on the merged findings.

## --lean: over-engineering audit

With `--lean`, run only this audit — no correctness, security, or performance review. The narrowness is the point: a focused "what can be cut" pass over the scope (the diff by default, the whole tree with `--repo`). It is read-only and advisory: it changes no files and gives no pass/fail verdict.

Report every over-engineering finding under exactly one of these tags, one per line:

| Tag | Flags |
|-----|-------|
| `delete:` | unused or speculative code — dead branches, never-called helpers, write-only config |
| `stdlib:` | hand-rolled logic the language's standard library already provides |
| `native:` | a dependency or code duplicating a feature the platform already provides |
| `yagni:` | premature generality — single-implementation abstraction, single-caller indirection, dead config knob |
| `shrink:` | logic that fits in fewer lines without losing clarity |

```
L<line>: <tag> <what to cut>. <replacement>. [path]
```

Before reporting a `yagni:` finding, read the deferred-debt markers with `moai mx query --kind DEBT` (read-only). A site that already carries `@MX:DEBT` is reported as `L<line>: yagni: <site> [already tracked @MX:DEBT — deferred]. [path]`. Never create or edit `@MX:DEBT` markers here.

Close with exactly one line: `net: -<N> lines possible` (diff scope), `net: -<N> lines, -<M> deps possible` (repo scope with removable dependencies), or `Lean already. Ship.` when nothing warrants removal. The over-engineering anti-patterns are catalogued in `.claude/skills/moai/references/anti-patterns.md`.

## --deep: deep security scan

With `--deep` (security-focused; composes with `--security`), go beyond the single pass: map entry points and trust boundaries for the scope (`--repo`, `--staged`, `--branch <base>`, or `--commit <SHA>`; the diff by default), enumerate the attack surface, then hunt for vulnerabilities across it. For each candidate, try to refute it before reporting: is the vulnerable path reachable from attacker-controlled input, is the impact real, and do existing defenses (validation, authorization, framework guards) already stop it? A candidate that survives is confirmed; one you could neither confirm nor refute goes under "Unconfirmed candidates", never mixed with confirmed findings. Confidence is `high` only when reachability, impact, and the lack of a defense are all verified.

For a large scope, hunting can be split by area into parallel read-only subagents as described under Large diffs. When a confirmed critical or high finding rests on a long chain of reasoning, one independent read-only subagent given only the claim and the relevant code (not your reasoning) is a useful cross-check.

Write the results to `.moai/reports/security-deepscan-<timestamp>/`, which is a report, not a SPEC:

```
.moai/reports/security-deepscan-<timestamp>/
  .gitignore        # one line: *   (keeps the results out of commits)
  report.md         # F1, F2, ... each with Severity, Confidence, Impact, Exploit scenario, Recommendation; then "## Unconfirmed candidates"
  findings.jsonl    # one finding per line
  revision.json     # what was scanned
  patches/          # only with --patch
```

```json
{"id":"F1","severity":"high","confidence":"medium","title":"...","impact":"...","exploit":"...","recommendation":"...","location":{"path":"...","hint":"..."}}
{"scanned_commit":"<sha>","working_tree_included":true,"scope":"repo|branch|commit|staged|file","generated_at":"<timestamp>"}
```

A scan with no confirmed findings still writes the directory and says "0 confirmed findings".

`--patch` (only with `--deep`, off by default): for each confirmed finding, draft the minimal fix in a scratch worktree (`git worktree add` to a temporary path, or a subagent with `isolation: "worktree"`), never in the user's working tree. Keep a patch only if it addresses that one finding, adds no new vulnerability, and leaves other behavior unchanged; otherwise write a short note explaining why. Save each as `patches/F<i>.patch` or `patches/F<i>.note.md`. Never apply, stage, commit, or push a drafted patch — the user applies them.
