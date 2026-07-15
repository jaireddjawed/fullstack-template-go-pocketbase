# Authentication: Clerk

This branch replaces PocketBase's native auth with [Clerk](https://clerk.com).
Clerk owns sign-up/sign-in/sessions in the Next.js app; the Go backend
verifies Clerk session tokens and maps each Clerk identity to a record in
the PocketBase `users` collection — so **collection API rules, relations
(`posts.owner`), and `apis.RequireAuth()` keep working unchanged**.

## Setup

1. Create a Clerk application at https://dashboard.clerk.com.
2. Frontend: copy `frontend/.env.example` to `frontend/.env.local` and fill
   in `NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY` and `CLERK_SECRET_KEY`.
3. Backend: export the same secret key before starting:

   ```sh
   CLERK_SECRET_KEY=sk_test_... make dev
   cd frontend && npm run dev
   ```

Without `CLERK_SECRET_KEY` the backend still runs (with a warning), but all
requests are anonymous.

## How it works

```
Browser ──(Clerk session JWT in Authorization header)──▶ Go backend
             internal/clerkauth middleware:
               1. verify JWT signature against Clerk's JWKS
               2. look up users record by clerk_id
                  └─ first sight: provision one (email/name from the
                     Clerk management API, random password, verified;
                     an existing record with the same email is linked
                     instead of duplicated)
               3. set e.Auth → rules/guards behave as native auth
```

Key pieces:

- `internal/clerkauth/clerkauth.go` — verifier + middleware + provisioning.
  The verifier is an interface, so tests inject a fake and exercise the full
  HTTP flow without Clerk credentials (`clerkauth_test.go`).
- `migrations/1752500001_add_clerk_id_to_users.go` — adds `clerk_id`
  (unique when non-empty) to `users`.
- `GET /api/app/me` — returns the PocketBase user record for the current
  token; the frontend uses it to learn its own PocketBase user id.

## Frontend integration

- `src/middleware.ts` — `clerkMiddleware` protects `/posts`.
- `src/app/layout.tsx` — `ClerkProvider`, `SignInButton`, `UserButton`.
- `src/lib/server-auth.ts` — `createServerClient()` returns a PocketBase SDK
  client that sends the Clerk token as a Bearer header on every request, so
  even the built-in CRUD API is authorized as the Clerk user:

  ```ts
  const pb = await createServerClient();
  const posts = await pb.collection(Collections.Posts).getFullList();
  // rules evaluate with @request.auth = the provisioned users record
  ```

## Notes

- PocketBase's own token auth still works in parallel (the middleware skips
  requests that already resolved a native token) — handy for the seeded
  demo user, scripts, and the admin dashboard.
- Users created by Clerk get a random password; password login for them is
  effectively disabled, Clerk is the source of truth.
- Deleting a user in Clerk does not delete the PocketBase record; add a
  Clerk webhook handler if you need that.
