# Migrations & Seeding

## Migrations

Schema lives in code as ordered, reversible Go migrations in `migrations/`.
They are registered via `init()` (the blank import in `cmd/app/main.go`) and applied
automatically on `serve`, or explicitly:

```sh
make migrate-up            # apply pending
make migrate-down          # revert the last one
make migration name=add_comments   # create a blank migration file
```

Each migration has an `up` and `down` function and runs in a transaction:

```go
func init() {
    m.Register(func(app core.App) error {
        // up: create/alter collections
        return nil
    }, func(app core.App) error {
        // down: revert it
        return nil
    })
}
```

See `migrations/1752500000_create_posts.go` for a full example: fields,
indexes, and API rules (row-level authorization) are all declared there.

### Automigrate

`migratecmd` is registered with `Automigrate: true`: when you edit
collections in the admin dashboard (http://127.0.0.1:8090/_/) during
development, PocketBase **writes the corresponding Go migration file into
`migrations/` for you**. Review and commit it like any other code. This is
the fastest way to design schema — click it together in the dashboard, get
the migration for free.

### Rules of thumb

- Never edit an already-committed migration; add a new one.
- `pb_data/` (the SQLite database itself) is gitignored — the schema is
  reproducible from migrations on any machine.
- Data backfills are just migrations that use the ORM
  (`app.FindAllRecords`, `app.Save`, raw `app.DB()` queries).

## Seeding

`internal/seed/seed.go` populates development data, exposed as a CLI command:

```sh
make seed    # = go run ./cmd/app seed
```

Seeders are **idempotent** — every record is looked up before being created,
so re-running is safe. Records are saved through `app.Save`, so hooks fire
exactly as they would in production (the demo posts get their slugs from the
slug hook, not from the seeder).

Seeded demo login: `demo@example.com` / `password123`.

To add seed data, extend `seed.Run` — group each concern into its own
`ensureX` function, mirroring Laravel's one-seeder-per-model convention.
