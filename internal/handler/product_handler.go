package handler

import (
	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/service"
	"shirt-store-api/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
    service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
    return &ProductHandler{service}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
    var req domain.CreateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    if err := h.service.CreateProduct(&req); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err.Error())
    }

    return response.Success(c, "Product created successfully", nil)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
    products, err := h.service.GetAllProducts()
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch products")
    }

    return response.Success(c, "Products retrieved successfully", products)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
    }

    product, err := h.service.GetProductByID(uint(id))
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, "Product not found")
    }

    return response.Success(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
    }

    var req domain.CreateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    if err := h.service.UpdateProduct(uint(id), &req); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err.Error())
    }

    return response.Success(c, "Product updated successfully", nil)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
    }

    if err := h.service.DeleteProduct(uint(id)); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err.Error())
    }

    return response.Success(c, "Product deleted successfully", nil)
}