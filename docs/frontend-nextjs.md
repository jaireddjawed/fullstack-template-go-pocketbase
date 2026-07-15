# Frontend: Next.js + PocketBase auth

This branch runs **two servers**: the Go/PocketBase backend (API only) and a
Next.js app that talks to it — via the typed PocketBase JS SDK for CRUD and
plain `fetch` for the custom `/api/app/*` Go endpoints.

## Running in development

```sh
make dev                        # backend on http://127.0.0.1:8090
cd frontend && npm run dev      # Next.js on http://localhost:3000
```

Open **http://localhost:3000**. First run: `make seed`, then log in with
`demo@example.com` / `password123`.

Configuration: `NEXT_PUBLIC_POCKETBASE_URL` (see `frontend/.env.example`,
defaults to `http://127.0.0.1:8090`).

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

## Authentication (PocketBase SDK + httpOnly cookie)

The SDK's auth store is serialized into an httpOnly `pb_auth` cookie so
**server components can act as the user** (`src/lib/server-auth.ts`):

- `login` server action (`src/lib/actions.ts`) calls
  `pb.collection("users").authWithPassword(...)` on the server and persists
  `{ token, record }` into the cookie.
- `createServerClient()` rebuilds a per-request SDK client from that cookie —
  never share a client across requests.
- `src/middleware.ts` redirects guests away from `/posts` (presence check
  only; the token is actually validated by PocketBase on use).
- `logout` deletes the cookie.

Client components that need PocketBase directly (e.g. realtime
subscriptions) can use `createClient()` — but then auth must be established
client-side too; prefer server components/actions for anything sensitive.

## Adding a page

1. Create `src/app/<route>/page.tsx` (server component by default).
2. Fetch data with `createServerClient()` (CRUD) or `fetch` + DTO types
   (custom endpoints).
3. Mutations go in server actions (`"use server"`), then `revalidatePath`.
4. Protect routes by extending the `matcher` in `src/middleware.ts`.
