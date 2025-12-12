package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IndigoCloud6/go-web-template/internal/config"
	"github.com/gin-gonic/gin"
)

func TestGenerateToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:          "test-secret-key",
		ExpirationHours: 24,
		Issuer:          "test-issuer",
	}

	token, err := GenerateToken(cfg, 123, "test@example.com")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.JWTConfig{
		Secret:          "test-secret-key",
		ExpirationHours: 24,
		Issuer:          "test-issuer",
	}

	// Generate a valid token
	token, err := GenerateToken(cfg, 123, "test@example.com")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Create a test router with the middleware
	r := gin.New()
	r.Use(JWTAuth(cfg))
	r.GET("/test", func(c *gin.Context) {
		userID, _ := GetUserIDFromContext(c)
		email, _ := GetEmailFromContext(c)
		c.JSON(http.StatusOK, gin.H{"user_id": userID, "email": email})
	})

	// Create a request with the token
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestJWTAuth_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.JWTConfig{
		Secret:          "test-secret-key",
		ExpirationHours: 24,
		Issuer:          "test-issuer",
	}

	r := gin.New()
	r.Use(JWTAuth(cfg))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.JWTConfig{
		Secret:          "test-secret-key",
		ExpirationHours: 24,
		Issuer:          "test-issuer",
	}

	r := gin.New()
	r.Use(JWTAuth(cfg))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWTAuth_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.JWTConfig{
		Secret:          "test-secret-key",
		ExpirationHours: 24,
		Issuer:          "test-issuer",
	}

	r := gin.New()
	r.Use(JWTAuth(cfg))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetUserIDFromContext_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID, exists := GetUserIDFromContext(c)
	if exists {
		t.Error("Expected exists to be false when user_id not set")
	}
	if userID != 0 {
		t.Errorf("Expected userID to be 0, got %d", userID)
	}
}

func TestGetEmailFromContext_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	email, exists := GetEmailFromContext(c)
	if exists {
		t.Error("Expected exists to be false when email not set")
	}
	if email != "" {
		t.Errorf("Expected email to be empty, got %s", email)
	}
}

func TestGetClaimsFromContext_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	claims, exists := GetClaimsFromContext(c)
	if exists {
		t.Error("Expected exists to be false when claims not set")
	}
	if claims != nil {
		t.Error("Expected claims to be nil")
	}
}

func TestGenerateToken_DefaultExpiration(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:          "test-secret-key",
		ExpirationHours: 0, // Should default to 24
		Issuer:          "test-issuer",
	}

	token, err := GenerateToken(cfg, 123, "test@example.com")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}
}
