package handler

import (
	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/service"
	"shirt-store-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
    return &UserHandler{service}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
    var req domain.RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    if err := h.service.Register(&req); err != nil {
        return response.Error(c, fiber.StatusBadRequest, err.Error())
    }

    return response.Success(c, "Registration successful", nil)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
    var req domain.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    token, err := h.service.Login(&req)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, err.Error())
    }

    return response.Success(c, "Login successful", fiber.Map{
        "token": token,
    })
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    
    user, err := h.service.GetUserByID(userID)
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, "User not found")
    }

    return response.Success(c, "Profile retrieved successfully", user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    
    user, err := h.service.GetUserByID(userID)
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, "User not found")
    }

    if err := c.BodyParser(user); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    if err := h.service.UpdateUser(user); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to update profile")
    }

    return response.Success(c, "Profile updated successfully", user)
}

func (h *UserHandler) MakeFirstAdmin(c *fiber.Ctx) error {
    var req struct {
        Email string `json:"email"`
    }

    if err := c.BodyParser(&req); err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
    }

    user, err := h.service.GetUserByEmail(req.Email)
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, "User not found")
    }

    user.Role = "admin"
    if err := h.service.UpdateUser(user); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to update user role")
    }

    return response.Success(c, "User role updated to admin successfully", nil)
}