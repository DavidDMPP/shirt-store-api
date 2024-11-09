// internal/handler/setting_handler.go
package handler

import (
	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/service"
	"shirt-store-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type SettingHandler struct {
    service *service.SettingService
}

func NewSettingHandler(service *service.SettingService) *SettingHandler {
    return &SettingHandler{service}
}

func (h *SettingHandler) GetMidtransConfig(c *fiber.Ctx) error {
    // Verify admin role (assuming middleware already checked)
    config, err := h.service.GetMidtransConfig()
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to get Midtrans configuration")
    }

    // Mask sensitive data
    maskedConfig := struct {
        Environment string `json:"environment"`
        IsActive    bool   `json:"is_active"`
    }{
        Environment: config.Environment,
        IsActive:    config.IsActive,
    }

    return response.Success(c, "Midtrans configuration retrieved successfully", maskedConfig)
}

func (h *SettingHandler) UpdateMidtransConfig(c *fiber.Ctx) error {
    var config domain.MidtransConfig
    if err := c.BodyParser(&config); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    // Validate
    if config.ServerKey == "" || config.ClientKey == "" {
        return response.Error(c, fiber.StatusBadRequest, "Server key and Client key are required")
    }

    if config.Environment != "sandbox" && config.Environment != "production" {
        return response.Error(c, fiber.StatusBadRequest, "Environment must be either 'sandbox' or 'production'")
    }

    if err := h.service.UpdateMidtransConfig(&config); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to update Midtrans configuration")
    }

    return response.Success(c, "Midtrans configuration updated successfully", nil)
}