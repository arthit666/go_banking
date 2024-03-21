package pocket

import (
	"time"

	"gorm.io/gorm"
)

type Pocket struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"create_at"`
	UpdatedAt   time.Time      `json:"update_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `json:"title" validate:"required"`
	Balance     float64        `json:"balance"`
	Description *string        `json:"description"`
	AccountID   uint           `json:"-" validate:"required"`
}

type PocketUpdate struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type PocketCreate struct {
	Title       string  `json:"title" validate:"required"`
	Balance     float64 `json:"balance"`
	Description *string `json:"description"`
}

type PocketTransfer struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"create_at"`
	From      uint      `json:"from" validate:"required"`
	To        uint      `json:"to" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,numeric,gt=0"`
	AccountID uint      `json:"-"`
}

type PocketTransferRequest struct {
	From   uint    `json:"from" validate:"required"`
	To     uint    `json:"to" validate:"required"`
	Amount float64 `json:"amount" validate:"required,numeric,gt=0"`
}

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *handler {
	return &handler{db}
}

type Err struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
