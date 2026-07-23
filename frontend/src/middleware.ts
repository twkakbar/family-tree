import { NextResponse, type NextRequest } from "next/server";

/**
 * UX-level route guard based on the presence of the refresh cookie.
 * Real auth is always enforced by the API.
 */
const REFRESH_COOKIE = "ft_refresh";
const AUTH_PAGES = new Set(["/login", "/register"]);
const PROTECTED_PREFIXES = ["/dashboard", "/studio", "/request-access", "/admin"];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const hasSessionCookie = request.cookies.has(REFRESH_COOKIE);

  const needsAuth = PROTECTED_PREFIXES.some(
    (prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`),
  );

  if (needsAuth && !hasSessionCookie) {
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
  matcher: [
    "/dashboard/:path*",
    "/studio/:path*",
    "/request-access/:path*",
    "/admin/:path*",
    "/login",
    "/register",
    "/forgot-password",
    "/reset-password",
  ],
};
