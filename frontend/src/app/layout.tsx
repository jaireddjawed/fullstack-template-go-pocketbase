import type { Metadata } from "next";
import Link from "next/link";
import { currentUser } from "@/lib/server-auth";
import { logout } from "@/lib/actions";
import "./globals.css";

export const metadata: Metadata = {
  title: "Fullstack Template",
  description: "Next.js + PocketBase (extended with Go)",
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const user = await currentUser();

  return (
    <html lang="en">
      <body>
        <nav>
          <Link href="/">Home</Link>
          {user && <Link href="/posts">Posts</Link>}
          <span className="spacer" />
          {user ? (
            <>
              <span className="muted">{user.email}</span>
              <form action={logout}>
                <button type="submit">Log out</button>
              </form>
            </>
          ) : (
            <>
              <Link href="/login">Log in</Link>
              <Link href="/signup">Sign up</Link>
            </>
          )}
        </nav>
        <main>{children}</main>
      </body>
    </html>
  );
}
