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

export interface FormState {
	error?: string;
	success?: string;
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

  if (!pb.authStore.record?.verified) {
    pb.authStore.clear();
    return { error: "Verify your email before logging in. You can request a new link below." };
  }

  await persistAuth(pb);
  redirect("/posts");
}

export async function logout(): Promise<void> {
	await clearAuth();
	redirect("/");
}

/** Server action: create a new PocketBase auth record. */
export async function signup(
	_prev: FormState,
	formData: FormData,
): Promise<FormState> {
	const name = String(formData.get("name") ?? "").trim();
	const email = String(formData.get("email") ?? "").trim();
	const password = String(formData.get("password") ?? "");
	const passwordConfirm = String(formData.get("passwordConfirm") ?? "");

	if (password !== passwordConfirm) {
		return { error: "Passwords do not match." };
	}

	try {
		const pb = await createServerClient();
		await pb.collection("users").create({
			name,
			email,
			password,
			passwordConfirm,
		});
		await pb.collection("users").requestVerification(email);
	} catch (err) {
		if (err instanceof ClientResponseError) {
			return { error: "We couldn't create that account. Check your details and try again." };
		}
		throw err;
	}

	return { success: "Account created. Check your email to verify it before logging in." };
}

/** Server action: ask PocketBase to email an account-verification link. */
export async function requestEmailVerification(
	_prev: FormState,
	formData: FormData,
): Promise<FormState> {
	const email = String(formData.get("email") ?? "").trim();

	try {
		await (await createServerClient()).collection("users").requestVerification(email);
	} catch (err) {
		if (!(err instanceof ClientResponseError)) throw err;
	}

	return { success: "If an account exists for that email, a verification link is on its way." };
}

/** Server action: verify an email using PocketBase's emailed token. */
export async function confirmEmailVerification(
	_prev: FormState,
	formData: FormData,
): Promise<FormState> {
	const token = String(formData.get("token") ?? "");
	if (!token) return { error: "This verification link is missing its token." };

	try {
		await (await createServerClient()).collection("users").confirmVerification(token);
	} catch (err) {
		if (err instanceof ClientResponseError) return { error: "This verification link is invalid or has expired." };
		throw err;
	}

	return { success: "Email verified. You can now log in." };
}

/** Server action: ask PocketBase to email a password-reset link. */
export async function requestPasswordReset(
	_prev: FormState,
	formData: FormData,
): Promise<FormState> {
	const email = String(formData.get("email") ?? "").trim();

	try {
		await (await createServerClient()).collection("users").requestPasswordReset(email);
	} catch (err) {
		if (!(err instanceof ClientResponseError)) {
			throw err;
		}
	}

	// Do not reveal whether an account exists for the submitted address.
	return {
		success: "If an account exists for that email, a password-reset link is on its way.",
	};
}

/** Server action: set a new password using the reset token from PocketBase. */
export async function confirmPasswordReset(
	_prev: FormState,
	formData: FormData,
): Promise<FormState> {
	const token = String(formData.get("token") ?? "");
	const password = String(formData.get("password") ?? "");
	const passwordConfirm = String(formData.get("passwordConfirm") ?? "");

	if (!token) {
		return { error: "This password-reset link is missing its token." };
	}
	if (password !== passwordConfirm) {
		return { error: "Passwords do not match." };
	}

	try {
		await (await createServerClient())
			.collection("users")
			.confirmPasswordReset(token, password, passwordConfirm);
	} catch (err) {
		if (err instanceof ClientResponseError) {
			return { error: "This password-reset link is invalid or has expired." };
		}
		throw err;
	}

	return { success: "Password updated. You can now log in." };
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
