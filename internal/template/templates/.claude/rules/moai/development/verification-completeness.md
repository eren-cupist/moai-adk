---
description: "Verification completeness — a check, gate, acceptance criterion, rule, or assertion is incomplete until its failure has been observed on a known input"
paths: "**/.moai/specs/**,**/.claude/rules/**,internal/template/templates/.claude/rules/**,internal/hook/**,scripts/**,**/.moai/astgrep-rules/**,**/.moai/hooks/**"
---

# Verification Completeness

A verification artifact — a test, check, gate, acceptance criterion or assertion — is complete only when its failure has been observed on a known input. A check whose red has never been seen proves nothing by turning green: the green could mean it works, that it matched nothing, or that it never ran. This rule sits next to claim integrity (`.claude/rules/moai/core/verification-claim-integrity.md`), which forbids reporting a verification that was not run; this one forbids relying on a verification that cannot fail.

## 1. The completion axis

### 1.1 Observed failure

Make the check fail on purpose — an input built to turn it red — and see the red before trusting the green. Watch for checks that print findings but exit with the status of whatever ran last (a script whose trailing `echo` returns 0), and for acceptance criteria that run a check but never use its result.

**An empty sweep asserts nothing.** A test selector that matches zero tests, a grep with zero hits read as a pass, a suite filtered down to no cases — each reports success. Establish how many things were actually checked before reading the verdict. Runners say so in their own words: go `[no tests to run]` / `[no test files]`, pytest `no tests ran`, jest `No tests found`, vitest `No test files found`, cargo a zero pass count.

### 1.2 The three-part check spec

A check specification states together: **(a) when** it must run to be meaningful (a check scheduled before its condition can differ passes vacuously), **(b) the input** that turns it red, observed at least once, and **(c) who sees the red** — which exit code, log level or report carries the failure. A red nobody sees never happened.

### 1.3 Continued firing

A check can stop running without anyone noticing: the event it subscribed to was removed, the installed binary predates the fix, a selector stopped matching. That is absent execution, not a suppressed failure, and nothing goes red. For each check, ask: if this stopped running tomorrow, what would look different? If the answer is nothing, give it a signal that its absence shows up unasked (a count of tests run, a required-check entry, a liveness line).

## 2. Two-cell adoption discipline

An acceptance criterion is adopted as a pair of cells written together: a **RED-now cell** — the criterion observed failing on the pre-implementation tree, pinned to that tree — and a **green path cell** naming the milestone that flips it and what the passing output becomes. One cell alone leaves the criterion unadopted.

RED must be red for the stated reason, and the RED cell says why. A criterion can be vacuous (green already), impossible (red forever, whatever the work does), or red for the wrong reason (failing because of files the work never touches); only stating the reason separates the last two. A green path that depends on someone fixing unrelated files does not measure this work.

**Mutant probe.** Before adopting a criterion, try to construct a change that satisfies the criterion while violating its requirement. If one exists, the criterion is too shallow. A rule checked only on invalid cases passes an all-matching mutant; one checked only on valid cases passes a nothing-matching one.

### 2.1 RED-now cell content — the four elements

For a criterion classified **release-blocking**, the RED-now cell carries all four elements:

- **the command** — a read-only shell invocation completed in a single invocation (no pipes, redirection, `&&`, `;` chains or subshells). A citation not of this form takes the undecidable disposition below.
- **its verbatim stdout** — the raw file bytes the command printed, not a rendered or summarized view.
- **its exit code**, as a separate field (the single-invocation form rules out `; echo $?`). Empty stdout with a non-zero exit is a complete observation.
- **the tree SHA** the measurement was taken on — a commit SHA, never a branch name. A document-level pin binds every criterion without its own pin; a criterion-level pin wins where present.

The elements live in a table cell or in a fenced evidence-ledger entry that a table cell cites by id; the ledger is preferred because tables mangle shell metacharacters. Either carrier is checked the same way.

The obligation is structural: it asks whether a command, its stdout, its exit code and a pinned SHA are present, never how the surrounding prose is worded.

**The undecidable disposition.** When a cited RED cannot be re-executed on the current tree — a historical event, an already-merged state, an external CI result — the criterion loses release-blocking eligibility, is classified as a regression-guard, and is not recorded as a pass.

## 3. Cross-layer revision sweep

When a criterion is rescoped, update the requirement and plan items it cites in the same pass. A revised criterion whose requirement still carries the old scope is false on arrival, and a plan that instructs the side a revised criterion forbids is the same defect one layer up.

## 4. Evidence pinning

Invariant claims — byte-unchanged, preserved surface, absence of something — pin the tree SHA where the evidence was collected, not a moving branch name; upstream movement otherwise makes unchanged work look changed. Re-measure before re-citing a divergence, and re-pin after a rebase. Claims about the mainline itself ("this branch descends from the remote main") keep the moving ref, because that movement is what they assert.

## 5. Corollaries

- To confirm a fix landed, run the pre-fix and post-fix forms and watch them diverge; a grep for the fixed text is satisfied by edits that fix nothing.
- A decision that rests on a rule no gate enforces names the rule file and clause in the artifact that carries the decision (verdict, report, progress record), so the application is reviewable later.
- When measuring a rule's effect, count defects that survived to an approved artifact, not audit findings: a better audit reports more findings at the same defect rate.
