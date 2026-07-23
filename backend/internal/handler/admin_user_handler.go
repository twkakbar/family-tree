package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/internal/middleware"
	"github.com/twkakbar/family-tree/backend/pkg/pagination"
	"github.com/twkakbar/family-tree/backend/pkg/response"
)

// AdminUserHandler serves /api/v1/admin/users.
type AdminUserHandler struct {
	users domain.AdminUserService
}

// NewAdminUserHandler builds admin user endpoints.
func NewAdminUserHandler(users domain.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{users: users}
}

type adminUpdateUserRequest struct {
	Name          *string `json:"name"`
	Role          *string `json:"role"`
	EmailVerified *bool   `json:"email_verified"`
}

// ListUsers handles GET /api/v1/admin/users.
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	params := pagination.Parse(c.Query("page"), c.Query("per_page"))
	items, total, err := h.users.List(c.Request.Context(), adminID, c.Query("search"), params.Page, params.PerPage)
	if err != nil {
		respondError(c, err)
		return
	}
	out := make([]UserResponse, len(items))
	for i := range items {
		out[i] = toUserResponse(&items[i])
	}
	response.OKWithMeta(c, http.StatusOK, out, pagination.NewMeta(params, total))
}

// UpdateUser handles PATCH /api/v1/admin/users/:id.
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid user id")
		return
	}
	var req adminUpdateUserRequest
	if !bindJSON(c, &req) {
		return
	}
	in := domain.AdminUpdateUserInput{
		Name:          req.Name,
		EmailVerified: req.EmailVerified,
	}
	if req.Role != nil {
		role := domain.Role(*req.Role)
		in.Role = &role
	}
	user, err := h.users.Update(c.Request.Context(), adminID, userID, in)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toUserResponse(user))
}

// DeleteUser handles DELETE /api/v1/admin/users/:id.
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid user id")
		return
	}
	if err := h.users.Delete(c.Request.Context(), adminID, userID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{"deleted": true})
}

// IssuePasswordReset handles POST /api/v1/admin/users/:id/password-reset.
func (h *AdminUserHandler) IssuePasswordReset(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid user id")
		return
	}
	raw, err := h.users.IssuePasswordReset(c.Request.Context(), adminID, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{
		"reset_token": raw,
		"expires_in":  3600,
		"hint":        "Share /reset-password?token=… with the user out-of-band (no email delivery yet).",
	})
}
