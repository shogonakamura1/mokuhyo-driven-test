package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/service"
	"github.com/mokuhyo-driven-test/api/pkg/auth"
)

type MeHandler struct {
	settingsService *service.SettingsService
}

func NewMeHandler(settingsService *service.SettingsService) *MeHandler {
	return &MeHandler{settingsService: settingsService}
}

func (h *MeHandler) GetMe(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	settings, err := h.settingsService.GetSettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.MeResponse{
		User: model.UserInfo{
			ID: userID,
		},
		Settings: *settings,
	})
}
