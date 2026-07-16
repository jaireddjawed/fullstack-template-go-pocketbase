import { useForm } from "@inertiajs/react";
import type { FormEvent } from "react";
import Layout from "@/components/Layout";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function Login() {
  const { data, setData, setError, clearErrors, post, processing, errors } = useForm({
    email: "demo@example.com",
    password: "password123",
  });

  function submit(e: FormEvent) {
    e.preventDefault();
    clearErrors();
    if (!data.email.trim()) return setError("email", "Email is required.");
    if (!/^\S+@\S+\.\S+$/.test(data.email)) return setError("email", "Enter a valid email address.");
    if (!data.password) return setError("password", "Password is required.");
    post("/login");
  }

  return (
    <Layout>
      <Card className="mx-auto w-full max-w-md">
        <CardHeader><CardTitle>Log in</CardTitle><CardDescription>Use the seeded account or your own credentials.</CardDescription></CardHeader>
        <CardContent>
          <form className="grid gap-4" onSubmit={submit} noValidate>
            <div className="grid gap-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" autoComplete="email" aria-invalid={Boolean(errors.email)} value={data.email} onChange={(e) => setData("email", e.target.value)} /></div>
            <div className="grid gap-2"><Label htmlFor="password">Password</Label><Input id="password" type="password" autoComplete="current-password" aria-invalid={Boolean(errors.password)} value={data.password} onChange={(e) => setData("password", e.target.value)} /></div>
            {(errors.email || errors.password) && <Alert variant="destructive"><AlertDescription>{errors.email || errors.password}</AlertDescription></Alert>}
            <Button type="submit" disabled={processing}>{processing ? "Logging in…" : "Log in"}</Button>
          </form>
          <p className="mt-4 text-sm text-muted-foreground">Defaults are the seeded demo account (<code>make seed</code>).</p>
        </CardContent>
      </Card>
    </Layout>
  );
}
