package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/internal/middleware"
	"github.com/twkakbar/family-tree/backend/pkg/pagination"
	"github.com/twkakbar/family-tree/backend/pkg/response"
)

// FamilyHandler adapts HTTP to the FamilyService.
type FamilyHandler struct {
	families domain.FamilyService
}

// NewFamilyHandler builds the family endpoints handler.
func NewFamilyHandler(families domain.FamilyService) *FamilyHandler {
	return &FamilyHandler{families: families}
}

// CreateFamily handles POST /api/v1/families.
func (h *FamilyHandler) CreateFamily(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}

	var req CreateFamilyRequest
	if !bindJSON(c, &req) {
		return
	}

	summary, err := h.families.Create(c.Request.Context(), userID, domain.CreateFamilyInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusCreated, toFamilyResponse(summary))
}

// ListFamilies handles GET /api/v1/families.
func (h *FamilyHandler) ListFamilies(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}

	params := pagination.Parse(c.Query("page"), c.Query("per_page"))
	summaries, total, err := h.families.ListForUser(c.Request.Context(), userID, params.Page, params.PerPage)
	if err != nil {
		respondError(c, err)
		return
	}

	items := make([]FamilyResponse, len(summaries))
	for i := range summaries {
		items[i] = toFamilyResponse(&summaries[i])
	}
	response.OKWithMeta(c, http.StatusOK, items, pagination.NewMeta(params, total))
}

// GetFamily handles GET /api/v1/families/:id.
func (h *FamilyHandler) GetFamily(c *gin.Context) {
	userID, familyID, ok := h.parseFamilyScope(c)
	if !ok {
		return
	}

	summary, err := h.families.GetByID(c.Request.Context(), userID, familyID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toFamilyResponse(summary))
}

// UpdateFamily handles PATCH /api/v1/families/:id.
func (h *FamilyHandler) UpdateFamily(c *gin.Context) {
	userID, familyID, ok := h.parseFamilyScope(c)
	if !ok {
		return
	}

	var req UpdateFamilyRequest
	if !bindJSON(c, &req) {
		return
	}

	summary, err := h.families.Update(c.Request.Context(), userID, familyID, domain.UpdateFamilyInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toFamilyResponse(summary))
}

// DeleteFamily handles DELETE /api/v1/families/:id.
func (h *FamilyHandler) DeleteFamily(c *gin.Context) {
	userID, familyID, ok := h.parseFamilyScope(c)
	if !ok {
		return
	}

	if err := h.families.Delete(c.Request.Context(), userID, familyID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{"deleted": true})
}

// ListCollaborators handles GET /api/v1/families/:id/collaborators.
func (h *FamilyHandler) ListCollaborators(c *gin.Context) {
	userID, familyID, ok := h.parseFamilyScope(c)
	if !ok {
		return
	}

	collabs, err := h.families.ListCollaborators(c.Request.Context(), userID, familyID)
	if err != nil {
		respondError(c, err)
		return
	}

	items := make([]CollaboratorResponse, len(collabs))
	for i := range collabs {
		items[i] = toCollaboratorResponse(&collabs[i])
	}
	response.OK(c, http.StatusOK, items)
}

// AddCollaborator handles POST /api/v1/families/:id/collaborators.
func (h *FamilyHandler) AddCollaborator(c *gin.Context) {
	userID, familyID, ok := h.parseFamilyScope(c)
	if !ok {
		return
	}

	var req AddCollaboratorRequest
	if !bindJSON(c, &req) {
		return
	}

	collab, err := h.families.AddCollaborator(c.Request.Context(), userID, familyID, domain.AddCollaboratorInput{
		Email: req.Email,
		Role:  domain.FamilyRole(req.Role),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusCreated, toCollaboratorResponse(collab))
}

// UpdateCollaborator handles PATCH /api/v1/families/:id/collaborators/:userId.
func (h *FamilyHandler) UpdateCollaborator(c *gin.Context) {
	userID, familyID, ok := h.parseFamilyScope(c)
	if !ok {
		return
	}

	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid user id")
		return
	}

	var req UpdateCollaboratorRequest
	if !bindJSON(c, &req) {
		return
	}

	collab, err := h.families.UpdateCollaboratorRole(c.Request.Context(), userID, familyID, targetUserID, domain.FamilyRole(req.Role))
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toCollaboratorResponse(collab))
}

// RemoveCollaborator handles DELETE /api/v1/families/:id/collaborators/:userId.
func (h *FamilyHandler) RemoveCollaborator(c *gin.Context) {
	userID, familyID, ok := h.parseFamilyScope(c)
	if !ok {
		return
	}

	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid user id")
		return
	}

	if err := h.families.RemoveCollaborator(c.Request.Context(), userID, familyID, targetUserID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{"removed": true})
}

func (h *FamilyHandler) parseFamilyScope(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return uuid.Nil, uuid.Nil, false
	}
	familyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid family id")
		return uuid.Nil, uuid.Nil, false
	}
	return userID, familyID, true
}

// ---- Family DTOs -----------------------------------------------------------

type CreateFamilyRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateFamilyRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=1000"`
}

type AddCollaboratorRequest struct {
	Email string `json:"email" binding:"required,email,max=254"`
	Role  string `json:"role" binding:"required,oneof=editor viewer"`
}

type UpdateCollaboratorRequest struct {
	Role string `json:"role" binding:"required,oneof=owner editor viewer"`
}

type FamilyResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Role         string    `json:"role"`
	PersonsCount int64     `json:"persons_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CollaboratorResponse struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func toFamilyResponse(summary *domain.FamilySummary) FamilyResponse {
	return FamilyResponse{
		ID:           summary.Family.ID.String(),
		Name:         summary.Family.Name,
		Description:  summary.Family.Description,
		Role:         string(summary.Role),
		PersonsCount: summary.PersonsCount,
		CreatedAt:    summary.Family.CreatedAt,
		UpdatedAt:    summary.Family.UpdatedAt,
	}
}

func toCollaboratorResponse(summary *domain.CollaboratorSummary) CollaboratorResponse {
	return CollaboratorResponse{
		UserID:    summary.User.ID.String(),
		Name:      summary.User.Name,
		Email:     summary.User.Email,
		Role:      string(summary.Collaborator.Role),
		CreatedAt: summary.Collaborator.CreatedAt,
	}
}
