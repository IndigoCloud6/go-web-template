package service

import (
	"context"

	"github.com/IndigoCloud6/go-web-template/internal/model"
	"github.com/IndigoCloud6/go-web-template/internal/repository"
	apperrors "github.com/IndigoCloud6/go-web-template/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Authenticate verifies user credentials and returns the user if valid
func (s *authService) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid email or password")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundErrorWithCause("user not found", err)
	}
	return user, nil
}
