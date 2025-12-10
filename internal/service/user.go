package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IndigoCloud6/go-web-template/internal/model"
	"github.com/IndigoCloud6/go-web-template/internal/repository"
	"github.com/IndigoCloud6/go-web-template/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for users
type UserService interface {
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	repo  repository.UserRepository
	redis *redis.Client
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, redis *redis.Client) UserService {
	return &userService{
		repo:  repo,
		redis: redis,
	}
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Age:      req.Age,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Clear user list cache
	s.redis.Del(ctx, "users:list:*")

	return user, nil
}

// GetByID retrieves a user by ID with caching
func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// Try to get from cache
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user model.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	// Get from database
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	userJSON, err := json.Marshal(user)
	if err != nil {
		logger.Warn("Failed to marshal user for caching", zap.Error(err))
	} else {
		s.redis.Set(ctx, cacheKey, userJSON, 5*time.Minute)
	}

	return user, nil
}

// List retrieves a list of users with pagination
func (s *userService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	users, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update updates a user
func (s *userService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email already exists
		existingUser, err := s.repo.GetByEmail(ctx, req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("email already exists")
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = string(hashedPassword)
	}
	if req.Age > 0 {
		user.Age = req.Age
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Clear cache
	cacheKey := fmt.Sprintf("user:%d", id)
	s.redis.Del(ctx, cacheKey)
	s.redis.Del(ctx, "users:list:*")

	return user, nil
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id uint) error {
	// Check if user exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Clear cache
	cacheKey := fmt.Sprintf("user:%d", id)
	s.redis.Del(ctx, cacheKey)
	s.redis.Del(ctx, "users:list:*")

	return nil
}
