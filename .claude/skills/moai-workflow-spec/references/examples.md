# SPEC Example (Tier M)

A complete Tier M SPEC in the shapes `moai spec lint` reads. The IDs are placeholders.

## spec.md

```markdown
---
id: SPEC-AUTH-001
title: "Refresh token rotation"
version: "0.1.0"
status: draft
created: YYYY-MM-DD
updated: YYYY-MM-DD
author: Author Name
priority: P1
phase: "v1.4.0"
module: "src/auth"
lifecycle: spec-anchored
tags: "auth, session, security"
tier: M
---

# Refresh token rotation

## HISTORY

| Version | Date | Author | Change |
|---------|------|--------|--------|
| 0.1.0 | YYYY-MM-DD | Author Name | Initial draft |

## Context

Refresh tokens never expire today, so a leaked token grants access until the user changes their password. Rotating the token on every refresh and rejecting reused tokens limits a leak to one use.

## Scope

The refresh endpoint and the session store. The login and logout endpoints keep their current behavior.

## Requirements

- REQ-AUTH-001: When a client presents a valid refresh token, the auth service shall issue a new access token and a new refresh token and invalidate the presented one.
- REQ-AUTH-002: When a client presents a refresh token that was already rotated, the auth service shall reject the request with 401 and revoke every refresh token in that session.
- REQ-AUTH-003: The auth service shall expire a refresh token 14 days after it was issued.
- REQ-AUTH-004: While the session store is unavailable, the auth service shall reject refresh requests with 503 and shall not issue tokens.

## Exclusions (What NOT to Build)

### Out of Scope — Login flow
- Changes to password login or MFA.

### Out of Scope — Device management UI
- A screen listing active sessions.
```

## acceptance.md

```markdown
# Acceptance Criteria — SPEC-AUTH-001

- AC-AUTH-001: Given a valid refresh token, When the client calls POST /auth/refresh, Then the response is 200 with a new access token and a new refresh token, and the old refresh token returns 401 on reuse. (maps REQ-AUTH-001)
- AC-AUTH-002: Given a refresh token that was already rotated, When the client calls POST /auth/refresh with it, Then the response is 401 and every refresh token of that session returns 401. (maps REQ-AUTH-002)
- AC-AUTH-003: Given a refresh token issued 14 days and 1 minute ago, When the client calls POST /auth/refresh, Then the response is 401. (maps REQ-AUTH-003)
- AC-AUTH-004: Given the session store is stopped, When the client calls POST /auth/refresh, Then the response is 503 and no token is written. (maps REQ-AUTH-004)

## Definition of Done

- The four criteria above pass in the auth integration test suite.
- Typecheck and lint pass with no new warnings.
```

## plan.md (excerpt)

```markdown
## Approach

Store a rotation counter per session; a refresh compares the presented token's counter with the stored one. Follows the existing session repository in `src/auth/session-repo.ts`.

## Milestones

### M1 — Rotation and reuse detection
Exit: AC-AUTH-001, AC-AUTH-002 green

### M2 — Expiry and store failure
Exit: AC-AUTH-003, AC-AUTH-004 green

## Risks

- Clients that retry a refresh after a timeout will trip reuse detection. Mitigation: accept the previous token for 10 seconds after rotation.

## Assumptions

- The 14-day lifetime matches the current access policy; no product decision was needed.

[NEEDS CLARIFICATION: revoke scope]
On reuse, revoke only the session or every session of the user? This changes whether other devices are logged out.
```
