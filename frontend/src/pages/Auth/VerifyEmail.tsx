import { Link, useForm, usePage } from "@inertiajs/react";
import type { FormEvent } from "react";
import Layout from "@/components/Layout";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function VerifyEmail() {
  const { data, setData, post, processing } = useForm({ email: "" });
  const { success } = usePage<{ success?: string }>().props;
  function submit(e: FormEvent) { e.preventDefault(); post("/verify-email"); }
  return <Layout><Card className="mx-auto w-full max-w-md"><CardHeader><CardTitle>Verify your email</CardTitle><CardDescription>Request a new verification link.</CardDescription></CardHeader><CardContent><form className="grid gap-4" onSubmit={submit}><div className="grid gap-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" value={data.email} onChange={(e) => setData("email", e.target.value)} /></div>{success && <Alert><AlertDescription>{success}</AlertDescription></Alert>}<Button type="submit" disabled={processing}>Send verification link</Button></form><p className="mt-4 text-sm"><Link className="underline" href="/login">Back to log in</Link></p></CardContent></Card></Layout>;
}
