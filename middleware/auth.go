package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Err struct {
	Massage string `json:"massage"`
}

func ExtractUserFromJWT(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(Err{Massage: "missing jwt token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(Err{Massage: "invalid jwt token"})
	}

	accountIDFloat, ok := claims["account_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(Err{Massage: "invalid jwt token"})
	}

	accountID := int(accountIDFloat)
	c.Locals("account_id", accountID)
	return c.Next()
}
