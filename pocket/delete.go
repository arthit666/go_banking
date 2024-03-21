package pocket

import (
	"errors"
	"strconv"

	"github.com/arthit666/make_app/balance"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// @Summary Delete a pocket
// @Description Delete a pocket by ID
// @Param id path int true "Pocket ID"
// @Tags pockets
// @Security  Bearer
// @Success 200 {object} pocket.SuccessResponse
// @Router /pockets/{id} [delete]
func (h *handler) DeletePocket(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
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

	err = balance.Add(h.DB, p.AccountID, p.Balance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + err.Error()})
	}

	tx := h.DB.Delete(&Pocket{}, id)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + tx.Error.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(SuccessResponse{Message: "delete pocket success"})

}
