import { NextResponse, type NextRequest } from "next/server";

/**
 * UX-level route guard based on the presence of the refresh cookie.
 * It cannot (and does not need to) validate the token - the API enforces
 * real authentication. This only avoids rendering pages that would
 * immediately bounce the user.
 */
const REFRESH_COOKIE = "ft_refresh";
const AUTH_PAGES = new Set(["/login", "/register"]);

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const hasSessionCookie = request.cookies.has(REFRESH_COOKIE);

  if (pathname.startsWith("/dashboard") && !hasSessionCookie) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("next", pathname);
    return NextResponse.redirect(loginUrl);
  }

  if (AUTH_PAGES.has(pathname) && hasSessionCookie) {
    return NextResponse.redirect(new URL("/dashboard", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/dashboard/:path*", "/login", "/register"],
};
