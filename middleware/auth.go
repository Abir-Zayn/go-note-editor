package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

var jwksCache jwk.Set
var jwksCacheTime time.Time

func getJWKS(supabaseURL string) (jwk.Set, error) {
	// Cache JWKS for 1 hour
	if jwksCache != nil && time.Since(jwksCacheTime) < time.Hour {
		return jwksCache, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jwksURL := supabaseURL + "/auth/v1/.well-known/jwks.json"
	set, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return nil, err
	}

	jwksCache = set
	jwksCacheTime = time.Now()
	return set, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		supabaseURL := os.Getenv("SUPABASE_URL")

		// Fetch JWKS
		keySet, err := getJWKS(supabaseURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch JWKS", "details": err.Error()})
			c.Abort()
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, jwt.ErrTokenMalformed
			}

			key, found := keySet.LookupKeyID(kid)
			if !found {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			var rawKey interface{}
			if err := key.Raw(&rawKey); err != nil {
				return nil, err
			}

			return rawKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
