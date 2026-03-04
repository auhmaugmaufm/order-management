package middleware

import (
	"strings"

	"github.com/auhmaugmaufm/event-driven-order/internal/auth"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(401).JSON(dto.ErrorResponse{Error: "unauthorized", Message: "missing token"})
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(dto.ErrorResponse{Error: "unauthorized", Message: "invalid format"})
		}
		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(dto.ErrorResponse{Error: "unauthorized", Message: "invalid token"})
		}
		c.Locals("ID", claims.ID)
		c.Locals("email", claims.Email)
		return c.Next()
	}
}
