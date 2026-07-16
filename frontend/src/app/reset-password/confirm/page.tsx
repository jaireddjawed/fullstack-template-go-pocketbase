import Link from "next/link";
import ResetPasswordConfirmForm from "./reset-password-confirm-form";

export default async function ResetPasswordConfirmPage({
	searchParams,
}: {
	searchParams: Promise<{ token?: string }>;
}) {
	const { token } = await searchParams;

	if (!token) {
		return (
			<>
				<h1>Invalid reset link</h1>
				<p>This password-reset link is missing its token.</p>
				<p><Link href="/reset-password">Request a new reset link</Link></p>
			</>
		);
	}

	return <ResetPasswordConfirmForm token={token} />;
}
