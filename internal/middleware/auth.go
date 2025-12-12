package middleware

import (
	"strings"
	"time"

	"github.com/IndigoCloud6/go-web-template/internal/config"
	apperrors "github.com/IndigoCloud6/go-web-template/pkg/errors"
	"github.com/IndigoCloud6/go-web-template/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTAuth creates a JWT authentication middleware
func JWTAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorFromAppError(c, apperrors.NewUnauthorizedError("authorization header is required"))
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.ErrorFromAppError(c, apperrors.NewUnauthorizedError("invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.NewUnauthorizedError("invalid signing method")
			}
			return []byte(cfg.Secret), nil
		})

		if err != nil {
			response.ErrorFromAppError(c, apperrors.NewUnauthorizedError("invalid or expired token"))
			c.Abort()
			return
		}

		if !token.Valid {
			response.ErrorFromAppError(c, apperrors.NewUnauthorizedError("invalid token"))
			c.Abort()
			return
		}

		// Store user information in context for later use
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(cfg *config.JWTConfig, userID uint, email string) (string, error) {
	expirationHours := cfg.ExpirationHours
	if expirationHours <= 0 {
		expirationHours = 24 // default to 24 hours
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// GetUserIDFromContext retrieves the user ID from the gin context
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetEmailFromContext retrieves the email from the gin context
func GetEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

// GetClaimsFromContext retrieves the claims from the gin context
func GetClaimsFromContext(c *gin.Context) (*Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	c1, ok := claims.(*Claims)
	return c1, ok
}
