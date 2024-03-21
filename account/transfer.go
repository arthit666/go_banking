package account

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// @Summary Transfer funds between accounts
// @Description Transfer funds from one accounts to another with account number
// @Tags accounts
// @Accept json
// @Produce json
// @Param transfer body account.AccountTransferRequest true "AccountTransferRequest data"
// @Success 201 {object} account.SuccessResponse
// @Security  Bearer
// @Router /accounts/transfer/ [post]
func (h *handler) Transfer(c *fiber.Ctx) error {
	tr := &AccountTransferRequest{}

	if err := c.BodyParser(tr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(tr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Err{Massage: "payload invalid: " + err.Error()})
	}

	t := &AccountTransfer{
		To:     tr.To,
		Amount: tr.Amount,
	}

	acc := c.Locals("account_id").(int)

	accStr := strconv.Itoa(acc)

	fpock, err := getById(accStr, h)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Err{Massage: "from account not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}

	t.From = fpock.AccountNumber

	tpock, err := getByAccountNumber(t.To, h)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Err{Massage: "target account not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}

	if err = transferBalance(h, fpock, tpock, t.Amount); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}

	tx := h.DB.Create(t)
	if tx.Error != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + tx.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{Message: "transfer success"})
}

func transferBalance(h *handler, from, to *Account, amount float64) error {
	amountDec := decimal.NewFromFloat(amount)

	bl := decimal.NewFromFloat(from.Balance)
	if bl.LessThan(amountDec) {
		return fmt.Errorf("insufficient balance in source account")
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
