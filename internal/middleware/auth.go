package middleware

import (
    "strings"
    "github.com/gofiber/fiber/v2"
    "github.com/mohammed-ayoub-js/endline-backend/internal/config"
    "github.com/mohammed-ayoub-js/endline-backend/internal/utils"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
        }

        claims, err := utils.ValidateAccessToken(parts[1], cfg)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
        }

        c.Locals("userID", claims.UserID)
        c.Locals("username", claims.Username)
        return c.Next()
    }
}