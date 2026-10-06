---
name: moai
description: >
  MoAI unified orchestrator for autonomous development. Routes natural
  language or subcommands (plan, run, sync, fix, review) to specialized
  agents.
allowed-tools: Agent, AskUserQuestion, Skill, TaskCreate, TaskUpdate, TaskList, TaskGet, Bash, Read, Write, Edit, Glob, Grep
argument-hint: "[subcommand] [args] | \"natural language task\""
---

## Pre-execution Context

!`git status --porcelain 2>/dev/null || true`
!`git branch --show-current 2>/dev/null || true`

## Essential Files

.moai/config/sections/*.yaml

---

## Authority References

Rules and constraints governing all workflows are always loaded from these sources. Do NOT duplicate their content here:

- Core identity, orchestration principles, agent catalog: CLAUDE.md
- Quality gates, security boundaries: .claude/rules/moai/core/moai-constitution.md
- SPEC workflow phases, token budgets: .claude/rules/moai/workflow/spec-workflow.md
- Development methodologies (DDD/TDD): .claude/rules/moai/workflow/spec-workflow.md (Run Phase section)
- Agent definitions: See CLAUDE.md Section 4.
- @MX tag rules and protocol: .claude/rules/moai/workflow/mx-tag-protocol.md

---

## Routing Observation Ledger

When dispatching a subcommand or workflow, the orchestrator records the routing decision to the append-only routing-ledger (`.moai/state/routing-ledger.jsonl`) via `moai harness ledger record` at dispatch time — the request text is piped via stdin and only a privacy-preserving digest is stored, never verbatim user text. As the routed pipeline reaches gate points, machine evidence is appended via `moai harness ledger evidence` (gate exits, audit verdicts, verify-log paths). Outcome is never supplied as an input; it is finalized from machine evidence only. This observation is opt-in and fail-open — it never blocks routing. NOTE: recording depends on the orchestrator actually invoking `moai harness ledger record` at dispatch; when the observability opt-in is ON but that record call is not emitted, the ledger stays empty — an un-recorded dispatch, NOT an opt-in-off no-op. Do not read an empty routing-ledger as 'opt-in disabled'.

---

## Intent Router

### Raw User Input

$ARGUMENTS

### Routing Instructions

[HARD] Route the Raw User Input above using the strict priority order below. Extract the FIRST WORD of the input for subcommand matching. All text after the subcommand keyword is CONTEXT to be passed to the matched workflow — it is NOT a routing signal and MUST NOT influence which workflow is selected.

## Execution Mode Flags (mutually exclusive)

- `--team`: Force agent-team of the Phase 4 4-mode catalog (`.claude/rules/moai/workflow/orchestration-mode-selection.md` §A), subject to its capability gate
- `--solo`: Force serial (sub-agent — single sequential agent per phase)
- No flag: The orchestrator auto-selects from the full 4-mode catalog at Phase 4; the complexity auto-select thresholds are stated once in `orchestration-mode-selection.md` §B.1 (machine source: `workflow.yaml` `auto_selection`) and are not restated here

The `--team` / `--solo` flags are forced overrides onto the catalog; the flag-free default resolves through the catalog decision tree (§B) and its capability gates. The `--mode` dispatch axis is a separate axis — see the crosswalk in `orchestration-mode-selection.md` §G.1 (correspondence, not merge).

### Priority 1: Explicit Subcommand Matching

[HARD] Extract the FIRST WORD from the Raw User Input section above. If it matches any subcommand below (or its alias), route to that workflow IMMEDIATELY. Do NOT analyze the remaining text for routing — it is context for the matched workflow:

[HARD] Mixed-language guard: FIRST-WORD subcommand matching applies only when (a) the input is pure ASCII/Latin, OR (b) the message is prefixed with a literal `/moai ` slash form. When the message contains non-Latin script (Korean/Japanese/Chinese/etc.) beyond the first token, do NOT route immediately on the leading English word — treat it as a possible embedded loanword and fall through to Priority 3 semantic classification of the ENTIRE message. Rationale: CJK technical writing embeds English loanwords such as 'goal', 'run', 'fix', 'plan' at sentence start; immediate first-word routing misfires on them.

- **plan** (aliases: spec): SPEC document creation workflow
- **run** (aliases: impl): DDD/TDD implementation workflow (per quality.yaml constitution.development_mode)
- **sync** (aliases: docs, pr): Documentation synchronization and PR creation
- **fix**: Auto-fix errors in a single pass
- **review** (aliases: code-review): Code review with security and MX tag compliance

### Priority 2: SPEC-ID Detection

Only if Priority 1 did not match: Check if the Raw User Input contains a pattern matching SPEC-XXX (such as SPEC-AUTH-001). If found, route to the **run** workflow automatically. The SPEC-ID becomes the target for DDD/TDD implementation.

### Priority 3: Natural Language Classification

Only if BOTH Priority 1 AND Priority 2 did not match: Classify the intent of the ENTIRE Raw User Input as natural language. This priority is NEVER reached when the first word matches a known subcommand.

[HARD] The cue words listed below are **English exemplars**, NOT literal-match requirements. Classify intent semantically for any `conversation_language` — a Korean, Japanese, Chinese, or other-language request expressing the same intent routes identically. Do not require the literal English tokens to appear.

- Planning and design language (design, architect, plan, spec, requirements, feature request) routes to **plan**
- Security language (security, audit, owasp, vulnerability, injection, xss, csrf) routes to **review** (with `--security` scope)
- Code-review language (review my code, code review, check my PR, look at my changes, take a look at my changes) routes to **review**
- Error and fix language (fix, error, bug, broken, failing, lint) routes to **fix**
- Documentation language (document, sync, docs, readme, changelog, PR) routes to **sync**
- Implementation language (implement, build, create, add, develop) with clear scope routes to **moai** (default autonomous)

### Priority 4: Default Behavior

If the intent remains ambiguous after all priority checks, use AskUserQuestion to present the top 2-3 matching workflows and let the user choose.

If the intent is clearly a development task with no specific routing signal, default to the **moai** workflow (plan -> run -> sync pipeline) for full autonomous execution.

---

## Workflow Quick Reference

### plan - SPEC Document Creation

Purpose: Create comprehensive specification documents using GEARS format with Research-Plan-Annotate cycle.
Phases: Deep Research (research.md) -> SPEC Planning -> Annotation Cycle (1-6 iterations) -> SPEC Creation -> Independent Review (plan-auditor)
Agents: manager-spec (primary), Explore (research), plan-auditor (quality gate), manager-git (conditional)
Skills: moai-workflow-spec (per delegation.yaml)
Flags: --branch, --resume SPEC-XXX, --issue (opt-in; default skips GitHub Issue creation per the late-branch opt-in policy)
For detailed orchestration: Read ${CLAUDE_SKILL_DIR}/workflows/plan.md

### run - DDD/TDD Implementation

Purpose: Implement SPEC requirements through configured development methodology.
Agents: manager-develop (cycle_type=ddd|tdd per quality.yaml, primary), manager-git
Skills: moai-workflow-tdd, moai-workflow-ddd (per delegation.yaml; cycle_type-selected) + domain moai-ref-* injected per mission
Flags: --resume SPEC-XXX, --team (experimental — Agent Teams re-allowed; see Execution Mode Flags)
For detailed orchestration: Read ${CLAUDE_SKILL_DIR}/workflows/run.md

### sync - Documentation Sync and PR

Purpose: Synchronize documentation with code changes and prepare pull requests.
Agents: manager-docs (primary), sync-auditor (quality gate), manager-git
Skills: moai-workflow-project (per delegation.yaml)
Modes: auto, force, status, project. Flags: --auto-merge, --merge (deprecated alias of --auto-merge), --skip-mx
For detailed orchestration: Read ${CLAUDE_SKILL_DIR}/workflows/sync.md

<!-- moai:contract-mode-start id="contract-signing-router" -->
Where `workflow.autonomy.mode: contract` — the Kickoff approval named here is the contract signature checked by `moai contract kickoff-check`; the progression mode is chosen when a goal is armed after that check passes. See `.claude/rules/moai/workflow/contract-autonomy.md` § The signing gate.

<!-- moai:contract-mode-end -->
### fix - Auto-Fix Errors

Purpose: Autonomously detect and fix LSP errors, linting issues, and type errors.
Agents: manager-develop (cycle_type=autofix), Agent(general-purpose) with domain whitelist (fixes)
Skills: moai-workflow-ddd (per delegation.yaml) + domain moai-ref-* injected per mission
Flags: --dry, --sequential, --level N, --resume, --team (experimental — Agent Teams re-allowed; see Execution Mode Flags)
For detailed orchestration: Read ${CLAUDE_SKILL_DIR}/workflows/fix.md

### review - Code Review

Purpose: Multi-perspective code review with security, performance, quality, and UX analysis.
Agents: sync-auditor (review), Agent(general-purpose) with security scope
Skills: moai-foundation-quality, moai-ref-owasp-checklist (per delegation.yaml; per-perspective ref skills injected per lens)
Flags: --staged, --branch, --security, --team (experimental — Agent Teams re-allowed; see Execution Mode Flags)
For detailed orchestration: Read ${CLAUDE_SKILL_DIR}/workflows/review.md

### (default) - MoAI Autonomous Workflow

Purpose: Full autonomous research -> plan -> annotate -> run -> sync pipeline.
Phases: Parallel Exploration (research.md) -> SPEC Generation -> Annotation Cycle -> Implementation -> Sync
Agents: Explore, manager-spec, plan-auditor (quality gate), manager-develop, manager-docs, manager-git, sync-auditor (quality gate)
Skills: moai-workflow-spec, moai-workflow-tdd (per delegation.yaml) + domain moai-ref-* injected per mission
Flags: --loop, --max N, --branch, --pr, --resume SPEC-XXX, --team (experimental — Agent Teams re-allowed; see Execution Mode Flags), --solo, --issue (opt-in; default skips GitHub Issue creation per the late-branch opt-in policy)
For detailed orchestration: Read ${CLAUDE_SKILL_DIR}/workflows/moai.md

---

## Execution Directive

When this skill is activated, execute the following steps in order:

Step 1 - Parse Arguments:
Extract subcommand keywords and flags from the Raw User Input. Recognized global flags: --resume [ID], --seq, --team, --solo. Also detect `ultrathink` keyword in the input text.

**CRITICAL: Deep analysis mode:**
- `ultrathink` keyword detected → Activate Claude's native extended reasoning (xhigh effort mode). This is native Claude behavior with no MCP dependency.

Step 1.5 - Flag-Subcommand Compatibility Validation:
[HARD] After parsing the subcommand and flags (Step 1), validate flag-subcommand compatibility BEFORE routing. If a forbidden combination is detected, STOP all further processing and output an error in the user's conversation_language. Do NOT proceed to Step 2.

Forbidden flag-subcommand combinations:

| Flag | Allowed subcommands | Forbidden subcommands |
|------|---------------------|------------------------|
| `--branch` | `plan`, default (autonomous) | `run`, `sync` |

Rationale: `--branch` creates the feature branch at SPEC initialization, so `/moai run` and `/moai sync` MUST operate on the branch `plan` already established — re-creating it mid-lifecycle corrupts the SPEC lifecycle and is rejected at the router level.

The retired `--worktree` flag is handled separately: a request carrying it is not a forbidden-combination error but a retired flag. Tell the user that plan no longer creates a workspace, and that entering one first is the replacement.

Error message template (Korean conversation_language; substitute the actual flag and subcommand):
```
에러: --branch 플래그는 /moai plan 전용입니다.
/moai run 과 /moai sync 는 plan 단계에서 만든 브랜치를 그대로 씁니다.

올바른 사용법:
  /moai plan SPEC-XXX --branch    (브랜치 생성)
  /moai run SPEC-XXX              (기존 브랜치 재사용)
  /moai sync SPEC-XXX             (기존 브랜치 재사용)

--branch 플래그를 뺀 형태로 다시 실행하세요.
```

Retired-flag message (`--worktree`):
```
안내: --worktree 플래그는 폐기됐습니다. plan 은 더 이상 작업 공간을 만들지 않습니다.

격리된 공간에서 작업하려면 먼저 들어간 뒤 plan 을 실행하세요:
  moai cc -w <이름>              (그 자리에서 진입)
  moai cg -w <이름> --spawn      (새 tmux 창, 현재 세션 유지)
  /moai plan "<설명>"
```

For English (`en` conversation_language), translate the message; the structure remains identical.

Step 2 - Route to Workflow:
Apply the Intent Router (Priority 1 through Priority 4) to determine the target workflow. If ambiguous, use AskUserQuestion to clarify with the user.

Step 2.2 - Record Routing Decision:
Immediately after routing resolves (Step 2), record the routing decision to the append-only routing-ledger (`.moai/state/routing-ledger.jsonl`) so that auto-invocation is observable. Run:

```
echo "<raw request text>" | moai harness ledger record --subcommand <matched> --mode <phase-4-mode> --tier <tier> --level <harness-level> --session <session-id>
```

The request text is piped via stdin and only a privacy-preserving digest is stored, never verbatim user text (policy source: § Routing Observation Ledger above). This step is opt-in and fail-open: if the `moai` CLI is absent from PATH or the command exits non-zero, log nothing and continue — it NEVER blocks routing, never gates the workflow, and never triggers a retry loop. An un-recorded dispatch is an observation gap, not an error.

[HARD] Beginner-Friendly Option Design:
All AskUserQuestion calls throughout MoAI workflows MUST follow these rules:
- The first option MUST always be the recommended choice, clearly marked with "(Recommended)" suffix — this is the `push`-mode branch; while `interview.recommendation_mode` is `pull` the suffix is withheld from every option and no option carries a preference claim (`.claude/rules/moai/core/askuser-protocol.md` § Recommendation Placement Principles)
- Every option MUST include a detailed description explaining what it does and its implications

Step 2.8 - Requirement Analysis & Completion Condition:
Before loading the workflow body (Step 3), produce a requirement-analysis record for the routed request:

1. **Requirement summary** (1-3 sentences): what the user asked for, restated in the orchestrator's own words.
2. **Completion condition**: the end state that means "done". Where the condition is machine-verifiable (test exit code, lint-clean state, grep count, bounded turn count), express it as one measurable end state + a stated check + a bound clause; the orchestrator evaluates the condition text per-turn.
3. **Pipeline contract**: `full-pipeline` (default natural-language route — run-phase completion auto-chains into sync) or `single-phase` (explicit `run`/`sync` subcommand — chaining is offered as the "(Recommended)" next-step option, never fired silently).
4. **Orchestration-shape pre-signal**: an early input to the Phase 4 4-mode selection (`orchestration-mode-selection.md` §A) — noted here, decided at Phase 4.

Trivial-scope exemption: skip this step entirely for `sync` status mode, and any Stage-1-Clarify exception per `askuser-protocol.md` § Ambiguity Triggers and Exceptions.
Socratic-first ordering: while intent clarity is below 100%, run the Socratic interview (per `askuser-protocol.md`) BEFORE deriving the completion condition — the condition encodes drained intent, never a guess.
A derived completion condition NEVER authorizes autonomous run-phase entry — Implementation Kickoff Approval remains mandatory at the plan→run boundary.

Step 3 - Load Workflow Details:
Read `workflows/<name>.md` for the target subcommand. (Agent Teams is experimental and re-allowed: a `--team` flag selects the Agent Teams layer, subject to the constraints in `.claude/rules/moai/workflow/orchestration-mode-selection.md` §C.1. Only the static layer stays retired, so there is no separate `team/<name>.md` workflow file — the same `workflows/<name>.md` is read either way. Historical: the retired era emitted `MODE_TEAM_UNAVAILABLE` and fell back to sub-agent mode; the sentinel is retained as documented history.)

Step 4 - Read Configuration:
Load relevant configuration from the .moai/config/sections/*.yaml section files as needed.

Step 5 - Initialize Task Tracking:
Use TaskCreate to register discovered work items with pending status.

Step 6 - Execute Workflow Phases:
Follow the workflow-specific phase instructions. Delegate all implementation to appropriate agents via Agent(). Collect user approvals at designated checkpoints via AskUserQuestion. Before each implementation/review Agent() spawn, inject 0-3 `At start, invoke Skill("<name>") for <reason>` lines per the delegation map (`.moai/config/sections/delegation.yaml`).

Step 7 - Track Progress:
Update task status using TaskUpdate as work progresses (pending to in_progress to completed).

Step 8 - Present Results:
Display results to the user in their conversation_language using Markdown format.

Step 9 - Declare Completion:
When all workflow phases complete successfully, state that the workflow is complete in the Completion Report (banner / prose) so the result is unambiguous.

Step 10 - Guide Next Steps:
Use AskUserQuestion to present the user with logical next actions based on the completed workflow.

---

Version: 2.8.0
Last Updated: 2026-07-07
