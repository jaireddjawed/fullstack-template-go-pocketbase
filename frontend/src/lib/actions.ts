"use server";

import { revalidatePath } from "next/cache";
import type { PublishPostResponse } from "@shared/types.gen";
import { POCKETBASE_URL } from "./pocketbase";
import { clerkToken } from "./server-auth";

/** Server action: call the custom Go endpoint as the Clerk user. */
export async function publishPost(id: string): Promise<PublishPostResponse> {
  const token = await clerkToken();

  const res = await fetch(`${POCKETBASE_URL}/api/app/posts/${id}/publish`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!res.ok) {
    throw new Error(`publish failed: ${res.status}`);
  }

  revalidatePath("/posts");
  return (await res.json()) as PublishPostResponse;
}
