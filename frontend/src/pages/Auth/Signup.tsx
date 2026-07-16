import { Link, useForm } from "@inertiajs/react";
import type { FormEvent } from "react";
import Layout from "@/components/Layout";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function Signup() {
  const { data, setData, post, processing, errors } = useForm({ name: "", email: "", password: "", passwordConfirm: "" });
  function submit(e: FormEvent) { e.preventDefault(); post("/signup"); }
  return <Layout><Card className="mx-auto w-full max-w-md"><CardHeader><CardTitle>Create account</CardTitle><CardDescription>Verify your email before signing in.</CardDescription></CardHeader><CardContent><form className="grid gap-4" onSubmit={submit}><div className="grid gap-2"><Label htmlFor="name">Name</Label><Input id="name" value={data.name} onChange={(e) => setData("name", e.target.value)} /></div><div className="grid gap-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" value={data.email} onChange={(e) => setData("email", e.target.value)} /></div><div className="grid gap-2"><Label htmlFor="password">Password</Label><Input id="password" type="password" value={data.password} onChange={(e) => setData("password", e.target.value)} /></div><div className="grid gap-2"><Label htmlFor="passwordConfirm">Confirm password</Label><Input id="passwordConfirm" type="password" value={data.passwordConfirm} onChange={(e) => setData("passwordConfirm", e.target.value)} /></div>{errors.email && <Alert variant="destructive"><AlertDescription>{errors.email}</AlertDescription></Alert>}<Button type="submit" disabled={processing}>Create account</Button></form><p className="mt-4 text-sm">Already have an account? <Link className="underline" href="/login">Log in</Link></p></CardContent></Card></Layout>;
}
