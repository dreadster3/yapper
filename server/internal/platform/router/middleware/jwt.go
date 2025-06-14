package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIdContextKey     = "user_id"
	UserClaimsContextKey = "user_claims"
)

type JWTMiddleware struct {
	jwksURL  string
	keyFunc  jwt.Keyfunc
	issuer   string
	audience string
}

type JWTConfig struct {
	JWKSUrl  string
	Issuer   string
	Audience string
}

func NewJWTMiddleware(config JWTConfig) (*JWTMiddleware, error) {
	jwks, err := keyfunc.NewDefault([]string{config.JWKSUrl})
	if err != nil {
		return nil, fmt.Errorf("failed to get JWK Set from URL: %w", err)
	}

	return &JWTMiddleware{
		jwksURL:  config.JWKSUrl,
		keyFunc:  jwks.Keyfunc,
		issuer:   config.Issuer,
		audience: config.Audience,
	}, nil
}

func (m *JWTMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}

		token, err := jwt.Parse(tokenString, m.keyFunc)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		if m.issuer != "" {
			if iss, ok := claims["iss"].(string); !ok || iss != m.issuer {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid issuer"})
				return
			}
		}

		if m.audience != "" {
			if aud, ok := claims["aud"].(string); !ok || aud != m.audience {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
				return
			}
		}

		c.Set(UserIdContextKey, claims["sub"])
		c.Set(UserClaimsContextKey, claims)
		c.Next()
	}
}
