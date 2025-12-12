package handler

import (
	"strconv"

	"github.com/IndigoCloud6/go-web-template/internal/model"
	"github.com/IndigoCloud6/go-web-template/internal/service"
	apperrors "github.com/IndigoCloud6/go-web-template/pkg/errors"
	"github.com/IndigoCloud6/go-web-template/pkg/response"
	"github.com/gin-gonic/gin"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param product body model.CreateProductRequest true "Product information"
// @Success 200 {object} response.Response{data=model.Product}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorFromAppError(c, apperrors.NewValidationErrorWithCause("invalid request body", err))
		return
	}

	product, err := h.productService.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}

	response.Success(c, product)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a product's information by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response{data=model.Product}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorFromAppError(c, apperrors.NewValidationError("invalid product id"))
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}

	response.Success(c, product)
}

// ListProducts godoc
// @Summary List products
// @Description Get a paginated list of products
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 500 {object} response.Response
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	products, total, err := h.productService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}

	response.Success(c, map[string]interface{}{
		"products":  products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product's information
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.UpdateProductRequest true "Product information"
// @Success 200 {object} response.Response{data=model.Product}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorFromAppError(c, apperrors.NewValidationError("invalid product id"))
		return
	}

	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorFromAppError(c, apperrors.NewValidationErrorWithCause("invalid request body", err))
		return
	}

	product, err := h.productService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}

	response.Success(c, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorFromAppError(c, apperrors.NewValidationError("invalid product id"))
		return
	}

	if err := h.productService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.ErrorFromAppError(c, err)
		return
	}

	response.SuccessWithMessage(c, "product deleted successfully", nil)
}
