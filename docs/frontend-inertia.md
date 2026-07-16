# Frontend: React + Inertia.js

This branch is a **monolith**: the Go backend renders the React frontend via
[Inertia.js](https://inertiajs.com) using the [gonertia](https://github.com/romsar/gonertia)
adapter. Pages get their data as props from Go — no client-side API calls,
no client router.

## Running in development

Two processes:

```sh
make dev                 # Go backend on http://127.0.0.1:8090
cd frontend && npm run dev   # Vite dev server (HMR)
```

Open **http://127.0.0.1:8090** (the Go server — not the Vite port). While
Vite runs it writes `frontend/.hot` (Laravel-style hot file); gonertia sees
it and points asset tags at the dev server for HMR. Stop Vite and the
backend falls back to the last production build in `frontend/dist`.

First run: `make seed`, then log in with `demo@example.com` / `password123`.

## Production

```sh
cd frontend && npm run build   # hashed assets + manifest into frontend/dist
make build                     # single Go binary; serves assets at /build/*
```

## How a page renders

```
GET /posts
  → internal/web/routes.go      route + middleware (auth cookie, inertia protocol)
  → internal/web/actions.go     postsIndex() builds typed props (internal/types)
  → gonertia                    renders frontend/root.html with the page JSON
  → frontend/src/pages/Posts/Index.tsx   receives props, typed by shared/types.gen.ts
```

- Page components live in `frontend/src/pages/<Component>.tsx`; the
  component name passed to `Render()` maps to that path.
- Props are structs from `internal/types`, so `make types` keeps page props
  typed end-to-end (see docs/type-sharing.md).

## Authentication (PocketBase, cookie-based)

Inertia pages can't attach `Authorization` headers, so this branch wraps
PocketBase auth in an **httpOnly cookie session** (`internal/web/middleware.go`):

- `POST /login` validates credentials against the `users` collection, mints
  a standard PocketBase auth token, and stores it in the `pb_auth` cookie.
- `loadAuthFromCookie` middleware resolves the cookie into `e.Auth` on every
  web request, so route guards and services see the same auth record the
  REST API would.
- `requireWebAuth` redirects guests to `/login` (the web equivalent of
  `apis.RequireAuth()`); `requireGuest` does the reverse.
- `POST /logout` clears the cookie.
- Signup sends PocketBase's verification email; unverified accounts cannot log
  in until they confirm the emailed link.

The PocketBase REST API (`/api/collections/...`) still works normally with
header tokens — useful for realtime subscriptions or mobile clients.

Set `Secure: true` on the cookie (internal/web/middleware.go) when deploying
behind HTTPS.

Set the `users` collection's verification-email template link to:

```text
http://127.0.0.1:8090/verify-email/confirm?token={TOKEN}
```

Use the deployed application origin in production.

## Adding a page

1. Add an action in `internal/web/actions.go` returning
   `i.Render(e.Response, e.Request, "Some/Page", props)`.
2. Declare the route in `internal/web/routes.go`.
3. Create `frontend/src/pages/Some/Page.tsx`.
4. If the props shape is new, add a DTO to `internal/types` and run `make types`.

## Form validation errors

Return `renderError(i, e.Response, e.Request, "Some/Page", gonertia.ValidationErrors{...})`
from an action; they surface in React through `useForm().errors` (see
`frontend/src/pages/Auth/Login.tsx`).
