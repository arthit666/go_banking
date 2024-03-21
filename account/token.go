package account

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// @Summary Refresh access token
// @Description Refresh the access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Header 200 {string} X-Refresh-Token "Refresh token bearer"
// @Success 200 {object} account.TokenResponse
// @Router /refresh [get]
func (h *handler) RefreshAccessToken(c *fiber.Ctx) error {
	refreshTokenStr := c.Get("X-Refresh-Token")
	if refreshTokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(Err{Massage: "missing refresh token"})
	}

	refreshTokenStr = strings.TrimPrefix(refreshTokenStr, "Bearer ")

	refreshToken, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !refreshToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(Err{Massage: "invalid refresh token"})
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(Err{Massage: "invalid refresh token"})
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(Err{Massage: "invalid refresh token"})
	}

	expTime := time.Unix(int64(exp), 0)
	if time.Now().After(expTime) {
		return c.Status(fiber.StatusUnauthorized).JSON(Err{Massage: "refresh token expired"})
	}

	act, err := GenerateToken(uint(claims["account_id"].(float64)), time.Now().Add(time.Minute*15).Unix())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}

	rft, err := GenerateToken(uint(claims["account_id"].(float64)), time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}

	return c.JSON(TokenResponse{
		AccessToken:  *act,
		RefreshToken: *rft,
	})
}

func GenerateToken(accId uint, exp int64) (*string, error) {
	tk := jwt.New(jwt.SigningMethodHS256)
	cl := tk.Claims.(jwt.MapClaims)
	cl["account_id"] = accId
	cl["exp"] = exp
	tkStr, err := tk.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	return &tkStr, nil
}
