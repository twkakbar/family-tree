# ADR-0012: Require email verification when SMTP is configured

- **Status:** Accepted
- **Date:** 2026-10-11

## Context

Accounts were marked verified at registration because there was no mail
delivery (ADR-0009). The product now has an SMTP relay and wants new
editors to confirm their address before they can sign in.

## Decision

1. When `SMTP_HOST`, credentials, `SMTP_FROM`, and `APP_PUBLIC_URL` are all
   set, registration creates an unverified user and emails a 24-hour link.
   No session is issued. Login with a correct password returns
   `EMAIL_NOT_VERIFIED` until the link is consumed.
2. Each address may receive at most **3** verification emails in a rolling
   hour (`email_dispatches`). The fourth returns `RATE_LIMITED`.
3. When SMTP is not configured, registration still auto-verifies and signs
   in. CI and laptops without mail keep the previous behavior.
4. Password reset and access-request mail stay out of band (ADR-0009).
   The bootstrap admin remains verified.

## Consequences

- Production must keep SMTP settings only in the VM `.env`, never in git.
- A failed send deletes a brand-new unverified account so the person can
  try again. An existing unverified account is left in place.
- The From address must be one the relay accepts.
