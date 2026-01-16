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
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

type JWKSClient struct {
	url        string
	httpClient *http.Client
	cache      map[string]*rsa.PublicKey
	mu         sync.RWMutex
	lastFetch  time.Time
	cacheTTL   time.Duration
}

func NewJWKSClient(jwksURL string) *JWKSClient {
	return &JWKSClient{
		url:        jwksURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      make(map[string]*rsa.PublicKey),
		cacheTTL:   1 * time.Hour,
	}
}

func (c *JWKSClient) fetchJWKS() (*JWKS, error) {
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

func (c *JWKSClient) getPublicKey(kid string) (*rsa.PublicKey, error) {
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

func (c *JWKSClient) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found in token header")
		}

		return c.getPublicKey(kid)
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func ExtractUserID(token *jwt.Token) (uuid.UUID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("sub claim not found")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	return userID, nil
}

func AuthMiddleware(jwksClient *JWKSClient) gin.HandlerFunc {
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

		userID, err := ExtractUserID(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to extract user ID", "details": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("token", token)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	return userID.(uuid.UUID), true
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	return userID, ok
}
