import Link from "next/link";
import VerifyEmailConfirmForm from "./verify-email-confirm-form";

export default async function VerifyEmailConfirmPage({ searchParams }: { searchParams: Promise<{ token?: string }> }) {
  const { token } = await searchParams;
  if (!token) return <><h1>Invalid verification link</h1><p>This link is missing its token.</p><Link href="/verify-email">Request a new link</Link></>;
  return <VerifyEmailConfirmForm token={token} />;
}
