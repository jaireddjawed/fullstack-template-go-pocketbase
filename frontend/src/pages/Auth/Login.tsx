import { useForm } from "@inertiajs/react";
import type { FormEvent } from "react";
import Layout from "@/components/Layout";

export default function Login() {
  const { data, setData, post, processing, errors } = useForm({
    email: "demo@example.com",
    password: "password123",
  });

  function submit(e: FormEvent) {
    e.preventDefault();
    post("/login");
  }

  return (
    <Layout>
      <h1>Log in</h1>
      <form className="stack" onSubmit={submit}>
        <input
          type="email"
          placeholder="Email"
          value={data.email}
          onChange={(e) => setData("email", e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          value={data.password}
          onChange={(e) => setData("password", e.target.value)}
        />
        {errors.email && <p className="error">{errors.email}</p>}
        <button type="submit" disabled={processing}>
          Log in
        </button>
      </form>
      <p className="muted">
        Defaults are the seeded demo account (<code>make seed</code>).
      </p>
    </Layout>
  );
}
