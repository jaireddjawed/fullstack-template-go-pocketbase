"use client";

import Link from "next/link";
import { useActionState, useTransition } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { confirmEmailVerification, type FormState } from "@/lib/actions";

const initialState: FormState = {};

export default function VerifyEmailConfirmForm({ token }: { token: string }) {
  const [state, dispatch] = useActionState(confirmEmailVerification, initialState);
  const [pending, startTransition] = useTransition();
  function submit() { const formData = new FormData(); formData.set("token", token); startTransition(() => dispatch(formData)); }
  return <Card className="mx-auto w-full max-w-md"><CardHeader><CardTitle>Verify your email</CardTitle><CardDescription>Confirm this email address to enable your account.</CardDescription></CardHeader><CardContent><form onSubmit={(event) => { event.preventDefault(); submit(); }} className="grid gap-4">{state.error && <Alert variant="destructive"><AlertDescription>{state.error}</AlertDescription></Alert>}{state.success && <Alert><AlertDescription>{state.success}</AlertDescription></Alert>}<Button type="submit" disabled={pending}>{pending ? "Verifying…" : "Verify email"}</Button></form><p className="mt-4 text-sm"><Link className="underline" href="/login">Back to log in</Link></p></CardContent></Card>;
}
