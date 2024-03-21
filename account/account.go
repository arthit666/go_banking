package account

import (
	"time"

	"github.com/arthit666/make_app/pocket"
	"gorm.io/gorm"
)

type Account struct {
	ID            uint            `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time       `json:"create_at"`
	UpdatedAt     time.Time       `json:"update_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
	Email         string          `json:"email" validate:"required,email" gorm:"unique"`
	Password      string          `json:"password" validate:"required"`
	Balance       float64         `json:"balance"`
	AccountNumber string          `json:"account_number"`
	PocketList    []pocket.Pocket `gorm:"ForeignKey:AccountID"`
}

type AccountResponse struct {
	ID            uint            `json:"id"`
	Email         string          `json:"email"`
	AccountNumber string          `json:"account_number"`
	Balance       float64         `json:"balance"`
	PocketList    []pocket.Pocket `json:"pocket_list"`
}

type AccountResponseList struct {
	AccountList []AccountResponse `json:"result"`
	Page        int               `json:"page"`
	TotalPage   int               `json:"total_page"`
	Count       int               `json:"count"`
	TotalCount  int64             `json:"total_count"`
}

type AccountRequest struct {
	Email    string  `json:"email" validate:"required"`
	Password string  `json:"password" validate:"required"`
	Balance  float64 `json:"balance"`
}

type AccountTransfer struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"create_at"`
	From      string    `json:"from"`
	To        string    `json:"to" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,numeric,gt=0"`
}

type AccountTransferRequest struct {
	To     string  `json:"to" validate:"required"`
	Amount float64 `json:"amount" validate:"required,numeric,gt=0"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *handler {
	return &handler{db}
}

type Err struct {
	Massage string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
