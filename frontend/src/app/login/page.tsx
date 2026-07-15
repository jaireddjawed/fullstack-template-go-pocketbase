"use client";

import { useActionState } from "react";
import { login, type LoginState } from "@/lib/actions";

const initialState: LoginState = {};

export default function LoginPage() {
  const [state, formAction, pending] = useActionState(login, initialState);

  return (
    <>
      <h1>Log in</h1>
      <form className="stack" action={formAction}>
        <input
          type="email"
          name="email"
          placeholder="Email"
          defaultValue="demo@example.com"
          required
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          defaultValue="password123"
          required
        />
        {state.error && <p className="error">{state.error}</p>}
        <button type="submit" disabled={pending}>
          Log in
        </button>
      </form>
      <p className="muted">
        Defaults are the seeded demo account (<code>make seed</code>).
      </p>
    </>
  );
}
