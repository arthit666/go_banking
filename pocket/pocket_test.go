package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Account struct {
	ID      uint
	Email   string
	Balance float64
}

func TestCreatePocket(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(Account{}, &Pocket{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("account_id", 1)
		return c.Next()
	})
	handler := New(tx)
	app.Post("/pockets", handler.CreatePocket)

	account := Account{
		Email:   "test@example.com",
		Balance: 1000,
	}
	tx.Create(&account)

	reqBody := PocketCreate{
		Title:   "Pocket 1",
		Balance: 500,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pockets", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var responseBody SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "create pocket success", responseBody.Message)

	var createdPocket Pocket
	tx.First(&createdPocket)
	assert.Equal(t, reqBody.Title, createdPocket.Title)
	assert.Equal(t, reqBody.Balance, createdPocket.Balance)
	assert.Equal(t, account.ID, createdPocket.AccountID)
}

func TestGetAllPockets(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{}, &Pocket{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("account_id", 1)
		return c.Next()
	})
	handler := New(tx)
	app.Get("/pockets", handler.GetAllPockets)

	account := Account{
		Email:   "test@example.com",
		Balance: 1000,
	}
	tx.Create(&account)

	pockets := []Pocket{
		{Title: "Pocket 1", Balance: 100, AccountID: 1},
		{Title: "Pocket 2", Balance: 200, AccountID: 1},
	}
	for _, pocket := range pockets {
		tx.Create(&pocket)
	}

	req := httptest.NewRequest(http.MethodGet, "/pockets", nil)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody []Pocket
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, len(pockets), len(responseBody))

	for i, pocket := range pockets {
		assert.Equal(t, pocket.Title, responseBody[i].Title)
		assert.Equal(t, pocket.Balance, responseBody[i].Balance)

	}
}

func TestGetPocketById(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{}, &Pocket{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("account_id", 1)
		return c.Next()
	})
	handler := New(tx)
	app.Get("/pockets/:id", handler.GetPocketById)

	account := Account{
		Email:   "test@example.com",
		Balance: 1000,
	}
	tx.Create(&account)

	pocket := Pocket{
		Title:     "Pocket 1",
		Balance:   100,
		AccountID: 1,
	}
	tx.Create(&pocket)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/pockets/%d", 1), nil)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody Pocket
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, pocket.Title, responseBody.Title)
	assert.Equal(t, pocket.Balance, responseBody.Balance)
}

func TestUpdatePocket(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{}, &Pocket{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("account_id", 1)
		return c.Next()
	})
	handler := New(tx)
	app.Put("/pockets/:id", handler.UpdatePocket)

	account := Account{
		Email:   "test@example.com",
		Balance: 1000,
	}
	tx.Create(&account)

	pocket := Pocket{
		Title:     "Pocket 1",
		Balance:   100,
		AccountID: 1,
	}
	tx.Create(&pocket)

	payload := PocketUpdate{
		Title: "Updated Pocket",
	}
	jsonPayload, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/pockets/%d", pocket.ID), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "update pocket success", responseBody.Message)

	var updatedPocket Pocket
	err = tx.First(&updatedPocket, pocket.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated Pocket", updatedPocket.Title)
}

func TestDeletePocket(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{}, &Pocket{})
	assert.NoError(t, err)

	tx := db.Begin()
	defer tx.Rollback()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("account_id", 1)
		return c.Next()
	})
	handler := New(tx)
	app.Delete("/pockets/:id", handler.DeletePocket)

	account := Account{
		Email:   "test@example.com",
		Balance: 1000,
	}
	tx.Create(&account)

	pocket := Pocket{
		Title:     "Pocket 1",
		Balance:   100,
		AccountID: 1,
	}
	tx.Create(&pocket)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/pockets/%d", pocket.ID), nil)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "delete pocket success", responseBody.Message)

	var updatedAccount Account
	err = tx.First(&updatedAccount, account.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, account.Balance+pocket.Balance, updatedAccount.Balance)
}

func TestTransfer(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Account{}, &Pocket{}, &PocketTransfer{})
	assert.NoError(t, err)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("account_id", 1)
		return c.Next()
	})
	handler := New(db)
	app.Post("/pockets/transfer", handler.Transfer)

	account := Account{
		Email:   "test@example.com",
		Balance: 1000,
	}
	db.Create(&account)

	fromPocket := Pocket{
		Title:     "From Pocket",
		Balance:   500,
		AccountID: 1,
	}
	db.Create(&fromPocket)

	toPocket := Pocket{
		Title:     "To Pocket",
		Balance:   200,
		AccountID: 1,
	}
	db.Create(&toPocket)

	payload := PocketTransferRequest{
		From:   1,
		To:     2,
		Amount: 100,
	}

	jsonPayload, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pockets/transfer", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var responseBody SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "transfer success", responseBody.Message)

	var updatedFromPocket Pocket
	err = db.First(&updatedFromPocket, fromPocket.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, float64(400), updatedFromPocket.Balance)

	var updatedToPocket Pocket
	err = db.First(&updatedToPocket, toPocket.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, float64(300), updatedToPocket.Balance)

}
