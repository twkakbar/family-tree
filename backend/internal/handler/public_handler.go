package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/internal/middleware"
	"github.com/twkakbar/family-tree/backend/pkg/response"
)

// PublicHandler serves unauthenticated royal-tree reads.
type PublicHandler struct {
	public domain.PublicTreeService
}

// NewPublicHandler builds public tree endpoints.
func NewPublicHandler(public domain.PublicTreeService) *PublicHandler {
	return &PublicHandler{public: public}
}

// GetMeta handles GET /api/v1/public/meta.
func (h *PublicHandler) GetMeta(c *gin.Context) {
	family, err := h.public.Meta(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{
		"id":          family.ID.String(),
		"name":        family.Name,
		"description": family.Description,
	})
}

// GetTree handles GET /api/v1/public/tree.
func (h *PublicHandler) GetTree(c *gin.Context) {
	tree, err := h.public.Tree(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	persons := make([]PersonResponse, len(tree.Persons))
	for i := range tree.Persons {
		persons[i] = toPersonResponse(&tree.Persons[i])
	}
	rels := make([]RelationshipResponse, len(tree.Relationships))
	for i := range tree.Relationships {
		rels[i] = toRelationshipResponse(&tree.Relationships[i])
	}
	response.OK(c, http.StatusOK, FamilyTreeResponse{Persons: persons, Relationships: rels})
}

// GetPerson handles GET /api/v1/public/persons/:id.
func (h *PublicHandler) GetPerson(c *gin.Context) {
	personID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid person id")
		return
	}
	detail, err := h.public.GetPerson(c.Request.Context(), personID)
	if err != nil {
		respondError(c, err)
		return
	}
	rels := make([]RelationshipResponse, len(detail.Relationships))
	for i := range detail.Relationships {
		rels[i] = toRelationshipResponse(&detail.Relationships[i])
	}
	response.OK(c, http.StatusOK, PersonDetailResponse{
		PersonResponse: toPersonResponse(&detail.Person),
		Relationships:  rels,
	})
}

// GetAncestors handles GET /api/v1/public/persons/:id/ancestors.
func (h *PublicHandler) GetAncestors(c *gin.Context) {
	h.traversal(c, true)
}

// GetDescendants handles GET /api/v1/public/persons/:id/descendants.
func (h *PublicHandler) GetDescendants(c *gin.Context) {
	h.traversal(c, false)
}

func (h *PublicHandler) traversal(c *gin.Context, ancestors bool) {
	personID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid person id")
		return
	}
	depth := parseDepth(c.Query("depth"))
	var items []domain.PersonWithDepth
	if ancestors {
		items, err = h.public.Ancestors(c.Request.Context(), personID, depth)
	} else {
		items, err = h.public.Descendants(c.Request.Context(), personID, depth)
	}
	if err != nil {
		respondError(c, err)
		return
	}
	out := make([]PersonDepthResponse, len(items))
	for i := range items {
		out[i] = PersonDepthResponse{
			PersonResponse: toPersonResponse(&items[i].Person),
			Depth:          items[i].Depth,
		}
	}
	response.OK(c, http.StatusOK, out)
}

// AccessRequestHandler serves edit-access request endpoints.
type AccessRequestHandler struct {
	access domain.AccessRequestService
	royal  domain.RoyalFamilyService
}

// NewAccessRequestHandler builds access-request endpoints.
func NewAccessRequestHandler(access domain.AccessRequestService, royal domain.RoyalFamilyService) *AccessRequestHandler {
	return &AccessRequestHandler{access: access, royal: royal}
}

type submitAccessRequestBody struct {
	Message string `json:"message" binding:"required,min=10,max=2000"`
}

// Submit handles POST /api/v1/access-requests.
func (h *AccessRequestHandler) Submit(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	var req submitAccessRequestBody
	if !bindJSON(c, &req) {
		return
	}
	created, err := h.access.Submit(c.Request.Context(), userID, req.Message)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusCreated, toAccessRequestResponse(created))
}

// MyLatest handles GET /api/v1/access-requests/me.
func (h *AccessRequestHandler) MyLatest(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	req, err := h.access.MyLatest(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toAccessRequestResponse(req))
}

// CanEdit handles GET /api/v1/studio/can-edit.
func (h *AccessRequestHandler) CanEdit(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	can, err := h.royal.CanEdit(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	familyID, err := h.royal.ID(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{
		"can_edit":  can,
		"family_id": familyID.String(),
	})
}

// ListPending handles GET /api/v1/admin/access-requests.
func (h *AccessRequestHandler) ListPending(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	items, err := h.access.ListPending(c.Request.Context(), adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	out := make([]AccessRequestAdminResponse, len(items))
	for i := range items {
		out[i] = toAccessRequestAdminResponse(&items[i])
	}
	response.OK(c, http.StatusOK, out)
}

// Approve handles POST /api/v1/admin/access-requests/:id/approve.
func (h *AccessRequestHandler) Approve(c *gin.Context) {
	h.review(c, true)
}

// Reject handles POST /api/v1/admin/access-requests/:id/reject.
func (h *AccessRequestHandler) Reject(c *gin.Context) {
	h.review(c, false)
}

func (h *AccessRequestHandler) review(c *gin.Context, approve bool) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid request id")
		return
	}
	var req *domain.EditAccessRequest
	if approve {
		req, err = h.access.Approve(c.Request.Context(), adminID, requestID)
	} else {
		req, err = h.access.Reject(c.Request.Context(), adminID, requestID)
	}
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toAccessRequestResponse(req))
}

type AccessRequestResponse struct {
	ID         string     `json:"id"`
	Status     string     `json:"status"`
	Message    string     `json:"message"`
	CreatedAt  time.Time  `json:"created_at"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
}

type AccessRequestAdminResponse struct {
	AccessRequestResponse
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func toAccessRequestResponse(req *domain.EditAccessRequest) AccessRequestResponse {
	return AccessRequestResponse{
		ID:         req.ID.String(),
		Status:     string(req.Status),
		Message:    req.Message,
		CreatedAt:  req.CreatedAt,
		ReviewedAt: req.ReviewedAt,
	}
}

func toAccessRequestAdminResponse(summary *domain.AccessRequestSummary) AccessRequestAdminResponse {
	return AccessRequestAdminResponse{
		AccessRequestResponse: toAccessRequestResponse(&summary.Request),
		UserID:                summary.User.ID.String(),
		Name:                  summary.User.Name,
		Email:                 summary.User.Email,
	}
}
