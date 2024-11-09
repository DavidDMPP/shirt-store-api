package handler

import (
	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/service"
	"shirt-store-api/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
    orderService   *service.OrderService
    paymentService *service.PaymentService
}

func NewOrderHandler(orderService *service.OrderService, paymentService *service.PaymentService) *OrderHandler {
    return &OrderHandler{
        orderService:   orderService,
        paymentService: paymentService,
    }
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)

    var req domain.CreateOrderRequest
    if err := c.BodyParser(&req); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    order, err := h.orderService.CreateOrder(userID, &req)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err.Error())
    }

    return response.Success(c, "Order created successfully", order)
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)

    orders, err := h.orderService.GetUserOrders(userID)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch orders")
    }

    return response.Success(c, "Orders retrieved successfully", orders)
}

func (h *OrderHandler) GetOrderDetail(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid order ID")
    }

    order, err := h.orderService.GetOrderByID(uint(id))
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, "Order not found")
    }

    userID := c.Locals("userID").(uint)
    if order.UserID != userID && c.Locals("role") != "admin" {
        return response.Error(c, fiber.StatusForbidden, "Access denied")
    }

    return response.Success(c, "Order retrieved successfully", order)
}

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
    orders, err := h.orderService.GetAllOrders()
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch orders")
    }

    return response.Success(c, "Orders retrieved successfully", orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid order ID")
    }

    var req struct {
        Status string `json:"status"`
    }
    if err := c.BodyParser(&req); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    if err := h.orderService.UpdateOrderStatus(uint(id), req.Status); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err.Error())
    }

    return response.Success(c, "Order status updated successfully", nil)
}