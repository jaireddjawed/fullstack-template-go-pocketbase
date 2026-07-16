"use client";

import Link from "next/link";
import { useActionState, useTransition } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { requestEmailVerification, type FormState } from "@/lib/actions";

type VerificationValues = { email: string };
const initialState: FormState = {};

export default function VerifyEmailPage() {
  const [state, dispatch] = useActionState(requestEmailVerification, initialState);
  const [pending, startTransition] = useTransition();
  const { register, handleSubmit, formState: { errors } } = useForm<VerificationValues>();

  function submit(values: VerificationValues) {
    const formData = new FormData();
    formData.set("email", values.email);
    startTransition(() => dispatch(formData));
  }

  return <Card className="mx-auto w-full max-w-md"><CardHeader><CardTitle>Verify your email</CardTitle><CardDescription>Enter your email to receive a new verification link.</CardDescription></CardHeader><CardContent><form className="grid gap-4" onSubmit={handleSubmit(submit)} noValidate><div className="grid gap-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" autoComplete="email" aria-invalid={Boolean(errors.email)} {...register("email", { required: "Email is required.", pattern: { value: /^\S+@\S+\.\S+$/, message: "Enter a valid email address." } })} />{errors.email && <p className="text-sm text-destructive">{errors.email.message}</p>}</div>{state.error && <Alert variant="destructive"><AlertDescription>{state.error}</AlertDescription></Alert>}{state.success && <Alert><AlertDescription>{state.success}</AlertDescription></Alert>}<Button type="submit" disabled={pending}>{pending ? "Sending…" : "Send verification link"}</Button></form><p className="mt-4 text-sm"><Link className="underline" href="/login">Back to log in</Link></p></CardContent></Card>;
}
