# Backend testing

```sh
make test    # = go test ./...
```

Tests run against a **real PocketBase app** with a throwaway data directory:
`testutil.NewApp(t)` creates a fresh, fully-migrated SQLite database per test
(all migrations in `migrations/` are applied automatically) and binds the
project's hooks and routes — the exact wiring production runs. No mocks.

## Helpers (`internal/testutil`)

| Helper | Purpose |
|---|---|
| `testutil.NewApp(t)` | migrated test app with hooks + routes bound; auto-cleanup |
| `testutil.CreateUser(t, app, email)` | insert a users record |
| `testutil.CreatePost(t, app, ownerID, title, published)` | insert a posts record |
| `testutil.AuthToken(t, record)` | valid auth token for `Authorization` headers |

## Service tests

Call services directly — they take a `core.App`, so the test app drops right
in (`internal/services/posts_test.go`):

```go
func TestStats(t *testing.T) {
    app := testutil.NewApp(t)
    user := testutil.CreateUser(t, app, "a@example.com")
    testutil.CreatePost(t, app, user.Id, "Post", true)

    stats, err := services.NewPostService(app).Stats()
    // assert...
}
```

Note that hooks fire in tests too — creating a post through `app.Save`
triggers the slug hook, which `TestSlugHookFiresOnCreate` verifies.

## API tests

Routes are tested end-to-end (routing, middleware, auth, JSON output) with
PocketBase's `tests.ApiScenario` (`internal/routes/routes_test.go`):

```go
scenario := tests.ApiScenario{
    Method:          "POST",
    URL:             "/api/app/posts/" + post.Id + "/publish",
    Headers:         map[string]string{"Authorization": testutil.AuthToken(t, owner)},
    ExpectedStatus:  200,
    ExpectedContent: []string{`"published":true`},
    DisableTestAppCleanup: true, // testutil.NewApp already registers cleanup
    TestAppFactory: func(tb testing.TB) *tests.TestApp { return app },
}
scenario.Test(t)
```

For every route, cover at least: the happy path, the unauthenticated case
(401), and the wrong-user case (403).

## What to test where

- **Services** — business rules, edge cases. Fast, most of your tests.
- **ApiScenario** — one per route/outcome: status codes, auth boundaries,
  response shape.
- **Collection API rules** (the row-level filters in migrations) can also be
  tested with `ApiScenario` against the built-in `/api/collections/...`
  endpoints if a rule is complex enough to warrant it.
