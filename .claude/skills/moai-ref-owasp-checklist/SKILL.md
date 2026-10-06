---
name: moai-ref-owasp-checklist
description: >
  Security review rubric for application code: the checks a /moai review security pass,
  a sync quality gate, or an implementer applies to authentication, authorization, input
  handling, secrets and HTTP surface, with the severity scale MoAI gates on. Not for
  pipeline, container or runtime API operations (moai-ref-secops).

when_to_use: >
  Load when implementing or reviewing code that handles user input, identity, sessions,
  secrets or external requests.

user-invocable: false
metadata:
  version: "2.0.0"
  category: "domain"
  status: "active"
  updated: "2026-03-30"
  tags: "owasp, security, checklist, review, reference"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 1200
---

# Application Security Rubric

## Severity

Report every finding with a severity and a confidence; filtering happens afterwards. The sync quality gate blocks on Critical and High and records Medium and Low as advisory.

| Severity | Meaning | Example |
|----------|---------|---------|
| Critical | Exploitable now, data or account compromise | injection, authentication bypass, secret committed to the repo |
| High | Exploitable with modest effort or insider position | missing object-ownership check, token readable by page scripts |
| Medium | Weakens defenses | no rate limit on login, verbose error details |
| Low | Hardening gap | missing security header on a non-sensitive route |

## Checks

Authorization and identity
- Every endpoint or server action that reads or mutates a resource checks that the caller owns or may access that specific object, not only that the caller is signed in.
- Admin and privileged functions check the role server-side.
- Identity is re-verified against the server-side source of truth (session store, token verification, identity provider) before authorization. A cached or client-decoded session value is not proof of identity.
- Checks in edge middleware, a proxy or an API gateway are a convenience. The handler that performs the mutation checks again.

Input and output
- Input is validated against a schema at the trust boundary (Zod, pydantic, a Go validator) with length and range limits.
- Database access uses parameterized queries; dynamic identifiers such as sort columns come from an allowlist.
- User-supplied content is rendered through the framework's escaping. Raw HTML insertion (`dangerouslySetInnerHTML`, `innerHTML`, template `safe` filters) only after sanitization.
- File paths built from input are normalized and confined to an allowed root. Uploads are checked by size and content type, not only by extension.
- Server-side fetches of user-supplied URLs go through an allowlist and cannot reach internal or cloud-metadata addresses.

Sessions, tokens and secrets
- Session and auth tokens live in `HttpOnly`, `Secure`, `SameSite` cookies, not in `localStorage`. The session ID rotates at login and is invalidated server-side at logout.
- Passwords are hashed with Argon2id or bcrypt; login failures do not reveal which field was wrong; login and reset endpoints are rate limited.
- Secrets come from the environment or a secret manager. None appear in source, committed config, logs or client bundles; client-exposed environment variables (for example `NEXT_PUBLIC_*`) hold nothing secret.
- Webhook receivers verify the provider's signature before trusting the payload. Scheduler-invoked HTTP endpoints require a shared-secret check with a constant-time compare.

HTTP surface
- Production responses do not include stack traces, internal paths or source maps.
- CORS allows specific origins; never a wildcard combined with credentials.
- Sensitive routes send `Strict-Transport-Security`, `Content-Security-Policy`, `X-Content-Type-Options: nosniff` and frame protection (`frame-ancestors` or `X-Frame-Options`).
- Responses expose only the fields the caller may see; serialize through an explicit response shape rather than returning database rows.
