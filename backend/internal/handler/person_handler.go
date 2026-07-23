package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/internal/middleware"
	"github.com/twkakbar/family-tree/backend/pkg/pagination"
	"github.com/twkakbar/family-tree/backend/pkg/response"
	uploadpkg "github.com/twkakbar/family-tree/backend/pkg/upload"
)

// PersonHandler adapts HTTP to the PersonService.
type PersonHandler struct {
	persons domain.PersonService
}

// NewPersonHandler builds the person and tree endpoints handler.
func NewPersonHandler(persons domain.PersonService) *PersonHandler {
	return &PersonHandler{persons: persons}
}

// ListPersons handles GET /api/v1/families/:id/persons.
func (h *PersonHandler) ListPersons(c *gin.Context) {
	userID, familyID, ok := parseFamilyScope(c)
	if !ok {
		return
	}

	params := pagination.Parse(c.Query("page"), c.Query("per_page"))
	items, total, err := h.persons.ListByFamily(
		c.Request.Context(),
		userID,
		familyID,
		c.Query("search"),
		params.Page,
		params.PerPage,
	)
	if err != nil {
		respondError(c, err)
		return
	}

	out := make([]PersonResponse, len(items))
	for i := range items {
		out[i] = toPersonResponse(&items[i])
	}
	response.OKWithMeta(c, http.StatusOK, out, pagination.NewMeta(params, total))
}

// CreatePerson handles POST /api/v1/families/:id/persons.
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	userID, familyID, ok := parseFamilyScope(c)
	if !ok {
		return
	}

	var req PersonMutationRequest
	if !bindJSON(c, &req) {
		return
	}
	input, err := req.toCreateInput()
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, err.Error())
		return
	}

	person, err := h.persons.Create(c.Request.Context(), userID, familyID, input)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusCreated, toPersonResponse(person))
}

// GetFamilyTree handles GET /api/v1/families/:id/tree.
func (h *PersonHandler) GetFamilyTree(c *gin.Context) {
	userID, familyID, ok := parseFamilyScope(c)
	if !ok {
		return
	}

	tree, err := h.persons.GetFamilyTree(c.Request.Context(), userID, familyID)
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

// CreateRelationship handles POST /api/v1/families/:id/relationships.
func (h *PersonHandler) CreateRelationship(c *gin.Context) {
	userID, familyID, ok := parseFamilyScope(c)
	if !ok {
		return
	}

	var req CreateRelationshipRequest
	if !bindJSON(c, &req) {
		return
	}

	fromID, err := uuid.Parse(req.FromPersonID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid from_person_id")
		return
	}
	toID, err := uuid.Parse(req.ToPersonID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid to_person_id")
		return
	}

	rel, err := h.persons.CreateRelationship(c.Request.Context(), userID, familyID, domain.CreateRelationshipInput{
		FromPersonID: fromID,
		ToPersonID:   toID,
		Type:         domain.RelationshipType(req.Type),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusCreated, toRelationshipResponse(rel))
}

// GetPerson handles GET /api/v1/persons/:id.
func (h *PersonHandler) GetPerson(c *gin.Context) {
	userID, personID, ok := parsePersonScope(c)
	if !ok {
		return
	}

	detail, err := h.persons.GetByID(c.Request.Context(), userID, personID)
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

// UpdatePerson handles PATCH /api/v1/persons/:id.
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	userID, personID, ok := parsePersonScope(c)
	if !ok {
		return
	}

	var req PersonMutationRequest
	if !bindJSON(c, &req) {
		return
	}
	input, err := req.toUpdateInput()
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, err.Error())
		return
	}

	person, err := h.persons.Update(c.Request.Context(), userID, personID, input)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toPersonResponse(person))
}

// DeletePerson handles DELETE /api/v1/persons/:id.
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	userID, personID, ok := parsePersonScope(c)
	if !ok {
		return
	}

	if err := h.persons.Delete(c.Request.Context(), userID, personID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{"deleted": true})
}

// UploadPhoto handles PUT /api/v1/persons/:id/photo (multipart field "photo").
func (h *PersonHandler) UploadPhoto(c *gin.Context) {
	userID, personID, ok := parsePersonScope(c)
	if !ok {
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "photo file is required (multipart field \"photo\")")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "unable to read uploaded file")
		return
	}
	defer src.Close()

	data, err := uploadpkg.ReadLimited(src, 12<<20) // hard cap; service enforces configured max
	if err != nil {
		respondError(c, domain.ErrInvalidUpload)
		return
	}

	person, err := h.persons.UpdatePhoto(c.Request.Context(), userID, personID, data)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, toPersonResponse(person))
}

// GetAncestors handles GET /api/v1/persons/:id/ancestors.
func (h *PersonHandler) GetAncestors(c *gin.Context) {
	h.getTraversal(c, true)
}

