"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import { ClientResponseError } from "pocketbase";
import type { PublishPostResponse } from "@shared/types.gen";
import { POCKETBASE_URL } from "./pocketbase";
import { clearAuth, createServerClient, persistAuth } from "./server-auth";

export interface LoginState {
  error?: string;
}

/** Server action: authenticate against PocketBase and persist the session. */
export async function login(
  _prev: LoginState,
  formData: FormData,
): Promise<LoginState> {
  const email = String(formData.get("email") ?? "");
  const password = String(formData.get("password") ?? "");

  const pb = await createServerClient();

  try {
    await pb.collection("users").authWithPassword(email, password);
  } catch (err) {
    if (err instanceof ClientResponseError) {
      return { error: "Invalid email or password." };
    }
    throw err;
  }

  await persistAuth(pb);
  redirect("/posts");
}

export async function logout(): Promise<void> {
  await clearAuth();
  redirect("/");
}

/** Server action: call the custom Go endpoint with the user's token. */
export async function publishPost(id: string): Promise<PublishPostResponse> {
  const pb = await createServerClient();

  const res = await fetch(`${POCKETBASE_URL}/api/app/posts/${id}/publish`, {
    method: "POST",
    headers: { Authorization: pb.authStore.token },
  });

  if (!res.ok) {
    throw new Error(`publish failed: ${res.status}`);
  }

  revalidatePath("/posts");
  return (await res.json()) as PublishPostResponse;
}
