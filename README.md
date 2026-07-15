# Fullstack Template — PocketBase × Go — `nextjs`

> **This branch:** Next.js frontend + PocketBase auth (JS SDK with httpOnly
> cookie for SSR). See [docs/frontend-nextjs.md](docs/frontend-nextjs.md).
>
> ```sh
> make dev                     # backend :8090
> cd frontend && npm install && npm run dev    # Next.js :3000 (open this one)
> make seed                    # demo@example.com / password123
> ```

A scaffold for full-stack applications: [PocketBase](https://pocketbase.io) used **as a Go framework** (not as a standalone executable), structured like a conventional backend framework — migrations, seeders, actions, services, routes, hooks, tests, and generated TypeScript types shared with the frontend.

## Branches

| Branch | Frontend | Auth |
|---|---|---|
| `main` | none (backend only, base for the others) | PocketBase |
| `react-inertia` | React + Inertia.js (Vite), served by the Go backend | PocketBase (cookie session) |
| `nextjs` | Next.js (separate dev server) | PocketBase (JS SDK) |
| `nextjs-clerk` | Next.js (separate dev server) | Clerk (verified by the Go backend) |

Start a new project by cloning the branch you want:

```sh
git clone -b nextjs <repo-url> my-app
```

## Quickstart (backend)

```sh
make dev                                     # start backend on http://127.0.0.1:8090
make superuser email=you@example.com pass=changeme123   # dashboard login
make seed                                    # demo user + posts
make test                                    # backend tests
```

- Admin dashboard: http://127.0.0.1:8090/_/
- PocketBase CRUD API: `/api/collections/...` (built in)
- Custom app routes: `/api/app/...` (defined in `internal/routes`)
- Demo login after seeding: `demo@example.com` / `password123`

## Project layout

```
main.go                  entrypoint
migrations/              Go migrations (schema as code, ordered, reversible)
internal/
  app/                   wires everything together (hooks, routes, commands)
  routes/                route declarations — like Laravel's routes/api.php
  actions/               HTTP handlers — like single-action controllers
  services/              business logic, framework-agnostic
  models/                typed record wrappers — like Eloquent models
  hooks/                 record event hooks — like model observers
  seed/                  database seeder
  commands/              custom CLI commands (e.g. `seed`)
  types/                 API DTOs — source of truth for shared TS types
  testutil/              test app + record factories
shared/                  generated TypeScript (types.gen.ts, pocketbase.gen.ts)
docs/                    documentation
```

## Everyday commands

| Command | What it does |
|---|---|
| `make dev` | run the backend (auto-applies pending migrations) |
| `make migration name=x` | create a new blank migration |
| `make migrate-up` / `migrate-down` | apply / revert migrations |
| `make seed` | seed development data (idempotent) |
| `make types` | regenerate shared TypeScript types |
| `make test` | run backend tests |
| `make build` | build a single production binary |

## Documentation

- [docs/architecture.md](docs/architecture.md) — how a request flows, what goes where
- [docs/database.md](docs/database.md) — migrations and seeding
- [docs/type-sharing.md](docs/type-sharing.md) — keeping frontend/backend types in sync
- [docs/testing.md](docs/testing.md) — writing backend tests
