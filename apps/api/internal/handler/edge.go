package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/service"
	"github.com/mokuhyo-driven-test/api/pkg/auth"
)

type EdgeHandler struct {
	edgeService    *service.EdgeService
	projectService *service.ProjectService
}

func NewEdgeHandler(edgeService *service.EdgeService, projectService *service.ProjectService) *EdgeHandler {
	return &EdgeHandler{
		edgeService:    edgeService,
		projectService: projectService,
	}
}

func (h *EdgeHandler) UpdateEdge(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	edgeID, err := uuid.Parse(c.Param("edgeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid edge ID"})
		return
	}

	// Check ownership
	owned, err := h.projectService.CheckOwnership(c.Request.Context(), projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !owned {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req model.UpdateEdgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.edgeService.UpdateEdge(c.Request.Context(), edgeID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *EdgeHandler) Reorder(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Check ownership
	owned, err := h.projectService.CheckOwnership(c.Request.Context(), projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !owned {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req model.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.edgeService.Reorder(c.Request.Context(), projectID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
