"use client";

import Link from "next/link";
import { useActionState, useTransition } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { signup, type FormState } from "@/lib/actions";

type SignupValues = { name: string; email: string; password: string; passwordConfirm: string };
const initialState: FormState = {};

export default function SignupPage() {
  const [state, dispatch] = useActionState(signup, initialState);
  const [pending, startTransition] = useTransition();
  const { register, handleSubmit, watch, formState: { errors } } = useForm<SignupValues>();

  function submit(values: SignupValues) {
    const formData = new FormData();
    for (const [key, value] of Object.entries(values)) formData.set(key, value);
    startTransition(() => dispatch(formData));
  }

  return (
    <Card className="mx-auto w-full max-w-md">
      <CardHeader><CardTitle>Create account</CardTitle><CardDescription>Create credentials for this PocketBase app.</CardDescription></CardHeader>
      <CardContent>
        <form className="grid gap-4" onSubmit={handleSubmit(submit)} noValidate>
          <div className="grid gap-2"><Label htmlFor="name">Name <span className="text-muted-foreground">(optional)</span></Label><Input id="name" autoComplete="name" {...register("name", { maxLength: { value: 100, message: "Name must be 100 characters or fewer." } })} />{errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}</div>
          <div className="grid gap-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" autoComplete="email" aria-invalid={Boolean(errors.email)} {...register("email", { required: "Email is required.", pattern: { value: /^\S+@\S+\.\S+$/, message: "Enter a valid email address." } })} />{errors.email && <p className="text-sm text-destructive">{errors.email.message}</p>}</div>
          <div className="grid gap-2"><Label htmlFor="password">Password</Label><Input id="password" type="password" autoComplete="new-password" aria-invalid={Boolean(errors.password)} {...register("password", { required: "Password is required.", minLength: { value: 8, message: "Password must be at least 8 characters." } })} />{errors.password && <p className="text-sm text-destructive">{errors.password.message}</p>}</div>
          <div className="grid gap-2"><Label htmlFor="passwordConfirm">Confirm password</Label><Input id="passwordConfirm" type="password" autoComplete="new-password" aria-invalid={Boolean(errors.passwordConfirm)} {...register("passwordConfirm", { required: "Confirm your password.", validate: (value) => value === watch("password") || "Passwords do not match." })} />{errors.passwordConfirm && <p className="text-sm text-destructive">{errors.passwordConfirm.message}</p>}</div>
          {state.error && <Alert variant="destructive"><AlertDescription>{state.error}</AlertDescription></Alert>}
          {state.success && <Alert><AlertDescription>{state.success}</AlertDescription></Alert>}
          <Button type="submit" disabled={pending}>{pending ? "Creating account…" : "Create account"}</Button>
        </form>
        <p className="mt-4 text-sm">Already have an account? <Link className="underline" href="/login">Log in</Link></p>
      </CardContent>
    </Card>
  );
}
