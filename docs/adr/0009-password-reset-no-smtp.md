# ADR-0009: Password reset and email verification without SMTP

- **Status:** Accepted
- **Date:** 2026-07-22

## Context

Phase 6 needs account recovery and an `email_verified` flag. The product still
has **no email delivery** (same constraint as access-request notifications).
Shipping SMTP would add provider secrets, templates, and bounce handling before
we need them.

## Decision

1. **Password reset tokens** are stored hashed (`password_reset_tokens`), 1h TTL,
   one-time use, revoke-on-consume and revoke-all-on-user.
2. **Public** `POST /auth/forgot-password` always returns success (no enumeration).
   In **non-production**, the raw `reset_token` may be included in the JSON so
   local/dev flows work without an admin.
3. **Admin** `POST /admin/users/:id/password-reset` always returns the raw token
   for out-of-band delivery.
4. **`email_verified_at`** is set on register/bootstrap and toggleable by admins.
   There is no “send verification email” endpoint until SMTP exists.
5. Users complete reset via `POST /auth/reset-password` and `/reset-password` UI.

## Consequences

- Admins (or local forgot-password) must share reset tokens manually.
- Production forgot-password creates a token but does not expose it in the
  response — operators use the admin issue endpoint.
- Adding SMTP later can keep the same token table and simply email the link.
