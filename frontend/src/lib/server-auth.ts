import "server-only";
import { cookies } from "next/headers";
import type { TypedPocketBase, UsersResponse } from "@shared/pocketbase.gen";
import { AUTH_COOKIE, createClient } from "./pocketbase";

/**
 * Server-side PocketBase client hydrated from the httpOnly auth cookie.
 * Use in server components, server actions, and route handlers.
 */
export async function createServerClient(): Promise<TypedPocketBase> {
  const pb = createClient();
  const cookie = (await cookies()).get(AUTH_COOKIE);

  if (cookie?.value) {
    try {
      const { token, record } = JSON.parse(cookie.value);
      pb.authStore.save(token, record);
    } catch {
      pb.authStore.clear();
    }
  }

  return pb;
}

/** The authenticated user, or null. */
export async function currentUser(): Promise<UsersResponse | null> {
  const pb = await createServerClient();
  return pb.authStore.isValid ? (pb.authStore.record as UsersResponse) : null;
}

/** Serializes the SDK auth store into the cookie after login. */
export async function persistAuth(pb: TypedPocketBase): Promise<void> {
  (await cookies()).set(
    AUTH_COOKIE,
    JSON.stringify({ token: pb.authStore.token, record: pb.authStore.record }),
    {
      httpOnly: true,
      sameSite: "lax",
      path: "/",
      maxAge: 7 * 24 * 60 * 60,
      secure: process.env.NODE_ENV === "production",
    },
  );
}

export async function clearAuth(): Promise<void> {
  (await cookies()).delete(AUTH_COOKIE);
}
