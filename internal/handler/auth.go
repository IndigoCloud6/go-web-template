package handler

import (
	"github.com/IndigoCloud6/go-web-template/internal/config"
	"github.com/IndigoCloud6/go-web-template/internal/middleware"
	"github.com/IndigoCloud6/go-web-template/internal/service"
	apperrors "github.com/IndigoCloud6/go-web-template/pkg/errors"
	"github.com/IndigoCloud6/go-web-template/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService service.AuthService
	jwtConfig   *config.JWTConfig
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, jwtConfig *config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtConfig:   jwtConfig,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorFromAppError(c, apperrors.NewValidationErrorWithCause("invalid request body", err))
		return
	}

	// Authenticate user
	user, err := h.authService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(h.jwtConfig, user.ID, user.Email)
	if err != nil {
		response.ErrorFromAppError(c, apperrors.NewInternalErrorWithCause("failed to generate token", err))
		return
	}

	resp := LoginResponse{
		Token: token,
	}
	resp.User.ID = user.ID
	resp.User.Name = user.Name
	resp.User.Email = user.Email

	response.Success(c, resp)
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh an existing valid JWT token
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.ErrorFromAppError(c, apperrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	email, exists := middleware.GetEmailFromContext(c)
	if !exists {
		response.ErrorFromAppError(c, apperrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	// Generate new token
	token, err := middleware.GenerateToken(h.jwtConfig, userID, email)
	if err != nil {
		response.ErrorFromAppError(c, apperrors.NewInternalErrorWithCause("failed to generate token", err))
		return
	}

	response.Success(c, map[string]string{
		"token": token,
	})
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.ErrorFromAppError(c, apperrors.NewUnauthorizedError("user not authenticated"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}

	response.Success(c, map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
	})
}
