// pkg/response/response.go
package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
    return c.JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func Error(c *fiber.Ctx, statusCode int, message string) error {
    return c.Status(statusCode).JSON(Response{
        Success: false,
        Message: message,
    })
}

func ValidationError(c *fiber.Ctx, err error) error {
    return c.Status(fiber.StatusBadRequest).JSON(Response{
        Success: false,
        Message: "Validation failed",
        Error:   err.Error(),
    })
}