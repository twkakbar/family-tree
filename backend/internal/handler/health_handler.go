package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler exposes liveness and readiness probes.
type HealthHandler struct {
	db      *gorm.DB
	version string
}

// NewHealthHandler builds the health endpoints handler.
func NewHealthHandler(db *gorm.DB, version string) *HealthHandler {
	return &HealthHandler{db: db, version: version}
}

// Live handles GET /healthz: process is up (no dependencies checked).
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "version": h.version})
}

// Ready handles GET /readyz: process can serve traffic (DB reachable).
func (h *HealthHandler) Ready(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err == nil {
		err = sqlDB.PingContext(c.Request.Context())
	}
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
