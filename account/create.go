package account

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Create a new account
// @Description Create a new account with the input payload
// @Tags accounts
// @Accept json
// @Produce json
// @Security  Bearer
// @Param account body account.AccountRequest true "Create account"
// @Success 201 {object} account.SuccessResponse
// @Router /accounts/ [post]
func (h *handler) CreateAccount(c *fiber.Ctx) error {
	a := &Account{}

	if err := c.BodyParser(a); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(a); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(Err{Massage: "payload invalid: " + err.Error()})
	}
	a.AccountNumber = randomAccountNumber()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}
	a.Password = string(hashedPassword)

	_, err = create(h, a)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Err{Massage: "error: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{Message: "create account success"})
}

func randomAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	id := ""
	for i := 0; i < 10; i++ {
		id += strconv.Itoa(rand.Intn(10))
	}
	return id
}

func create(h *handler, a *Account) (*Account, error) {
	tx := h.DB.Create(a)
	if tx.Error != nil {
		return nil, fmt.Errorf(tx.Error.Error())
	}
	return a, nil
}
