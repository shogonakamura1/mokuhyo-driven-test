package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mokuhyo-driven-test/api/internal/service"
	"github.com/mokuhyo-driven-test/api/pkg/auth"
)

type AuthHandler struct {
	authService *service.AuthService
	oauthConfig *auth.GoogleOAuthConfig
}

func NewAuthHandler(authService *service.AuthService, oauthConfig *auth.GoogleOAuthConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		oauthConfig: oauthConfig,
	}
}

type GoogleAuthRequest struct {
	Code string `json:"code" binding:"required"`
}

type GoogleAuthResponse struct {
	IDToken string `json:"id_token"`
	User    struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture,omitempty"`
	} `json:"user"`
}

// HandleGoogleAuth はGoogle認証コードをIDトークンに交換します
func (h *AuthHandler) HandleGoogleAuth(c *gin.Context) {
	var req GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 認証コードをIDトークンに交換
	token, err := h.oauthConfig.ExchangeCodeForToken(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to exchange code for token", "details": err.Error()})
		return
	}

	// IDトークンを取得
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "id_token not found in token response"})
		return
	}

	// IDトークンを検証してユーザー情報を取得
	jwksClient := auth.NewGoogleJWKSClient()
	verifiedToken, err := jwksClient.VerifyToken(idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to verify token", "details": err.Error()})
		return
	}

	email, name, picture, err := auth.ExtractGoogleUserInfo(verifiedToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to extract user info", "details": err.Error()})
		return
	}

	googleUserID, err := auth.ExtractGoogleUserID(verifiedToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to extract user ID", "details": err.Error()})
		return
	}

	// pictureを*stringに変換（空文字列の場合はnil）
	var picturePtr *string
	if picture != "" {
		picturePtr = &picture
	}

	// ユーザーを作成または取得
	user, err := h.authService.GetOrCreateUser(c.Request.Context(), googleUserID, email, name, picturePtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get or create user", "details": err.Error()})
		return
	}

	// user.Pictureをstringに変換（nilの場合は空文字列）
	pictureStr := ""
	if user.Picture != nil {
		pictureStr = *user.Picture
	}

	c.JSON(http.StatusOK, GoogleAuthResponse{
		IDToken: idToken,
		User: struct {
			ID      string `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture,omitempty"`
		}{
			ID:      user.ID.String(),
			Email:   user.Email,
			Name:    user.Name,
			Picture: pictureStr,
		},
	})
}
