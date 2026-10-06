---
name: moai-ref-api-patterns
description: >
  Review checklist for HTTP API changes (REST, GraphQL, RPC handlers, Next.js route
  handlers and server actions): contract stability, error shape, pagination, validation
  and status codes. Used when implementing or reviewing backend endpoints. Not for security
  review (moai-ref-owasp-checklist), schema design (moai-domain-database) or UI work.

when_to_use: >
  Load when adding or changing an endpoint, its request or response shape, or its error
  handling.

user-invocable: false
metadata:
  version: "2.0.0"
  category: "domain"
  status: "active"
  updated: "2026-03-30"
  tags: "api, rest, graphql, contract, review, reference"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 800
---

# API Change Checklist

Match the conventions the service already uses (URL style, envelope, error codes, pagination) before introducing new ones; consistency across endpoints matters more than any one convention. Where the API is generated from or checked against a schema (OpenAPI, GraphQL SDL, protobuf), change the schema first and regenerate clients in the same change.

- Contract: removing or renaming a field, changing its type or nullability, changing a default, tightening validation, or changing an error code is a breaking change for existing callers. It needs a new version or a migration path agreed with the consumers, and the SPEC should say which.
- Errors: one error shape across the service, with a stable machine-readable code, a human message, per-field details for validation errors, and a request or trace ID. No stack traces or internal identifiers in production responses.
- Status codes: 400 for malformed input, 422 for well-formed input that fails validation (or 400 if the service already uses it for both), 401 versus 403 kept distinct, 404 for objects the caller may not see, 409 for state conflicts, 429 with `Retry-After`. Success with an error inside a 200 body is a finding.
- Validation: request bodies, query parameters and path parameters are parsed through a schema at the handler boundary, and the parsed value, not the raw input, is what the rest of the code uses.
- Collections: every list endpoint is paginated with a server-enforced maximum page size; use cursors for large or frequently changing collections.
- Idempotency: retried writes (payments, creation from unreliable clients) accept an idempotency key or are naturally idempotent.
- Tests: each new or changed endpoint has a test for its success response shape and its main error responses.
