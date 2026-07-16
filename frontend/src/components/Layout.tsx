import { Link, router, usePage } from "@inertiajs/react";
import type { ReactNode } from "react";
import type { SharedProps } from "@/types";
import { Button } from "@/components/ui/button";

export default function Layout({ children }: { children: ReactNode }) {
  const { auth } = usePage<SharedProps>().props;

  return (
    <>
      <nav className="mb-8 flex items-center gap-4 border-b pb-4">
        <Link href="/">Home</Link>
        {auth && <Link href="/posts">Posts</Link>}
        <span className="flex-1" />
        {auth ? (
          <>
            <span className="text-sm text-muted-foreground">{auth.email}</span>
            <Button variant="outline" onClick={() => router.post("/logout")}>Log out</Button>
          </>
        ) : (
          <><Link href="/login">Log in</Link><Link href="/signup">Sign up</Link></>
        )}
      </nav>
      <main>{children}</main>
    </>
  );
}
