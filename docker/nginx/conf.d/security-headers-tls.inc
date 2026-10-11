# HTTPS-only additions. Include alongside security-headers.inc in the
# same server or location context (nginx sends both).
# HSTS is safe only after the certificate for this host is valid.

add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
add_header Content-Security-Policy "default-src 'self'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'; object-src 'none'; img-src 'self' data: blob:; font-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' 'wasm-unsafe-eval'; worker-src 'self' blob:; connect-src 'self'; media-src 'self'; upgrade-insecure-requests" always;
