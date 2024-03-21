package pocket

import (
	"fmt"

	"github.com/arthit666/make_app/balance"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// @Summary Create a new pocket
// @Description Create a new pocket with the provided data
// @Tags pockets
// @Accept json
// @Produce json
// @Param pocket body pocket.PocketCreate true "PocketCreate data"
// @Success 201 {object} pocket.SuccessResponse
// @Security  Bearer
// @Router /pockets/ [post]
func (h *handler) CreatePocket(c *fiber.Ctx) error {
	pc := &PocketCreate{}

	if err := c.BodyParser(pc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(pc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Err{Message: "payload invalid: " + err.Error()})
	}

	acc := c.Locals("account_id").(int)

	p := &Pocket{
		Title:       pc.Title,
		AccountID:   uint(acc),
		Balance:     pc.Balance,
		Description: pc.Description,
	}

	err := balance.Deduct(h.DB, p.AccountID, p.Balance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + err.Error()})
	}

	_, err = create(h, p)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{Message: "create pocket success"})
}

func create(h *handler, p *Pocket) (*Pocket, error) {
	tx := h.DB.Create(p)
	if tx.Error != nil {
		return nil, fmt.Errorf(tx.Error.Error())
	}
	return p, nil
}
