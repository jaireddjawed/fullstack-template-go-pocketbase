import { Link, useForm, usePage } from "@inertiajs/react";
import type { FormEvent } from "react";
import Layout from "@/components/Layout";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function ConfirmEmail() {
  const { token } = usePage<{ token?: string }>().props;
  const { post, processing, errors } = useForm({ token: token ?? "" });
  function submit(e: FormEvent) { e.preventDefault(); post("/verify-email/confirm"); }
  return <Layout><Card className="mx-auto w-full max-w-md"><CardHeader><CardTitle>Verify your email</CardTitle><CardDescription>Confirm this email address to enable your account.</CardDescription></CardHeader><CardContent><form className="grid gap-4" onSubmit={submit}>{errors.token && <Alert variant="destructive"><AlertDescription>{errors.token}</AlertDescription></Alert>}<Button type="submit" disabled={processing || !token}>Verify email</Button></form><p className="mt-4 text-sm"><Link className="underline" href="/login">Back to log in</Link></p></CardContent></Card></Layout>;
}
