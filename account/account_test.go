package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arthit666/make_app/pocket"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateAccount(t *testing.T) {
	//Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	handler := New(tx)
	app.Post("/accounts", handler.CreateAccount)

	reqBody := AccountRequest{
		Email:    "test@example.com",
		Password: "password123",
		Balance:  1000,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Acc
	resp, err := app.Test(req)
	assert.NoError(t, err)

	//Assertions
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	var responseBody SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "create account success", responseBody.Message)

	var createdAccount Account
	tx.Where("email = ?", reqBody.Email).First(&createdAccount)
	assert.NotEmpty(t, createdAccount.ID)
	assert.Equal(t, reqBody.Email, createdAccount.Email)
	assert.Equal(t, reqBody.Balance, createdAccount.Balance)
	assert.NotEmpty(t, createdAccount.AccountNumber)

	err = bcrypt.CompareHashAndPassword([]byte(createdAccount.Password), []byte(reqBody.Password))
	assert.NoError(t, err)

}

func TestGetAllAccounts(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{}, &pocket.Pocket{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	handler := New(tx)
	app.Get("/accounts", handler.GetAllAccounts)

	accounts := []Account{
		{Email: "user1@example.com", AccountNumber: "1234567890", Balance: 1000},
		{Email: "user2@example.com", AccountNumber: "9876543210", Balance: 500},
	}
	for _, account := range accounts {
		tx.Create(&account)
	}

	pockets := []pocket.Pocket{
		{Title: "Pocket 1", Balance: 100, AccountID: 1},
		{Title: "Pocket 2", Balance: 200, AccountID: 2},
	}
	for _, p := range pockets {
		tx.Create(&p)
	}

	// Act
	req := httptest.NewRequest(http.MethodGet, "/accounts?page=1&count=10", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody AccountResponseList
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	assert.Equal(t, 1, responseBody.Page)
	assert.Equal(t, 10, responseBody.Count)
	assert.Equal(t, 1, responseBody.TotalPage)
	assert.Equal(t, int64(2), responseBody.TotalCount)

	assert.Equal(t, 2, len(responseBody.AccountList))

	assert.Equal(t, accounts[0].Email, responseBody.AccountList[0].Email)
	assert.Equal(t, accounts[0].AccountNumber, responseBody.AccountList[0].AccountNumber)
	assert.Equal(t, accounts[0].Balance, responseBody.AccountList[0].Balance)
	assert.Equal(t, 1, len(responseBody.AccountList[0].PocketList))
	assert.Equal(t, pockets[0].Title, responseBody.AccountList[0].PocketList[0].Title)
	assert.Equal(t, pockets[0].Balance, responseBody.AccountList[0].PocketList[0].Balance)

}
func TestGetAccountById(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{}, &pocket.Pocket{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	handler := New(tx)
	app.Get("/accounts/:id", handler.GetAccountById)

	account := Account{
		Email:         "user@example.com",
		AccountNumber: "1234567890",
		Balance:       1000,
	}
	tx.Create(&account)

	pockets := []pocket.Pocket{
		{Title: "Pocket 1", Balance: 100, AccountID: 1},
		{Title: "Pocket 2", Balance: 200, AccountID: 1},
	}
	for _, p := range pockets {
		tx.Create(&p)
	}

	// Act
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", account.ID), nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody AccountResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	assert.Equal(t, account.ID, responseBody.ID)
	assert.Equal(t, account.Email, responseBody.Email)
	assert.Equal(t, account.AccountNumber, responseBody.AccountNumber)
	assert.Equal(t, account.Balance, responseBody.Balance)
	assert.Equal(t, 2, len(responseBody.PocketList))

	assert.Equal(t, pockets[0].Title, responseBody.PocketList[0].Title)
	assert.Equal(t, pockets[0].Balance, responseBody.PocketList[0].Balance)
	assert.Equal(t, pockets[1].Title, responseBody.PocketList[1].Title)
	assert.Equal(t, pockets[1].Balance, responseBody.PocketList[1].Balance)
}

func TestLogin(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	handler := New(tx)
	app.Post("/login", handler.Login)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	account := Account{
		Email:    "user@example.com",
		Password: string(hashedPassword),
	}
	tx.Create(&account)

	reqBody := Login{
		Email:    "user@example.com",
		Password: password,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	// Act
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody.AccessToken)
	assert.NotEmpty(t, responseBody.RefreshToken)

}

func TestAccountTransfer(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{}, &AccountTransfer{}, &pocket.Pocket{})
	assert.NoError(t, err)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("account_id", 1)
		return c.Next()
	})
	handler := New(db)
	app.Config()
	app.Post("/accounts/transfer", handler.Transfer)

	fromAccount := &Account{
		Email:         "test@test.com",
		AccountNumber: "1234567890",
		Balance:       1000,
	}
	db.Create(fromAccount)

	toAccount := &Account{
		Email:         "test2@test.com",
		AccountNumber: "9876543210",
		Balance:       500,
	}
	db.Create(toAccount)

	reqBody := AccountTransferRequest{
		To:     toAccount.AccountNumber,
		Amount: 500,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("account_id", fromAccount.ID)
		return c.Next()
	})

	//Acc
	resp, err := app.Test(req)
	assert.NoError(t, err)

	//Assertions
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	var updatedFromAccount Account
	db.First(&updatedFromAccount, fromAccount.ID)
	assert.Equal(t, 500.0, updatedFromAccount.Balance)

	var updatedToAccount Account
	db.First(&updatedToAccount, toAccount.ID)
	assert.Equal(t, 1000.0, updatedToAccount.Balance)

	var transferRecord AccountTransfer
	db.First(&transferRecord)
	assert.Equal(t, fromAccount.AccountNumber, transferRecord.From)
	assert.Equal(t, toAccount.AccountNumber, transferRecord.To)
	assert.Equal(t, 500.0, transferRecord.Amount)

}
