package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/pkg/response"
)

// TransferHandler serves JSON tree export/import.
type TransferHandler struct {
	transfer domain.TransferService
}

// NewTransferHandler builds export/import endpoints.
func NewTransferHandler(transfer domain.TransferService) *TransferHandler {
	return &TransferHandler{transfer: transfer}
}

// ExportFamily handles GET /api/v1/families/:id/export?format=json.
func (h *TransferHandler) ExportFamily(c *gin.Context) {
	userID, familyID, ok := parseFamilyScope(c)
	if !ok {
		return
	}
	format := c.DefaultQuery("format", "json")
	if format != "json" {
		respondError(c, domain.ErrUnsupportedFormat)
		return
	}

	doc, err := h.transfer.ExportJSON(c.Request.Context(), userID, familyID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.Header("Content-Disposition", `attachment; filename="kesultanan-aceh-export.json"`)
	response.OK(c, http.StatusOK, doc)
}

// ImportFamily handles POST /api/v1/families/:id/import.
func (h *TransferHandler) ImportFamily(c *gin.Context) {
	userID, familyID, ok := parseFamilyScope(c)
	if !ok {
		return
	}

	var doc domain.TreeTransfer
	if !bindJSON(c, &doc) {
		return
	}

	result, err := h.transfer.ImportJSON(c.Request.Context(), userID, familyID, &doc)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, result)
}
