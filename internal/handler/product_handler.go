package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/ahmadeko2017/backed-golang-tugas/internal/dto"
	"github.com/ahmadeko2017/backed-golang-tugas/internal/entity"
	"github.com/ahmadeko2017/backed-golang-tugas/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with name, price, stock, and category_id
// @Tags products
// @Accept json
// @Produce json
// @Param product body ProductCreateRequest true "Product Data"
// @Success 201 {object} entity.Product
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}

	if err := h.service.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get a list of all products, optionally filtered by name, with pagination
// @Tags products
// @Produce json
// @Param name query string false "Search by product name"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	name := c.Query("name")
	page, limit := dto.GetPaginationParams(c)

	products, total, err := h.service.GetAllProducts(name, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := dto.PaginatedResponse{
		Data: products,
		Meta: dto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Get details of a specific product
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 404 {object} map[string]string "Product not found"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product's details
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body entity.Product true "Product Data"
// @Success 200 {object} entity.Product
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateProduct(uint(id), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string "Message"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.service.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
