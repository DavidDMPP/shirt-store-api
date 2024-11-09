package handler

import (
	"shirt-store-api/internal/service"
	"shirt-store-api/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
    paymentService *service.PaymentService
    orderService   *service.OrderService
}

func NewPaymentHandler(paymentService *service.PaymentService, orderService *service.OrderService) *PaymentHandler {
    return &PaymentHandler{
        paymentService: paymentService,
        orderService:   orderService,
    }
}

func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
    // Convert string ID to uint
    idStr := c.Params("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid order ID")
    }

    order, err := h.orderService.GetOrderByID(uint(id))
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, "Order not found")
    }

    userID := c.Locals("userID").(uint)
    user, err := h.orderService.GetUserByID(userID)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to get user details")
    }

    paymentURL, err := h.paymentService.CreatePayment(order, user)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to create payment")
    }

    return response.Success(c, "Payment URL generated successfully", fiber.Map{
        "payment_url": paymentURL,
        "order_id":   order.ID,
    })
}

func (h *PaymentHandler) HandlePaymentNotification(c *fiber.Ctx) error {
    var notification map[string]interface{}
    if err := c.BodyParser(&notification); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid notification payload")
    }

    if err := h.paymentService.HandlePaymentNotification(notification); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to process payment notification")
    }

    return response.Success(c, "Payment notification processed successfully", nil)
}

func (h *PaymentHandler) CheckPaymentStatus(c *fiber.Ctx) error {
    orderID := c.Params("id")
    if orderID == "" {
        return response.Error(c, fiber.StatusBadRequest, "Order ID is required")
    }

    status, err := h.paymentService.GetPaymentStatus(orderID)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to check payment status")
    }

    return response.Success(c, "Payment status retrieved successfully", fiber.Map{
        "order_id": orderID,
        "status":   status,
    })
}