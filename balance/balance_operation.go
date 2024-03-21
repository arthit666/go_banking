package balance

import (
	"fmt"

	"gorm.io/gorm"
)

type Account struct {
	ID      uint
	Balance float64
}

func Deduct(db *gorm.DB, accountID uint, amount float64) error {
	var acc Account
	tx := db.First(&acc, accountID)
	if tx.Error != nil {
		return fmt.Errorf("error getting account: %s", tx.Error.Error())
	}

	if acc.Balance < amount {
		return fmt.Errorf("insufficient balance in account")
	}

	acc.Balance -= amount

	tx = db.Save(&acc)
	if tx.Error != nil {
		return fmt.Errorf("error updating account balance: %s", tx.Error.Error())
	}

	return nil
}

func Add(db *gorm.DB, accountID uint, amount float64) error {
	var acc Account
	tx := db.First(&acc, accountID)
	if tx.Error != nil {
		return fmt.Errorf("error getting account: %s", tx.Error.Error())
	}

	acc.Balance += amount

	tx = db.Save(&acc)
	if tx.Error != nil {
		return fmt.Errorf("error updating account balance: %s", tx.Error.Error())
	}

	return nil
}
