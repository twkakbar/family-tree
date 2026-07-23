package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/internal/middleware"
	"github.com/twkakbar/family-tree/backend/pkg/response"
)

// Stable machine-readable error codes. Registry: docs/api/README.md.
const (
	codeValidation         = "VALIDATION_ERROR"
	codeBadRequest         = "BAD_REQUEST"
	codeUnauthorized       = "UNAUTHORIZED"
	codeInvalidCredentials = "INVALID_CREDENTIALS"
	codeTokenInvalid       = "TOKEN_INVALID"
	codeTokenExpired       = "TOKEN_EXPIRED"
	codeForbidden          = "FORBIDDEN"
	codeNotFound           = "NOT_FOUND"
	codeConflict           = "CONFLICT"
	codeInternal           = "INTERNAL_ERROR"
)

// bindJSON binds and validates the request body. On failure it writes the
// 422/400 response and returns false; handlers simply return early.
func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			response.ErrorWithDetails(c, http.StatusUnprocessableEntity,
				codeValidation, "validation failed", fieldErrors(verrs))
			return false
		}
		response.Error(c, http.StatusBadRequest, codeBadRequest, "malformed request body")
		return false
	}
	return true
}

// fieldErrors converts validator errors into {snake_case_field: message}.
func fieldErrors(verrs validator.ValidationErrors) map[string]string {
	out := make(map[string]string, len(verrs))
	for _, fe := range verrs {
		out[toSnake(fe.Field())] = validationMessage(fe)
	}
	return out
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "strong_password":
		return "must be at least 8 characters and contain both letters and digits"
	default:
		return "is invalid"
	}
}

func toSnake(field string) string {
	var b strings.Builder
	for i, r := range field {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r - 'A' + 'a')
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// respondError maps domain errors to HTTP responses. It is the single place
// deciding status codes for expected failures; unexpected errors are logged
// with the request id and surface as an opaque 500.
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		response.Error(c, http.StatusNotFound, codeNotFound, "resource not found")
	case errors.Is(err, domain.ErrEmailTaken):
		response.Error(c, http.StatusConflict, codeConflict, "email is already registered")
	case errors.Is(err, domain.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, codeInvalidCredentials, "invalid email or password")
	case errors.Is(err, domain.ErrTokenExpired):
		response.Error(c, http.StatusUnauthorized, codeTokenExpired, "token has expired")
	case errors.Is(err, domain.ErrTokenInvalid):
		response.Error(c, http.StatusUnauthorized, codeTokenInvalid, "token is invalid")
	case errors.Is(err, domain.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, codeUnauthorized, "authentication required")
	case errors.Is(err, domain.ErrForbidden):
		response.Error(c, http.StatusForbidden, codeForbidden, "permission denied")
	case errors.Is(err, domain.ErrAlreadyCollaborator):
		response.Error(c, http.StatusConflict, codeConflict, "user is already a collaborator")
	case errors.Is(err, domain.ErrLastOwner):
		response.Error(c, http.StatusConflict, codeConflict, "family must retain at least one owner")
	case errors.Is(err, domain.ErrInvalidRelationship):
		response.Error(c, http.StatusUnprocessableEntity, codeValidation, "invalid relationship or person data")
	case errors.Is(err, domain.ErrRelationshipExists):
		response.Error(c, http.StatusConflict, codeConflict, "relationship already exists")
	case errors.Is(err, domain.ErrAccessRequestPending):
		response.Error(c, http.StatusConflict, codeConflict, "access request already pending")
	case errors.Is(err, domain.ErrAlreadyEditor):
		response.Error(c, http.StatusConflict, codeConflict, "user already has edit access")
	case errors.Is(err, domain.ErrInvalidUpload):
		response.Error(c, http.StatusUnprocessableEntity, codeValidation, "invalid photo upload (use JPEG, PNG, or WebP within size limit)")
	case errors.Is(err, domain.ErrUnsupportedFormat):
		response.Error(c, http.StatusUnprocessableEntity, codeValidation, "unsupported or invalid tree transfer document")
	default:
		slog.Error("unhandled error",
			"error", err,
			"request_id", middleware.RequestIDFromContext(c),
			"path", c.FullPath(),
		)
		response.Error(c, http.StatusInternalServerError, codeInternal, "internal server error")
	}
}
