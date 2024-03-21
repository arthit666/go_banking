package pocket

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// @Summary Transfer funds between pockets
// @Description Transfer funds from one pocket to another with pocket id
// @Tags pockets
// @Accept json
// @Produce json
// @Param transfer body pocket.PocketTransferRequest true "PocketTransferRequest data"
// @Success 201 {object} pocket.SuccessResponse
// @Security  Bearer
// @Router /pockets/transfer/ [post]
func (h *handler) Transfer(c *fiber.Ctx) error {

	tr := &PocketTransferRequest{}

	if err := c.BodyParser(tr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(tr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Err{Message: "payload invalid: " + err.Error()})
	}

	acc := c.Locals("account_id").(int)
	accStr := strconv.Itoa(acc)

	t := &PocketTransfer{
		From:      tr.From,
		To:        tr.To,
		Amount:    tr.Amount,
		AccountID: uint(acc),
	}

	fpock, err := getById(strconv.Itoa(int(t.From)), accStr, h)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Err{Message: "from pocket not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + err.Error()})
	}

	tpock, err := getById(strconv.Itoa(int(t.To)), accStr, h)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Err{Message: "target pocket not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + err.Error()})
	}

	if err = transferBalance(h, fpock, tpock, t.Amount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + err.Error()})
	}

	tx := h.DB.Create(t)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Message: "error: " + tx.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{Message: "transfer success"})
}

func transferBalance(h *handler, from, to *Pocket, amount float64) error {
	amountDec := decimal.NewFromFloat(amount)

	bl := decimal.NewFromFloat(from.Balance)
	if bl.LessThan(amountDec) {
		return fmt.Errorf("insufficient balance in source pocket")
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	bl = bl.Sub(amountDec)
	from.Balance, _ = bl.Float64()
	if err := tx.Save(from).Error; err != nil {
		tx.Rollback()
		return err
	}

	bl = decimal.NewFromFloat(to.Balance)
	bl = bl.Add(amountDec)
	to.Balance, _ = bl.Float64()
	if err := tx.Save(to).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
