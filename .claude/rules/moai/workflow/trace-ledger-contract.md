---
description: "Executable workflow trace evidence and performance summary contract"
paths:
  - ".claude/skills/moai/workflows/plan*.md"
  - ".claude/skills/moai/workflows/plan/**"
  - ".claude/skills/moai/workflows/run*.md"
  - ".claude/skills/moai/workflows/run/**"
  - ".claude/skills/moai/workflows/sync*.md"
  - ".claude/skills/moai/workflows/sync/**"
  - ".claude/hooks/moai/trace-ledger.sh"
---

# Workflow trace ledger contract

Workflow tracing is opt-in. When `MOAI_TRACE_PHASES=1` is set, record each workflow stage by calling `.claude/hooks/moai/trace-ledger.sh record`, which appends to `.moai/state/workflow-trace.jsonl`. When the flag is off, record nothing and claim no trace. Only rows the helper wrote count as evidence — a comment in a workflow file, a pseudocode line, or an empty ledger is not a trace.

Each JSONL row is append-only and contains:

- `stage` and `event` (`stage_start`, `stage_end`, `tool_start`, or `tool_end`)
- `parent_id` for the session or agent that owns the event
- `input_hash` for the exact routing or verification input
- `cohort` for the comparable S/M/L, cold/warm, or injected-failure cohort
- `cache_hit` (`hit`, `miss`, or `unknown`)
- `retry_cause` (`none` when no retry occurred)
- `duration_ms`, measured by the caller

The helper rejects unsafe field values and non-numeric durations. `trace-ledger.sh summary` counts `tool_end` calls and reports p50/p95 durations plus cache, retry, and cohort counts. Compare cohorts only when their rows share the same input and cohort identity; missing rows are an observation gap, not a zero cost. Routing decisions stay in `.moai/state/routing-ledger.jsonl`; this ledger adds timing and trace identity without changing that schema.
