// Package domain contains the enterprise core of the application: entities,
// domain errors, and the repository/service contracts implemented by the
// outer layers. It must not import from any other internal package.
package domain

import "errors"

// Sentinel errors returned by services for expected failure modes. Handlers
// map them to HTTP status codes centrally; repositories translate driver
// errors into them. Anything not in this list is treated as an internal
// error (HTTP 500).
var (
	// ErrNotFound signals that the requested resource does not exist
	// (or is soft-deleted).
	ErrNotFound = errors.New("resource not found")

	// ErrEmailTaken signals a registration/update with an email that is
	// already in use by an active account.
	ErrEmailTaken = errors.New("email already registered")

	// ErrInvalidCredentials signals a failed login. It deliberately does not
	// distinguish "unknown email" from "wrong password".
	ErrInvalidCredentials = errors.New("invalid email or password")

	// ErrTokenInvalid signals a missing, malformed, unknown or reused token.
	ErrTokenInvalid = errors.New("token is invalid")

	// ErrTokenExpired signals a well-formed but expired token.
	ErrTokenExpired = errors.New("token has expired")

	// ErrUnauthorized signals a request lacking valid authentication.
	ErrUnauthorized = errors.New("authentication required")

	// ErrForbidden signals an authenticated request lacking permission.
	ErrForbidden = errors.New("permission denied")
)
