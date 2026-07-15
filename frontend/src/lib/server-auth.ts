import "server-only";
import { auth } from "@clerk/nextjs/server";
import type { TypedPocketBase } from "@shared/pocketbase.gen";
import type { AuthUser } from "@shared/types.gen";
import { createClient, POCKETBASE_URL } from "./pocketbase";

/**
 * Returns the user's Clerk session token, or null when signed out.
 * The Go backend verifies it and maps it to a PocketBase user record
 * (internal/clerkauth), so it works as a drop-in Authorization header for
 * both the PocketBase CRUD API and the custom /api/app/* routes.
 */
export async function clerkToken(): Promise<string | null> {
  const { getToken } = await auth();
  return getToken();
}

/**
 * Server-side PocketBase client that authenticates every request with the
 * viewer's Clerk token. Create one per request.
 */
export async function createServerClient(): Promise<TypedPocketBase> {
  const pb = createClient();
  const token = await clerkToken();

  if (token) {
    pb.beforeSend = (url, options) => {
      options.headers = { ...options.headers, Authorization: `Bearer ${token}` };
      return { url, options };
    };
  }

  return pb;
}

/**
 * The PocketBase user record for the current Clerk session (provisioned by
 * the backend on first request), or null when signed out.
 */
export async function currentUser(): Promise<AuthUser | null> {
  const token = await clerkToken();
  if (!token) return null;

  const res = await fetch(`${POCKETBASE_URL}/api/app/me`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });
  if (!res.ok) return null;

  return (await res.json()) as AuthUser;
}
