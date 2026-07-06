// Package response defines the JSON envelope used by every API endpoint:
//
//	{"success": true,  "data": …, "meta": …}
//	{"success": false, "error": {"code": …, "message": …, "details": …}}
//
// It is framework-thin (gin helpers only) and domain-agnostic; mapping
// domain errors to HTTP codes happens in the handler layer.
package response

import "github.com/gin-gonic/gin"

// Envelope is the top-level shape of every API response.
type Envelope struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
	Meta    any       `json:"meta,omitempty"`
}

// APIError is the machine-readable error payload. Code values are stable and
// registered in docs/api/README.md.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// OK writes a success envelope with the given status and payload.
func OK(c *gin.Context, status int, data any) {
	c.JSON(status, Envelope{Success: true, Data: data})
}

// OKWithMeta writes a success envelope including metadata (e.g. pagination).
func OKWithMeta(c *gin.Context, status int, data, meta any) {
	c.JSON(status, Envelope{Success: true, Data: data, Meta: meta})
}

// Error writes an error envelope and aborts the request chain.
func Error(c *gin.Context, status int, code, message string) {
	ErrorWithDetails(c, status, code, message, nil)
}

// ErrorWithDetails writes an error envelope with structured details
// (e.g. per-field validation messages) and aborts the request chain.
func ErrorWithDetails(c *gin.Context, status int, code, message string, details any) {
	c.AbortWithStatusJSON(status, Envelope{
		Success: false,
		Error:   &APIError{Code: code, Message: message, Details: details},
	})
}
