---
name: db-migration-conventions
description: Apply this repo's SQL migration conventions — naming, reversibility, idempotency, locking-safe DDL, and how schema constraints should mirror Business strong-type invariants. Use when writing or editing files under business/sdk/migrate/sql.
---

# DB Migration Conventions

This skill defines how to write a migration file in this repo. It is the
authoring counterpart to the SQL-specific checks in `boundary-keeper` (schema
invariants) and `error-tripwire` (transactional failure visibility), which only
review a finished diff. Apply this whenever creating or editing a file under
`business/sdk/migrate/sql`.

## File naming and pairing

- Every migration is a sequential, zero-padded, timestamp- or counter-prefixed
  file pair: an "up" file that applies the change and a matching "down" file
  that reverses it, unless this repo's migration tool uses single-file
  `-- +migrate Up` / `-- +migrate Down` sections — match whatever convention
  already exists in `business/sdk/migrate/sql`, don't introduce a second
  naming scheme.
- One logical schema change per migration file. Don't bundle an unrelated
  index addition into a migration whose name describes a column rename.

## Down migrations must actually reverse

A down migration that is empty, or that drops a table instead of precisely
undoing the up migration's change, is not acceptable. If a column was added
nullable with a default, the down migration drops that exact column — it does
not also touch unrelated columns "while we're in there."

## Idempotency and safety

- Prefer `IF NOT EXISTS` / `IF EXISTS` guards on `CREATE`/`DROP` only when the
  migration tool doesn't already guarantee each migration runs exactly once —
  check the existing migrations for the established pattern before adding
  guards that duplicate what the tool already enforces.
- Never write a migration that silently no-ops on conflict in a way that masks
  a real problem (e.g. `ON CONFLICT DO NOTHING` on a migration meant to
  backfill data) — if rows are expected to already exist, that's fine; if the
  migration's job is to guarantee a specific end state, verify it got there.

## Locking and destructive changes

- Adding a `NOT NULL` column to an existing table needs either a `DEFAULT` (so
  the rewrite is non-blocking on the databases this repo targets) or a
  multi-step migration (add nullable → backfill → add constraint) — don't add
  `NOT NULL` without a default to a table that can have existing rows.
- `DROP COLUMN` / `DROP TABLE` are destructive and irreversible once data is
  gone — confirm in the migration's preceding comment that a backup/backfill
  plan exists, or that the column was never populated in production.
- Adding an index on a large existing table should use the
  non-locking/concurrent form your database supports (e.g. `CREATE INDEX
  CONCURRENTLY` on Postgres) rather than a plain blocking `CREATE INDEX`,
  unless the table is known-small or this is a fresh table created in the same
  migration.

## Schema mirrors Business invariants

Per `layered-architecture-types`, the Business layer enforces invariants with
strong types; the schema should make the *same* invalid states impossible at
the storage level wherever the database can express it, rather than relying on
Business validation alone:

- An enum-like Business strong type (`business/types/<x>/parse.go`) with a
  fixed value set → a `CHECK` constraint or DB-native enum on the column, not
  an unconstrained `text`/`varchar`.
- A required relationship → a `NOT NULL` foreign key with the correct
  `ON DELETE` behavior (`CASCADE`, `RESTRICT`, or `SET NULL` — pick
  deliberately, don't default to whichever Postgres assumes).
- A field that's "NOT NULL DEFAULT '{}'::jsonb" instead of nullable JSONB, per
  the pointer-avoidance guidance in `layered-architecture-types`, when the
  Business layer would otherwise need a wrapper type just to handle NULL.

## Transactions

- A migration that performs multiple statements that must succeed or fail
  together is wrapped in a transaction, unless one of the statements is of a
  kind that can't run inside a transaction on this database (e.g. certain
  concurrent index builds on Postgres) — call that out in a comment rather than
  silently wrapping (or not wrapping) without explanation.
- Seed-data migrations that depend on a previous migration's schema use
  `INSERT ... ON CONFLICT` (or the equivalent upsert) so re-running in a fresh
  environment doesn't fail on existing rows, consistent with the idempotency
  rule above.

## What NOT to do

- Don't edit a migration that has already been merged to `main` — once it can
  have run anywhere, ship a new migration to make the further change.
- Don't reference application code, ORMs, or Business types from a migration
  file — migrations are plain SQL (or this repo's chosen migration DSL), with
  no dependency on `business/*` Go packages.
- Don't add a comment claiming the migration is reversible if the down file
  doesn't actually restore the prior state — see `doc-drift-check` for how this
  gets caught in review; better to get it right at authoring time.