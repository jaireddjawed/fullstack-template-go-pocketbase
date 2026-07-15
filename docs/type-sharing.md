# Type sharing between Go and TypeScript

Frontend and backend types stay in sync through **code generation from a
single source of truth** — never by hand-maintaining parallel definitions.
One command regenerates everything:

```sh
make types
```

Both outputs land in `shared/` and are committed, so the frontend always
compiles against the current backend contract. If a Go type changes and the
frontend wasn't updated, `tsc` fails — that's the sync guarantee.

## Two generators, two sources of truth

### 1. Custom API DTOs → tygo

`internal/types` holds the request/response structs for every custom
`/api/app/*` route. [tygo](https://github.com/gzuidhof/tygo) converts them to
TypeScript interfaces in `shared/types.gen.ts` (config: `tygo.yaml`).

```go
// Go (internal/types/api.go)          // Generated TS (shared/types.gen.ts)
type PostStats struct {                 export interface PostStats {
    Total     int64 `json:"total"`        total: number;
    Published int64 `json:"published"`    published: number;
    Drafts    int64 `json:"drafts"`       drafts: number;
}                                       }
```

**Workflow:** change the Go struct → `make types` → the frontend picks up the
new shape. Never edit `*.gen.ts` files.

### 2. PocketBase collections → pocketbase-typegen

[pocketbase-typegen](https://github.com/patmood/pocketbase-typegen) reads the
migrated SQLite database and emits typed records for every collection
(`UsersRecord`, `PostsRecord`, typed `Collections` enum, ...) into
`shared/pocketbase.gen.ts`. These give you full typing on the PocketBase JS
SDK:

```ts
import PocketBase from "pocketbase";
import { type TypedPocketBase, Collections } from "@shared/pocketbase.gen";

const pb = new PocketBase(url) as TypedPocketBase;
const posts = await pb.collection(Collections.Posts).getFullList(); // PostsRecord[]
```

`make types` runs `migrate up` first, so the generated types always reflect
the latest migrations.

## Frontend wiring

Frontend branches map `@shared/*` to `../shared/*` in `tsconfig.json`, so
imports look the same everywhere:

```ts
import type { PostStats } from "@shared/types.gen";
```

## When schema or DTOs change

1. Write the migration / change the Go struct.
2. `make types`
3. Fix any TypeScript errors the change surfaces.
4. Commit the regenerated `shared/*.gen.ts` together with the backend change.
