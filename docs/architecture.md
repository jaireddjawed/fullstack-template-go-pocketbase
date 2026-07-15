# Architecture

PocketBase is embedded as a Go library. You get its admin dashboard, auth,
CRUD REST API, realtime subscriptions, and file storage for free, and extend
it with your own schema, routes, and logic — organized like a conventional
MVC-ish framework.

## How a request flows

### Built-in CRUD (`/api/collections/...`)

PocketBase handles these itself. Authorization is declared per collection as
**API rules** (see `migrations/1752500000_create_posts.go`), e.g.:

```
ListRule: published = true || owner = @request.auth.id
```

Rules are row-level filters evaluated against the authenticated request —
most standard CRUD needs no Go code at all.

### Custom routes (`/api/app/...`)

For anything beyond CRUD, requests flow:

```
routes  →  actions  →  services  →  PocketBase ORM (core.App)
```

| Layer | Package | Responsibility | Laravel analogue |
|---|---|---|---|
| Routes | `internal/routes` | map URL + method + middleware to an action | `routes/api.php` |
| Actions | `internal/actions` | parse request, call service, shape HTTP response | single-action controllers |
| Services | `internal/services` | business logic; no HTTP knowledge | service classes |
| Models | `internal/models` | typed wrappers over records (PocketBase "record proxies") | Eloquent models |
| Hooks | `internal/hooks` | react to record events from *any* source (API, dashboard, Go code) | model observers |

Example: `POST /api/app/posts/{id}/publish`

1. `internal/routes/routes.go` declares the route and binds `apis.RequireAuth()`.
2. `internal/actions/posts.go#PublishPost` reads the path param and `e.Auth`,
   calls the service, and maps `services.ErrNotOwner` to a 403.
3. `internal/services/posts.go#Publish` loads a `models.Post`, checks
   `post.IsOwnedBy(...)`, and saves it.
4. Saving fires hooks (`internal/hooks`) — e.g. the slug generator.

## Models

`internal/models` wraps records in typed structs via PocketBase's record
proxy pattern (`core.BaseRecordProxy`). A model *is* its record — pass it to
`app.Save()` directly and every hook and validation fires as usual:

```go
post, _ := models.FindPostByID(app, id)
post.SetPublished(true)
app.Save(post)
```

Models are optional sugar: plain `core.Record` access is always available
(`record.GetString("title")`), and hooks receive plain records that you can
wrap with `models.NewPost(e.Record)` when convenient. Add accessors as you
need them rather than exhaustively mirroring every field.

## Wiring

`internal/app/app.go` is the composition root:

- `New()` builds the production app (migrate command, custom CLI commands,
  hooks, routes).
- `Bind(app core.App)` attaches hooks + routes to *any* app instance —
  production and tests share this, so tests exercise the real wiring.

## Adding a feature (checklist)

1. **Schema** — `make migration name=create_comments`, edit the file in
   `migrations/`, model it on the posts migration. Set API rules there.
2. **DTOs** — add request/response structs to `internal/types`, run `make types`.
3. **Model** (optional) — `internal/models/comment.go` if Go code touches the
   collection's fields in more than one place.
4. **Service** — create `internal/services/comments.go` with the logic.
5. **Action + route** — handler in `internal/actions`, declare it in
   `internal/routes/routes.go`.
6. **Hooks** — if the logic must also run for dashboard/CRUD writes, put it
   in `internal/hooks` instead of the action.
7. **Tests** — service tests + an `ApiScenario` per route (see docs/testing.md).
8. **Seed** — extend `internal/seed` if the dev environment needs sample data.
