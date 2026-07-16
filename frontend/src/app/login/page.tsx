"use client";

import Link from "next/link";
import { useActionState, useTransition } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { login, type LoginState } from "@/lib/actions";

type LoginValues = { email: string; password: string };
const initialState: LoginState = {};

export default function LoginPage() {
  const [state, dispatch] = useActionState(login, initialState);
  const [pending, startTransition] = useTransition();
  const { register, handleSubmit, formState: { errors } } = useForm<LoginValues>({
    defaultValues: { email: "demo@example.com", password: "password123" },
  });

  function submit(values: LoginValues) {
    const formData = new FormData();
    formData.set("email", values.email);
    formData.set("password", values.password);
    startTransition(() => dispatch(formData));
  }

  return (
    <Card className="mx-auto w-full max-w-md">
      <CardHeader>
        <CardTitle>Log in</CardTitle>
        <CardDescription>Use the seeded account or your own credentials.</CardDescription>
      </CardHeader>
      <CardContent>
        <form className="grid gap-4" onSubmit={handleSubmit(submit)} noValidate>
          <div className="grid gap-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" autoComplete="email" aria-invalid={Boolean(errors.email)} {...register("email", { required: "Email is required.", pattern: { value: /^\S+@\S+\.\S+$/, message: "Enter a valid email address." } })} />
            {errors.email && <p className="text-sm text-destructive">{errors.email.message}</p>}
          </div>
          <div className="grid gap-2">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" autoComplete="current-password" aria-invalid={Boolean(errors.password)} {...register("password", { required: "Password is required." })} />
            {errors.password && <p className="text-sm text-destructive">{errors.password.message}</p>}
          </div>
          {state.error && <Alert variant="destructive"><AlertDescription>{state.error}</AlertDescription></Alert>}
          <Button type="submit" disabled={pending}>{pending ? "Logging in…" : "Log in"}</Button>
        </form>
        <p className="mt-4 text-sm text-muted-foreground">Defaults: <code>demo@example.com</code> / <code>password123</code></p>
        <p className="mt-4 text-sm">Need an account? <Link className="underline" href="/signup">Sign up</Link>{" · "}<Link className="underline" href="/verify-email">Verify email</Link>{" · "}<Link className="underline" href="/reset-password">Forgot your password?</Link></p>
      </CardContent>
    </Card>
  );
}