// GetDescendants handles GET /api/v1/persons/:id/descendants.
func (h *PersonHandler) GetDescendants(c *gin.Context) {
	h.getTraversal(c, false)
}

// DeleteRelationship handles DELETE /api/v1/relationships/:id.
func (h *PersonHandler) DeleteRelationship(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}

	relID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid relationship id")
		return
	}

	if err := h.persons.DeleteRelationship(c.Request.Context(), userID, relID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *PersonHandler) getTraversal(c *gin.Context, ancestors bool) {
	userID, personID, ok := parsePersonScope(c)
	if !ok {
		return
	}

	depth := parseDepth(c.Query("depth"))
	var (
		items []domain.PersonWithDepth
		err   error
	)
	if ancestors {
		items, err = h.persons.GetAncestors(c.Request.Context(), userID, personID, depth)
	} else {
		items, err = h.persons.GetDescendants(c.Request.Context(), userID, personID, depth)
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

func parseFamilyScope(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
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

func parsePersonScope(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return uuid.Nil, uuid.Nil, false
	}
	personID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, codeBadRequest, "invalid person id")
		return uuid.Nil, uuid.Nil, false
	}
	return userID, personID, true
}

func parseDepth(raw string) int {
	if raw == "" {
		return 5
	}
	var depth int
	if _, err := fmt.Sscanf(raw, "%d", &depth); err != nil || depth < 1 {
		return 5
	}
	return depth
}

// ---- Person DTOs -----------------------------------------------------------

type PersonMutationRequest struct {
	FirstName string  `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string  `json:"last_name" binding:"max=100"`
	Gender    string  `json:"gender" binding:"required,oneof=male female other unknown"`
	BirthDate *string `json:"birth_date"`
	DeathDate *string `json:"death_date"`
	Bio       string  `json:"bio" binding:"max=5000"`
}

func (r PersonMutationRequest) toCreateInput() (domain.CreatePersonInput, error) {
	return r.toMutationInput()
}

func (r PersonMutationRequest) toUpdateInput() (domain.UpdatePersonInput, error) {
	in, err := r.toMutationInput()
	return domain.UpdatePersonInput(in), err
}

func (r PersonMutationRequest) toMutationInput() (domain.CreatePersonInput, error) {
	birth, err := parseOptionalDate(r.BirthDate)
	if err != nil {
		return domain.CreatePersonInput{}, err
	}
	death, err := parseOptionalDate(r.DeathDate)
	if err != nil {
		return domain.CreatePersonInput{}, err
	}
	return domain.CreatePersonInput{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Gender:    domain.Gender(r.Gender),
		BirthDate: birth,
		DeathDate: death,
		Bio:       r.Bio,
	}, nil
}

type CreateRelationshipRequest struct {
	FromPersonID string `json:"from_person_id" binding:"required"`
	ToPersonID   string `json:"to_person_id" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=parent spouse"`
}

type PersonResponse struct {
	ID        string    `json:"id"`
	FamilyID  string    `json:"family_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	BirthDate *string   `json:"birth_date"`
	DeathDate *string   `json:"death_date"`
	Bio       string    `json:"bio"`
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PersonDetailResponse struct {
	PersonResponse
	Relationships []RelationshipResponse `json:"relationships"`
}

type PersonDepthResponse struct {
	PersonResponse
	Depth int `json:"depth"`
}

type RelationshipResponse struct {
	ID           string    `json:"id"`
	FamilyID     string    `json:"family_id"`
	FromPersonID string    `json:"from_person_id"`
	ToPersonID   string    `json:"to_person_id"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
}

type FamilyTreeResponse struct {
	Persons       []PersonResponse       `json:"persons"`
	Relationships []RelationshipResponse `json:"relationships"`
}

func toPersonResponse(p *domain.Person) PersonResponse {
	photoURL := ""
	if p.PhotoPath != "" {
		photoURL = "/uploads/" + p.PhotoPath
	}
	return PersonResponse{
		ID:        p.ID.String(),
		FamilyID:  p.FamilyID.String(),
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Gender:    string(p.Gender),
		BirthDate: formatOptionalDate(p.BirthDate),
		DeathDate: formatOptionalDate(p.DeathDate),
		Bio:       p.Bio,
		PhotoURL:  photoURL,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toRelationshipResponse(r *domain.Relationship) RelationshipResponse {
	return RelationshipResponse{
		ID:           r.ID.String(),
		FamilyID:     r.FamilyID.String(),
		FromPersonID: r.FromPersonID.String(),
		ToPersonID:   r.ToPersonID.String(),
		Type:         string(r.Type),
		CreatedAt:    r.CreatedAt,
	}
}

func parseOptionalDate(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *raw)
	if err != nil {
		return nil, errInvalidDate
	}
	return &t, nil
}

func formatOptionalDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

var errInvalidDate = &dateParseError{}

type dateParseError struct{}

func (e *dateParseError) Error() string { return "dates must use YYYY-MM-DD format" }
