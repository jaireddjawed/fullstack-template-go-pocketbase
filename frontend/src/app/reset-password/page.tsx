"use client";

import Link from "next/link";
import { useActionState } from "react";
import { requestPasswordReset, type FormState } from "@/lib/actions";

const initialState: FormState = {};

export default function ResetPasswordPage() {
	const [state, formAction, pending] = useActionState(requestPasswordReset, initialState);

	return (
		<>
			<h1>Reset password</h1>
			<p className="muted">Enter your email and we’ll send a password-reset link.</p>
			<form className="stack" action={formAction}>
				<input type="email" name="email" placeholder="Email" autoComplete="email" required />
				{state.error && <p className="error">{state.error}</p>}
				{state.success && <p className="success">{state.success}</p>}
				<button type="submit" disabled={pending}>Send reset link</button>
			</form>
			<p>Remembered it? <Link href="/login">Log in</Link></p>
		</>
	);
}
