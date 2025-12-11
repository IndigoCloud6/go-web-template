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
)

// ProductService handles business logic for products
type ProductService interface {
	Create(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error)
	GetByID(ctx context.Context, id uint) (*model.Product, error)
	List(ctx context.Context, page, pageSize int) ([]*model.Product, int64, error)
	Update(ctx context.Context, id uint, req *model.UpdateProductRequest) (*model.Product, error)
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	repo  repository.ProductRepository
	redis *redis.Client
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository, redis *redis.Client) ProductService {
	return &productService{
		repo:  repo,
		redis: redis,
	}
}

// Create creates a new product
func (s *productService) Create(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Clear product list cache
	keys, _ := s.redis.Keys(ctx, "products:list:*").Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}

	return product, nil
}

// GetByID retrieves a product by ID with caching
func (s *productService) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	// Try to get from cache
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var product model.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	// Get from database
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	productJSON, err := json.Marshal(product)
	if err != nil {
		logger.Warn("Failed to marshal product for caching", zap.Error(err))
	} else {
		s.redis.Set(ctx, cacheKey, productJSON, 5*time.Minute)
	}

	return product, nil
}

// List retrieves a list of products with pagination
func (s *productService) List(ctx context.Context, page, pageSize int) ([]*model.Product, int64, error) {
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

	products, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update updates a product
func (s *productService) Update(ctx context.Context, id uint, req *model.UpdateProductRequest) (*model.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	// Clear cache
	cacheKey := fmt.Sprintf("product:%d", id)
	s.redis.Del(ctx, cacheKey)

	// Clear product list cache
	keys, _ := s.redis.Keys(ctx, "products:list:*").Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}

	return product, nil
}

// Delete deletes a product
func (s *productService) Delete(ctx context.Context, id uint) error {
	// Check if product exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Clear cache
	cacheKey := fmt.Sprintf("product:%d", id)
	s.redis.Del(ctx, cacheKey)

	// Clear product list cache
	keys, _ := s.redis.Keys(ctx, "products:list:*").Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}

	return nil
}
