"use client";

import Link from "next/link";
import { useActionState } from "react";
import { signup, type FormState } from "@/lib/actions";

const initialState: FormState = {};

export default function SignupPage() {
	const [state, formAction, pending] = useActionState(signup, initialState);

	return (
		<>
			<h1>Create account</h1>
			<form className="stack" action={formAction}>
				<input type="text" name="name" placeholder="Name" autoComplete="name" />
				<input type="email" name="email" placeholder="Email" autoComplete="email" required />
				<input type="password" name="password" placeholder="Password" autoComplete="new-password" required />
				<input type="password" name="passwordConfirm" placeholder="Confirm password" autoComplete="new-password" required />
				{state.error && <p className="error">{state.error}</p>}
				{state.success && <p className="success">{state.success}</p>}
				<button type="submit" disabled={pending}>Create account</button>
			</form>
			<p>Already have an account? <Link href="/login">Log in</Link></p>
		</>
	);
}
