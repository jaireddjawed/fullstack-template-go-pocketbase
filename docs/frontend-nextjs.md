# Frontend: Next.js + PocketBase auth

This branch runs **two servers**: the Go/PocketBase backend (API only) and a
Next.js app that talks to it — via the typed PocketBase JS SDK for CRUD and
plain `fetch` for the custom `/api/app/*` Go endpoints.

## Running in development

```sh
make dev                        # backend on http://127.0.0.1:8090
cd frontend && npm run dev      # Next.js on http://localhost:3000
```

Open **http://localhost:3000** and sign in through Clerk (see
[auth-clerk.md](auth-clerk.md) for key setup). `make seed` still creates
demo posts, owned by the native demo user.

Configuration: `frontend/.env.example` — `NEXT_PUBLIC_POCKETBASE_URL` plus
the Clerk keys.

## Data access patterns

**Built-in CRUD** — the SDK, typed by the generated collection types:

```ts
const pb = await createServerClient();
const posts = await pb.collection(Collections.Posts).getFullList();  // PostsResponse[]
```

Authorization happens in the backend's collection API rules — the frontend
never filters for security, it just sees what the rules allow.

**Custom Go endpoints** — `fetch` with the generated DTO types:

```ts
const stats = (await res.json()) as PostStats;   // from @shared/types.gen
```

See `src/lib/actions.ts#publishPost` for an authenticated example.

## Authentication

Auth on this branch is handled by **Clerk** — see [auth-clerk.md](auth-clerk.md).
`createServerClient()` attaches the viewer's Clerk session token to every
PocketBase request, and the Go backend resolves it to a users record.

## Adding a page

1. Create `src/app/<route>/page.tsx` (server component by default).
2. Fetch data with `createServerClient()` (CRUD) or `fetch` + DTO types
   (custom endpoints).
3. Mutations go in server actions (`"use server"`), then `revalidatePath`.
4. Protect routes by extending `isProtected` in `src/middleware.ts`.
