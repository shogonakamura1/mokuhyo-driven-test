package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/service"
	"github.com/mokuhyo-driven-test/api/pkg/auth"
)

type NodeHandler struct {
	nodeService    *service.NodeService
	projectService *service.ProjectService
}

func NewNodeHandler(nodeService *service.NodeService, projectService *service.ProjectService) *NodeHandler {
	return &NodeHandler{
		nodeService:    nodeService,
		projectService: projectService,
	}
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
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

	var req model.CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node, edge, err := h.nodeService.CreateNode(c.Request.Context(), projectID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": node, "edge": edge})
}

func (h *NodeHandler) UpdateNode(c *gin.Context) {
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

	nodeID, err := uuid.Parse(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node ID"})
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

	var req model.UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.nodeService.UpdateNode(c.Request.Context(), nodeID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *NodeHandler) DeleteNode(c *gin.Context) {
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

	nodeID, err := uuid.Parse(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node ID"})
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

	if err := h.nodeService.DeleteNode(c.Request.Context(), projectID, nodeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
