package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/BhavaniNBL/ecommerce-backend/proto/productpb"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// func (h *ProductHandler) CreateProduct(c *gin.Context) {
// 	var product model.Product
// 	if err := c.ShouldBindJSON(&product); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := h.service.Create(&product); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, product)
// }

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set UUID if missing (safe fallback)
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	// Force timestamps
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Optional: Validation (ensure price > 0, category non-empty, etc.)
	if product.Price <= 0 || product.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price and Category must be provided"})
		return
	}

	if err := h.service.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	log.Println("Fetching product with ID:", id)
	product, err := h.service.GetByID(id)
	if err != nil {
		log.Println("Handler error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product.ToGRPC())
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.ID, _ = uuid.Parse(id)
	if err := h.service.Update(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	// c.Status(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// func (h *ProductHandler) ListProducts(c *gin.Context) {
// 	filters := map[string]string{
// 		"name":     c.Query("name"),
// 		"category": c.Query("category"),
// 	}
// 	products, err := h.service.List(filters)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
// 		return
// 	}

// 	// Convert to gRPC format
// 	var grpcProducts []*productpb.Product
// 	for _, p := range products {
// 		grpcProducts = append(grpcProducts, p.ToGRPC())
// 	}
// 	c.JSON(http.StatusOK, grpcProducts)
// }

func (h *ProductHandler) ListProducts(c *gin.Context) {
	resp, err := h.service.ListProducts(context.Background(), &productpb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var grpcProducts []*productpb.Product
	for _, p := range resp.Products {
		grpcProducts = append(grpcProducts, p)
	}

	c.JSON(http.StatusOK, grpcProducts)
}
