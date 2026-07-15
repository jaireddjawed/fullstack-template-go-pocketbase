import { NextResponse, type NextRequest } from "next/server";
import { AUTH_COOKIE } from "@/lib/pocketbase";

// Presence check only — the token itself is validated by PocketBase on
// every request that uses it.
export function middleware(request: NextRequest) {
  if (!request.cookies.get(AUTH_COOKIE)?.value) {
    const login = new URL("/login", request.url);
    return NextResponse.redirect(login);
  }
  return NextResponse.next();
}

export const config = {
  matcher: ["/posts/:path*"],
};
