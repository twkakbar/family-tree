package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/internal/middleware"
	"github.com/twkakbar/family-tree/backend/pkg/response"
)

// UserHandler adapts HTTP to the UserService.
type UserHandler struct {
	users domain.UserService
}

// NewUserHandler builds the user endpoints handler.
func NewUserHandler(users domain.UserService) *UserHandler {
	return &UserHandler{users: users}
}

// GetMe handles GET /api/v1/users/me.
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}

	user, err := h.users.GetByID(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toUserResponse(user))
}

// UpdateMe handles PATCH /api/v1/users/me.
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}

	var req UpdateMeRequest
	if !bindJSON(c, &req) {
		return
	}

	user, err := h.users.UpdateProfile(c.Request.Context(), userID, domain.UpdateProfileInput{
		Name: req.Name,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toUserResponse(user))
}
