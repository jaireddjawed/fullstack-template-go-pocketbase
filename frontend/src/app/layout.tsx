import type { Metadata } from "next";
import Link from "next/link";
import { ClerkProvider, SignInButton, UserButton } from "@clerk/nextjs";
import { auth } from "@clerk/nextjs/server";
import "./globals.css";
import { Geist } from "next/font/google";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

const geist = Geist({subsets:['latin'],variable:'--font-sans'});

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
      <html lang="en" className={cn("font-sans", geist.variable)}>
        <body className="mx-auto min-h-screen max-w-4xl bg-background px-6 py-8 text-foreground">
          <nav className="mb-8 flex items-center gap-4 border-b pb-4">
            <Link href="/">Home</Link>
            {userId && <Link href="/posts">Posts</Link>}
            <span className="flex-1" />
            {userId ? (
              <UserButton />
            ) : (
              <SignInButton mode="modal">
                <Button>Sign in</Button>
              </SignInButton>
            )}
          </nav>
          <main>{children}</main>
        </body>
      </html>
    </ClerkProvider>
  );
}
