---
name: moai-ref-secops
description: >
  Defensive review rubric for the operational layer: CI/CD pipelines, infrastructure-as-code,
  container images and Kubernetes workloads, and the runtime controls on a deployed API.
  Used by /moai review security hunts and the sync security gate. Not for application-code
  security (moai-ref-owasp-checklist), API design (moai-ref-api-patterns), or offensive
  techniques.

when_to_use: >
  Load when a change touches CI workflows, Dockerfiles, Kubernetes or Terraform manifests,
  gateway or WAF configuration, or API rate and query limits.

user-invocable: false
metadata:
  version: "2.0.0"
  category: "domain"
  status: "active"
  updated: "2026-06-24"
  tags: "devsecops, container, kubernetes, cicd, iac, api-operations, reference"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 1000
---

# Operational Security Rubric

Report every finding with a severity and a confidence, using the scale in moai-ref-owasp-checklist (the sync gate blocks on Critical and High). Describe the misconfiguration, how to detect it and how to fix it; do not write exploitation steps.

Pipeline (CI/CD)
- Workflow tokens and cloud credentials have the least permissions the job needs; prefer short-lived OIDC credentials over stored long-lived keys.
- Third-party actions and images are pinned to a commit SHA or digest, not a moving tag.
- Workflows triggered by forks or external pull requests do not receive secrets, and do not run untrusted code on a shared self-hosted runner.
- Secrets are never echoed into logs, written to artifacts or baked into build outputs; a secret scanner runs on the repository and on pull requests.
- Dependency audit and static analysis run in the pipeline, and their failures are not silently ignored.

Infrastructure as code
- Manifests are scanned for misconfiguration before apply: public storage buckets, wildcard IAM actions or principals, open security groups, unencrypted data stores, disabled audit logging.
- State files and plan outputs that contain secrets are stored encrypted with restricted access.

Containers and Kubernetes
- Images use a minimal, pinned base, are scanned before deploy, and contain no secrets in any layer (check build args and copied files).
- Containers run as non-root with a read-only root filesystem, dropped capabilities, `allowPrivilegeEscalation: false`, and a seccomp profile; no privileged mode, host networking, or host-path mounts without a stated need.
- ServiceAccounts follow least privilege: no `cluster-admin` binding or wildcard Role for a workload, and token automount is disabled where the pod does not call the API.
- Pod Security admission (or an equivalent policy engine) enforces the baseline; resource requests and limits are set.

Running API
- Every endpoint enforces object- and function-level authorization on the server, so access to another tenant's object fails even with a valid session.
- Rate limits, request-size caps and timeouts are enforced at the gateway or in the service, with stricter limits on authentication and expensive endpoints.
- GraphQL endpoints cap query depth and complexity and disable introspection in production unless intended.
- WAF or gateway rules have been tuned against real traffic before switching to blocking mode, and logs make denied requests and authorization failures observable.
- Old API versions and undocumented endpoints are inventoried and retired rather than left reachable.
