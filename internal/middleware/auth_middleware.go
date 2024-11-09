package middleware

import (
	"shirt-store-api/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "success": false,
                "message": "Authorization header is required",
            })
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "success": false,
                "message": "Invalid authorization format",
            })
        }

        claims, err := jwt.ValidateToken(tokenParts[1])
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "success": false,
                "message": "Invalid or expired token",
            })
        }

        c.Locals("userID", claims.UserID)
        c.Locals("role", claims.Role)

        return c.Next()
    }
}

func AdminMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("role")
        if role != "admin" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "success": false,
                "message": "Admin access required",
            })
        }
        return c.Next()
    }
}