---
name: moai-ref-testing-pyramid
description: >
  Rubric for judging whether a change's tests sit at the right level (unit, integration,
  end-to-end) and test the right things; used when writing tests and when sync-auditor
  assesses test adequacy. Not for coverage thresholds or how tests are run
  (moai-workflow-testing).

when_to_use: >
  Load when deciding what kind of test a change needs, or when reviewing a change's test mix.

user-invocable: false
metadata:
  version: "2.0.0"
  category: "domain"
  status: "active"
  updated: "2026-03-30"
  tags: "testing, pyramid, test-levels, review, reference"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 800
---

# Test Level Rubric

Coverage targets and test commands live in moai-workflow-testing and quality.yaml; this page is only about which level a test belongs at and what it should assert.

Choosing the level:

- Pure logic, transformations, validation, and error mapping: unit tests, fast and with no network, database or filesystem beyond a temp dir.
- Behavior that only exists across a boundary (SQL against a real schema, an HTTP handler with its middleware, a queue consumer): integration tests against the real dependency, for example a test container or the framework's test server. Mock only what you do not own and cannot run.
- A user-visible flow that must keep working end to end (sign-in, checkout, the main path through a feature): a small number of end-to-end tests, not one per edge case.

When reviewing, flag each of these as a finding with severity and confidence:

- Most of the change's tests are end-to-end while the logic underneath has no unit tests.
- A unit test reaches a real external service, or depends on test order or wall-clock time.
- An integration boundary the change introduced or altered (schema, API contract, serialization) has no integration test.
- Assertions check internal calls or private state rather than observable behavior.
- Tests cover framework or third-party behavior, or generated code, instead of the project's own logic.
- A test is flaky, skipped without a reason, or passes without asserting anything meaningful.
