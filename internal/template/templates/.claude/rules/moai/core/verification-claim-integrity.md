---
description: "Evidence rules for verification and completion reports, SPEC acceptance criteria and audit verdicts: the five-section report, moving-ref attribution, tool provenance, and refused commands"
paths: "**/.moai/specs/**,**/.moai/reports/**"
---

# Verification-Claim Integrity

The base rule is in `AGENTS.md`, section "Evidence": claim a verification, a completion, a defect or the premise of a recommendation only after observing it with the tool that decides it, in this run, against this tree. A missing failure signal is not a pass, a grep match is a hypothesis rather than a verified defect, and a reference existing does not show that the referenced thing is still live — check that before recommending to keep something. This file adds the detail that SPECs, audits and reports need.

## 1. Report format

Verification and completion reports, including the auditors' verdicts, use five sections:

| Section | Carries |
|---|---|
| **Claim** | what is asserted, one claim per row or sentence |
| **Evidence** | the command that was run and its verbatim deciding output; a summary is not evidence |
| **Baseline-attribution** | the tree and commit it was measured against, in this run |
| **Gaps** | what was not observed; an empty Gaps section asserts that nothing was left unobserved |
| **Residual-risk** | what could still be wrong despite what was observed |

When a command a verification depended on was refused rather than run (worktree guard, permission deny, policy gate), say so in Gaps together with what you did instead. Reading the source as a fallback is fine and the verdict may still pass, but the report must not present the inference as a measurement. Refused tool calls are also recorded as `tool_failure:<tool>:<category>` rows in `.moai/lessons-inbox.jsonl`.

## 2. Attribution

An attributed claim names the command and the output observed in this run. A figure remembered from another task, package or point in time is a carry-over: report it under Gaps, not as a claim.

A measurement made by the project's own tooling has two coordinates: the tree it read and the build that judged it. A tool found on `PATH` may be an installed build older than the tree, and a stale build passes silently on checks it does not have. Build the tool from the tree and invoke it by path, or state that the installed build's commit is not behind the tree's `HEAD`. This does not apply where there is no repository to compare against.

When an acceptance criterion depends on a baseline measured before the change, commit the baseline before the change: git cannot show authoring order within one commit.

### 2.1 Moving refs

A claim measured against a moving ref (`origin/main`, `origin/HEAD`, a line number) stays textually the same while its truth drifts. First decide what the ref is to the claim. If the claim was measured at that address, it is an ANCHOR: substituting today's SHA keeps its meaning for a later reader, and re-running it next week should give the same answer. If the claim is about the moving thing itself (what mainline carries, which tip to start from), it is a SUBJECT. Then pick the remedy:

| Remedy | Class | Form |
|---|---|---|
| R1 pin | ANCHOR, value known now | replace the ref with the resolved SHA |
| R2 freeze at pre-flight | ANCHOR, value not known yet (the usual run-phase case) | `BASELINE_SHA=$(git rev-parse origin/main)` before the first run-phase commit; decide the criterion against it and record the value |
| R3 declare | SUBJECT, narrative | keep the ref and add `<!-- moving-ref-ok: <reason> -->` on that line or the line above; an empty reason does not suppress |
| R4 measure at read time | SUBJECT, a reader will act on it | lead with the command that decides the claim; any value follows as a dated reference |

If a claim reads false only because upstream commits moved the ref, remediate as above; if this work's own commits made it false, fix the work instead. `moai spec lint` flags moving refs and prints these remedies.
