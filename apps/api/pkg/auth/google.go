package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleJWKSClient はGoogleのJWKSを使用してJWTトークンを検証するクライアントです
type GoogleJWKSClient struct {
	url        string
	httpClient *http.Client
	cache      map[string]*rsa.PublicKey
	mu         sync.RWMutex
	lastFetch  time.Time
	cacheTTL   time.Duration
}

// NewGoogleJWKSClient は新しいGoogle JWKSクライアントを作成します
func NewGoogleJWKSClient() *GoogleJWKSClient {
	return &GoogleJWKSClient{
		url:        "https://www.googleapis.com/oauth2/v3/certs",
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      make(map[string]*rsa.PublicKey),
		cacheTTL:   1 * time.Hour,
	}
}

func (c *GoogleJWKSClient) fetchJWKS() (*JWKS, error) {
	resp, err := c.httpClient.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func (c *GoogleJWKSClient) getPublicKey(kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	if key, ok := c.cache[kid]; ok && time.Since(c.lastFetch) < c.cacheTTL {
		c.mu.RUnlock()
		return key, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double check
	if key, ok := c.cache[kid]; ok && time.Since(c.lastFetch) < c.cacheTTL {
		return key, nil
	}

	jwks, err := c.fetchJWKS()
	if err != nil {
		return nil, err
	}

	// Clear cache and rebuild
	c.cache = make(map[string]*rsa.PublicKey)
	for _, jwk := range jwks.Keys {
		if jwk.Kty != "RSA" {
			continue
		}

		nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			continue
		}

		eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
		if err != nil {
			continue
		}

		var eInt int
		for _, b := range eBytes {
			eInt = eInt<<8 | int(b)
		}

		publicKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: eInt,
		}

		c.cache[jwk.Kid] = publicKey
	}

	c.lastFetch = time.Now()

	if key, ok := c.cache[kid]; ok {
		return key, nil
	}

	return nil, errors.New("key not found")
}

func (c *GoogleJWKSClient) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// GoogleのIDトークンはRS256を使用
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found in token header")
		}

		return c.getPublicKey(kid)
	}, jwt.WithValidMethods([]string{"RS256"}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// ExtractGoogleUserID はGoogle IDトークンからユーザーIDを抽出します
// Google IDトークンのsubクレームはGoogleアカウントの一意なIDです
func ExtractGoogleUserID(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("sub claim not found")
	}

	return sub, nil
}

// ExtractGoogleUserInfo はGoogle IDトークンからユーザー情報を抽出します
func ExtractGoogleUserInfo(token *jwt.Token) (email, name, picture string, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", errors.New("invalid token claims")
	}

	email, _ = claims["email"].(string)
	name, _ = claims["name"].(string)
	picture, _ = claims["picture"].(string)

	return email, name, picture, nil
}

// GoogleAuthMiddleware はGoogle IDトークンを検証するミドルウェアです
// このミドルウェアは、Google User IDからデータベースのユーザーIDを取得するために
// AuthServiceを使用する必要があります
func GoogleAuthMiddleware(jwksClient *GoogleJWKSClient, getUserIDByGoogleID func(string) (uuid.UUID, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwksClient.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		googleUserID, err := ExtractGoogleUserID(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to extract user ID", "details": err.Error()})
			c.Abort()
			return
		}

		// Google User IDからデータベースのユーザーIDを取得
		userID, err := getUserIDByGoogleID(googleUserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get user ID", "details": err.Error()})
			c.Abort()
			return
		}

		// データベースのユーザーIDをコンテキストに保存
		c.Set("user_id", userID)
		c.Set("google_user_id", googleUserID)
		c.Set("token", token)
		c.Next()
	}
}

// GetGoogleUserID はコンテキストからGoogle User IDを取得します
func GetGoogleUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("google_user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GoogleOAuthConfig はGoogle OAuth設定を保持します
type GoogleOAuthConfig struct {
	config *oauth2.Config
}

// NewGoogleOAuthConfig は新しいGoogle OAuth設定を作成します
func NewGoogleOAuthConfig(clientID, clientSecret, redirectURL string) *GoogleOAuthConfig {
	return &GoogleOAuthConfig{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		},
	}
}

// ExchangeCodeForToken は認証コードをIDトークンに交換します
func (c *GoogleOAuthConfig) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.config.Exchange(ctx, code)
}
