---
name: moai-domain-database
description: >
  Review checklist for schema changes, migrations, queries and ORM models (SQL databases,
  document stores, caches): migration safety on live data, constraints, indexes for the
  queries actually issued, and transaction and connection handling. Not for API contracts
  (moai-ref-api-patterns) or application security (moai-ref-owasp-checklist).

when_to_use: >
  Load when adding or changing a migration, a model or ORM mapping, a query, or a caching
  layer.

license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read, Write, Edit, Bash(psql:*), Bash(mysql:*), Bash(sqlite3:*), Bash(mongosh:*), Bash(redis-cli:*), Bash(npm:*), Bash(npx:*), Bash(prisma:*), Bash(neonctl:*), Bash(firebase:*), Bash(supabase:*), Grep, Glob
user-invocable: false
metadata:
  version: "3.0.0"
  category: "domain"
  status: "active"
  updated: "2026-04-25"
  tags: "database, schema, migrations, queries, review"
  author: "MoAI-ADK Team"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 900
---

# Database Change Checklist

Use the project's migration tool and its file naming, and generate migrations with it rather than hand-numbering them. Never edit a migration that has already been applied to a shared environment; add a new one. Running a migration or any write against a shared or production database is an outward-facing action that needs the user's approval.

Migrations
- Each migration is reversible, or the irreversibility is stated in the SPEC and the commit.
- Changes to live tables are backward compatible with the code currently deployed: add a nullable column or a column with a default, backfill, switch the code, then add constraints or drop the old column in a later migration.
- Long operations on large tables (adding an index, rewriting a column type, adding a constraint with validation) use the database's online form where it has one, for example `CREATE INDEX CONCURRENTLY` in PostgreSQL.
- Data backfills run in batches and are restartable.

Schema and models
- Integrity lives in the database: primary keys, foreign keys, `NOT NULL`, unique and check constraints match the domain rules the code assumes.
- ORM model definitions and the migration agree on column types, defaults (including generated primary keys) and nullability; a mismatch is a finding.
- Money uses exact decimal types; timestamps are stored with time zone or as UTC.

Queries
- Every new query pattern on a growing table is backed by an index whose leading columns match its filter and sort; check the plan (`EXPLAIN`) for anything on a hot path.
- No per-row queries inside loops (N+1); load related data with a join, a batched `IN` query or the ORM's eager loading.
- Queries are parameterized; dynamic sort or filter columns come from an allowlist.
- Unbounded reads are paginated or streamed.

Transactions, connections and caches
- Multi-statement writes that must succeed together share one transaction with the isolation level the invariant needs.
- Connection pools are sized for the deployment (serverless and edge runtimes need a pooler or an HTTP driver).
- Cached values have an explicit TTL and an invalidation path on write; the cache is never the only copy of data.
