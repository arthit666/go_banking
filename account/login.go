package account

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login account
// @Description Authenticate account and obtain access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param account body account.Login true "Login credentials"
// @Success 200 {object} account.TokenResponse
// @Router /login/ [post]
func (h *handler) Login(c *fiber.Ctx) error {
	req := &Login{}
	acc := &Account{}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Err{Massage: "payload invalid: " + err.Error()})
	}

	tx := h.DB.Where("email = ?", req.Email).First(acc)
	if tx.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(Err{Massage: tx.Error.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.ErrUnauthorized)
	}

	act, err := GenerateToken(acc.ID, time.Now().Add(time.Minute*15).Unix())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}

	rft, err := GenerateToken(acc.ID, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}

	return c.JSON(TokenResponse{
		AccessToken:  *act,
		RefreshToken: *rft,
	})
}
