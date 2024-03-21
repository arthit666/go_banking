package pocket

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// @Summary Update a pocket
// @Description Update a pocket with the provided data
// @Tags pockets
// @Accept json
// @Produce json
// @Param id path int true "Pocket ID"
// @Param pocket body pocket.PocketUpdate true "PocketUpdate data"
// @Success 201 {object} pocket.SuccessResponse
// @Security  Bearer
// @Router /pockets/{id} [put]
func (h *handler) UpdatePocket(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	pr := &PocketUpdate{}
	if err := c.BodyParser(pr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(pr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Err{Message: "payload invalid: " + err.Error()})
	}

	acc := c.Locals("account_id").(int)
	accStr := strconv.Itoa(acc)

	p, err := getById(c.Params("id"), accStr, h)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Err{Message: "pocket not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + err.Error()})
	}

	p.ID = uint(id)
	p.Title = pr.Title
	p.Description = pr.Description

	tx := h.DB.Model(p).Updates(*p)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "update error: " + tx.Error.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(SuccessResponse{Message: "update pocket success"})

}
