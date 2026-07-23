# Shared security headers for every response that uses `include`.
# Note: nginx does not inherit add_header into locations that set their own
# add_header — those locations must include this file again.

add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "DENY" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Permissions-Policy "camera=(), microphone=(), geolocation=(), payment=()" always;
add_header Cross-Origin-Opener-Policy "same-origin" always;
add_header Cross-Origin-Resource-Policy "same-origin" always;

# Content-Security-Policy (Phase 7 / ADR-0010).
# Pragmatic allowances for Next.js App Router + React Three Fiber:
#   - 'unsafe-inline' styles (Tailwind / Next)
#   - 'unsafe-inline' + 'unsafe-eval' + wasm-unsafe-eval scripts (Next + Three/WebGL)
#   - blob: workers (Three.js)
# Tighten further when nonce-based Next CSP is practical.
add_header Content-Security-Policy "default-src 'self'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'; object-src 'none'; img-src 'self' data: blob:; font-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' 'wasm-unsafe-eval'; worker-src 'self' blob:; connect-src 'self'; media-src 'self'" always;
