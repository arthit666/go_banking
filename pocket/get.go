package pocket

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// @Summary Get all pocket
// @Description Get details of all pocket
// @Tags pockets
// @Accept  json
// @Produce  json
// @Security  Bearer
// @Success 200 {array} pocket.Pocket
// @Router /pockets/ [get]
func (h *handler) GetAllPockets(c *fiber.Ctx) error {
	p := []Pocket{}
	acc := c.Locals("account_id").(int)
	accStr := strconv.Itoa(acc)
	tx := h.DB.Where("account_id = ?", accStr).Find(&p)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + tx.Error.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(p)
}

// @Summary Get pocket by ID
// @Description Get pocket by ID
// @Tags pockets
// @Produce json
// @Param id path int true "Pocket ID"
// @Security  Bearer
// @Router /pockets/{id} [get]
func (h *handler) GetPocketById(c *fiber.Ctx) error {
	id := c.Params("id")
	acc := c.Locals("account_id").(int)
	accStr := strconv.Itoa(acc)
	p, err := getById(id, accStr, h)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Err{Message: "pocket not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(p)
}

func getById(id string, accStr string, h *handler) (*Pocket, error) {
	p := &Pocket{}
	tx := h.DB.Where("account_id = ?", accStr).First(p, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return p, nil
}
