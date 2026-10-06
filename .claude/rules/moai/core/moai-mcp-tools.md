# moai MCP tools

The 47 tools exposed by the self-hosted `moai` MCP server (`.mcp.json` → `{command: "moai", args: ["mcp-server"]}`) carry MoAI's CLI capabilities and are prefixed `mcp__moai__`; the Go producer is `internal/cli/mcp_server.go`. When a tool is in your `tools:` list, prefer it over the equivalent `moai` CLI call: it backs the same implementation, returns structured output and avoids shell quoting. Use the CLI when the tool is not in your list.

## The `project_root` input

Twenty-two tools accept an optional `project_root` string: `spec_progress`, `spec_audit`, `spec_drift`, `verify_snapshot`, `verify_trend`, `codex_audit`, `codex_review`, `codex_task`, `glm_audit`, `glm_review`, `claude_audit`, `audit_multi`, `graph_file_api`, `graph_find_code`, `graph_shortest_path`, `graph_trace_calls`, `factory_decide`, `todo_add`, `todo_list`, `factory_next`, `factory_stage`, and `factory_complete`. It names the tree the call should act on. `codex_task`, `factory_next`, `factory_stage` and `factory_complete` require it and refuse a call without it.

When working inside a worktree, pass your own `git rev-parse --show-toplevel` as `project_root`. The server is a long-lived process whose working directory and environment point at the primary checkout, so without it the call silently acts on the primary checkout — a SPEC that exists only on your branch is simply absent, not reported missing. In the primary checkout, omit it. A path that is not a MoAI project root (or a registered linked worktree of one) is rejected rather than falling back, and accepted paths are canonicalized before use. For `audit_multi`, the root reaches every backend in the fan-out, so all reviewers see the same tree.
