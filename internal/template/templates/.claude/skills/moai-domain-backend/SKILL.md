---
name: moai-domain-backend
description: >
  Backend implementation and review checklist for server code (Go, Python, Node/TypeScript
  services and BFFs): service boundaries, outbound calls, transactions, background work,
  configuration and observability, with pointers to the API, security and database
  references. Not for UI work or for API contract detail (moai-ref-api-patterns).

when_to_use: >
  Load when implementing or reviewing server-side logic, service integrations or
  authentication flows.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Write, Edit, Bash(npm:*), Bash(npx:*), Bash(node:*), Bash(uv:*), Bash(pip:*), Bash(pytest:*), Bash(ruff:*), Bash(docker:*), Bash(curl:*), Bash(go:*), Bash(cargo:*), Grep, Glob
user-invocable: false
metadata:
  version: "2.0.0"
  category: "domain"
  status: "active"
  updated: "2026-01-11"
  modularized: "false"
  tags: "backend, server, services, review"
  author: "MoAI-ADK Team"
---

# Backend Checklist

Start from the service's existing structure: its layering, error type, logger, config loader and test helpers. New code uses those rather than introducing parallel ones. Related references: moai-ref-api-patterns for endpoint contracts, moai-ref-owasp-checklist for security, moai-domain-database for schema and query changes, and the language rule under `.claude/rules/moai/languages/` for commands.

- Boundaries: handlers parse and validate input, call domain logic, and map results and errors to responses. Domain logic does not import HTTP or framework types.
- Outbound calls: every network call has a timeout or deadline, propagates the request context or cancellation, and handles non-2xx and malformed responses from the upstream. Retries are bounded, use backoff, and apply only to idempotent operations.
- Transactions: multi-step writes that must succeed together run in one transaction; no network calls are made while a transaction holds locks.
- Background work: jobs and consumers are idempotent or deduplicated, record failures somewhere visible, and stop cleanly on shutdown.
- Configuration: settings and secrets come from the environment or the project's config loader, are validated at startup, and fail fast when missing.
- Errors: errors are wrapped with context where they cross a layer, logged once at the boundary that handles them, and never swallowed.
- Observability: logs are structured, carry a request or trace ID, and contain no secrets, tokens or personal data.
- Concurrency: shared state is guarded, goroutines, tasks and promises have an owner that waits for or cancels them, and concurrent paths are tested with the race detector where available.
