import { Link, router, usePage } from "@inertiajs/react";
import type { ReactNode } from "react";
import type { SharedProps } from "@/types";

export default function Layout({ children }: { children: ReactNode }) {
  const { auth } = usePage<SharedProps>().props;

  return (
    <>
      <nav>
        <Link href="/">Home</Link>
        {auth && <Link href="/posts">Posts</Link>}
        <span className="spacer" />
        {auth ? (
          <>
            <span className="muted">{auth.email}</span>
            <button onClick={() => router.post("/logout")}>Log out</button>
          </>
        ) : (
          <Link href="/login">Log in</Link>
        )}
      </nav>
      <main>{children}</main>
    </>
  );
}
