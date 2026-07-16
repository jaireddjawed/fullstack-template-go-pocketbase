"use client";

import Link from "next/link";
import { useActionState } from "react";
import { confirmPasswordReset, type FormState } from "@/lib/actions";

const initialState: FormState = {};

export default function ResetPasswordConfirmForm({ token }: { token: string }) {
	const [state, formAction, pending] = useActionState(confirmPasswordReset, initialState);

	return (
		<>
			<h1>Choose a new password</h1>
			<form className="stack" action={formAction}>
				<input type="hidden" name="token" value={token} />
				<input type="password" name="password" placeholder="New password" autoComplete="new-password" required />
				<input type="password" name="passwordConfirm" placeholder="Confirm new password" autoComplete="new-password" required />
				{state.error && <p className="error">{state.error}</p>}
				{state.success && <p className="success">{state.success}</p>}
				<button type="submit" disabled={pending}>Update password</button>
			</form>
			<p><Link href="/login">Back to log in</Link></p>
		</>
	);
}
