"use client";

import Link from "next/link";
import { useActionState, useTransition } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { confirmPasswordReset, type FormState } from "@/lib/actions";

type ResetValues = { password: string; passwordConfirm: string };
const initialState: FormState = {};

export default function ResetPasswordConfirmForm({ token }: { token: string }) {
  const [state, dispatch] = useActionState(confirmPasswordReset, initialState);
  const [pending, startTransition] = useTransition();
  const { register, handleSubmit, watch, formState: { errors } } = useForm<ResetValues>();

  function submit(values: ResetValues) {
    const formData = new FormData();
    formData.set("token", token);
    formData.set("password", values.password);
    formData.set("passwordConfirm", values.passwordConfirm);
    startTransition(() => dispatch(formData));
  }

  return (
    <Card className="mx-auto w-full max-w-md">
      <CardHeader><CardTitle>Choose a new password</CardTitle><CardDescription>Your new password must be at least 8 characters.</CardDescription></CardHeader>
      <CardContent>
        <form className="grid gap-4" onSubmit={handleSubmit(submit)} noValidate>
          <div className="grid gap-2"><Label htmlFor="password">New password</Label><Input id="password" type="password" autoComplete="new-password" aria-invalid={Boolean(errors.password)} {...register("password", { required: "Password is required.", minLength: { value: 8, message: "Password must be at least 8 characters." } })} />{errors.password && <p className="text-sm text-destructive">{errors.password.message}</p>}</div>
          <div className="grid gap-2"><Label htmlFor="passwordConfirm">Confirm password</Label><Input id="passwordConfirm" type="password" autoComplete="new-password" aria-invalid={Boolean(errors.passwordConfirm)} {...register("passwordConfirm", { required: "Confirm your password.", validate: (value) => value === watch("password") || "Passwords do not match." })} />{errors.passwordConfirm && <p className="text-sm text-destructive">{errors.passwordConfirm.message}</p>}</div>
          {state.error && <Alert variant="destructive"><AlertDescription>{state.error}</AlertDescription></Alert>}
          {state.success && <Alert><AlertDescription>{state.success}</AlertDescription></Alert>}
          <Button type="submit" disabled={pending}>{pending ? "Updating…" : "Update password"}</Button>
        </form>
        <p className="mt-4 text-sm"><Link className="underline" href="/login">Back to log in</Link></p>
      </CardContent>
    </Card>
  );
}
