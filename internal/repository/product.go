package repository

import (
	"context"

	"github.com/IndigoCloud6/go-web-template/internal/model"
	"gorm.io/gorm"
)

// ProductRepository handles database operations for products
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id uint) (*model.Product, error)
	List(ctx context.Context, offset, limit int) ([]*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// List retrieves a list of products with pagination
func (r *productRepository) List(ctx context.Context, offset, limit int) ([]*model.Product, error) {
	var products []*model.Product
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Update updates a product
func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete deletes a product by ID
func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}

// Count returns the total number of products
func (r *productRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&count).Error
	return count, err
}
