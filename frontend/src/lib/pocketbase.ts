import PocketBase from "pocketbase";
import type { TypedPocketBase } from "@shared/pocketbase.gen";

export const POCKETBASE_URL =
  process.env.NEXT_PUBLIC_POCKETBASE_URL ?? "http://127.0.0.1:8090";

export const AUTH_COOKIE = "pb_auth";

/**
 * A fresh SDK client. On the server, create one per request and hydrate it
 * from the auth cookie (see server-auth.ts) — never share a client between
 * requests, the auth store is stateful.
 */
export function createClient(): TypedPocketBase {
  return new PocketBase(POCKETBASE_URL) as TypedPocketBase;
}
