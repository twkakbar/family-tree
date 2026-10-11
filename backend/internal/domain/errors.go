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

	// ErrAlreadyCollaborator signals an invite for a user who already has
	// access to the family.
	ErrAlreadyCollaborator = errors.New("user is already a collaborator")

	// ErrLastOwner signals an attempt to remove or demote the sole owner of
	// a family.
	ErrLastOwner = errors.New("family must retain at least one owner")

	// ErrInvalidRelationship signals a relationship that violates tree
	// invariants (wrong family, cycle, too many parents, self-edge).
	ErrInvalidRelationship = errors.New("invalid relationship")

	// ErrRelationshipExists signals a duplicate edge of the same type.
	ErrRelationshipExists = errors.New("relationship already exists")

	// ErrAccessRequestPending signals the user already has a pending request.
	ErrAccessRequestPending = errors.New("access request already pending")

	// ErrAlreadyEditor signals the user already has edit rights on the royal family.
	ErrAlreadyEditor = errors.New("user already has edit access")

	// ErrInvalidUpload signals a rejected photo (type, size, or empty file).
	ErrInvalidUpload = errors.New("invalid upload")

	// ErrUnsupportedFormat signals an export/import document the server rejects.
	ErrUnsupportedFormat = errors.New("unsupported transfer format")

	// ErrEmailNotVerified signals a correct password for an account that has
	// not confirmed its email yet.
	ErrEmailNotVerified = errors.New("email is not verified")

	// ErrRateLimited signals too many verification emails in the current hour.
	ErrRateLimited = errors.New("too many verification emails")
)
