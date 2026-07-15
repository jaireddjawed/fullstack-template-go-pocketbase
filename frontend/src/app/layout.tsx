import type { Metadata } from "next";
import Link from "next/link";
import { ClerkProvider, SignInButton, UserButton } from "@clerk/nextjs";
import { auth } from "@clerk/nextjs/server";
import "./globals.css";

export const metadata: Metadata = {
  title: "Fullstack Template",
  description: "Next.js + Clerk + PocketBase (extended with Go)",
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { userId } = await auth();

  return (
    <ClerkProvider>
      <html lang="en">
        <body>
          <nav>
            <Link href="/">Home</Link>
            {userId && <Link href="/posts">Posts</Link>}
            <span className="spacer" />
            {userId ? (
              <UserButton />
            ) : (
              <SignInButton mode="modal">
                <button>Sign in</button>
              </SignInButton>
            )}
          </nav>
          <main>{children}</main>
        </body>
      </html>
    </ClerkProvider>
  );
}
