package account

import (
	"errors"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// @Summary Get all accounts
// @Description Get details of all accounts
// @Tags accounts
// @Accept  json
// @Produce  json
// @Security  Bearer
// @Success 200 {array} AccountResponseList
// @Router /accounts/ [get]
func (h *handler) GetAllAccounts(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Err{Massage: "invalid page parameter"})
	}
	limit, err := strconv.Atoi(c.Query("count", "10"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Err{Massage: "invalid page_size parameter"})
	}

	offset := (page - 1) * limit

	acc := []Account{}
	tx := h.DB.Offset(offset).Limit(limit).Preload("PocketList").Find(&acc)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + tx.Error.Error()})
	}

	var totalCount int64
	h.DB.Model(&Account{}).Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	accr := []AccountResponse{}

	for _, v := range acc {
		accr = append(accr, AccountResponse{
			ID:            v.ID,
			Email:         v.Email,
			AccountNumber: v.AccountNumber,
			PocketList:    v.PocketList,
			Balance:       v.Balance,
		})
	}

	resp := AccountResponseList{
		AccountList: accr,
		Page:        page,
		Count:       limit,
		TotalPage:   totalPages,
		TotalCount:  totalCount,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// @Summary Get account by ID
// @Description Get account by ID
// @Tags accounts
// @Produce json
// @Param id path int true "Account ID"
// @Security  Bearer
// @Success 200 {array} account.AccountResponse
// @Router /accounts/{id} [get]
func (h *handler) GetAccountById(c *fiber.Ctx) error {
	id := c.Params("id")
	acc, err := getById(id, h)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Err{Massage: "account not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}
	res := AccountResponse{
		ID:            acc.ID,
		Email:         acc.Email,
		AccountNumber: acc.AccountNumber,
		PocketList:    acc.PocketList,
		Balance:       acc.Balance,
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func getById(id string, h *handler) (*Account, error) {
	p := &Account{}
	tx := h.DB.Preload("PocketList").First(p, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return p, nil
}

func getByAccountNumber(an string, h *handler) (*Account, error) {
	p := &Account{}
	tx := h.DB.Preload("PocketList").First(p, "account_number", an)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return p, nil
}
