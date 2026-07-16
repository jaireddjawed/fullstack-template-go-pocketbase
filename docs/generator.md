# Project generator (`fullstack-template init`)

An interactive TUI (Bubble Tea) that scaffolds a new project from this
template by picking a stack:

```
fullstack-template init
  → project name / Go module path
  → frontend:  none · Next.js · React + Inertia
  → auth:      PocketBase · Clerk · WorkOS (coming soon)
  → database:  PocketBase (SQLite)
  → extras:    email verification · shadcn/ui
  → confirm → generate
```

## Running

From a checkout of this repo (clones your local branches — works offline):

```sh
make generator          # builds ./fullstack-template
./fullstack-template init
```

Or directly: `go run ./cmd/fullstack-template init`. Outside a checkout,
install it with:

```sh
go install github.com/jaireddjawed/fullstack-template-go-pocketbase/cmd/fullstack-template@latest
```

Non-interactive (scripts, CI):

```sh
go run ./cmd/fullstack-template init --no-input \
  --name my-app --module github.com/you/my-app \
  --frontend next --auth clerk --extras shadcn
```

Flags: `--name`, `--module`, `--frontend`, `--auth`, `--database`,
`--extras` (csv), `--dir` (default `./<name>`), `--repo` (template source,
defaults to the local checkout when run inside it, else GitHub), `--no-input`.

## How generation works

Each supported stack lives on a template branch; the generator resolves the
selection matrix to a branch and post-processes a clone of it:

| frontend | auth | branch |
|---|---|---|
| none | PocketBase | `main` |
| Next.js | PocketBase | `nextjs` |
| Next.js | Clerk | `nextjs-clerk` |
| React + Inertia | PocketBase | `react-inertia` |

Steps (see `internal/scaffold/generate.go`):

1. `git clone --depth 1 --branch <branch>` into the target directory.
2. Strip `.git` and rewrite the Go module path everywhere
   (`.go`, `go.mod`, docs, `tygo.yaml`, TS configs).
3. Create Docker files:
   - `Dockerfile` (+ frontend image for Next.js, asset-embedding image for
     Inertia), `.dockerignore`, `docker-compose.yml`.
4. Apply extras:
   - **Email verification** — a migration setting the users collection's
     `AuthRule` to `verified = true` (PocketBase auth only; Clerk manages
     its own verification).
   - **shadcn/ui** — no files; the post-generate checklist prints the
     `npx shadcn init` steps (requires a frontend).
5. Fresh `git init` + initial commit, then print next steps.

## Architecture

- `internal/scaffold` — all logic: option matrix + availability rules
  (`config.go`), generation pipeline (`generate.go`), extras (`extras.go`).
  Fully unit-tested, including an end-to-end test that generates a project
  from the local repo and asserts on the result.
- `cmd/fullstack-template/tui` — the Bubble Tea wizard. It only *collects*
  a `scaffold.Config`; disabled options (unsupported combos, coming-soon
  items) are greyed out and unselectable. Model logic is tested by driving
  `Update` with key events.
- `cmd/fullstack-template` — CLI entry: flags, repo detection, and the
  non-interactive path.

## Adding a new option

1. Add the enum value and its `Choice` (with availability rules) in
   `internal/scaffold/config.go`; extend `branchMatrix` if it maps to a new
   template branch.
2. Implement it: a new branch for stack options, or an `applyExtra` case
   (generated files / next-step instructions) for extras.
3. Add matrix + availability tests in `scaffold_test.go`.

The TUI picks up new choices automatically — it renders whatever the
`*Choices()` functions return.
