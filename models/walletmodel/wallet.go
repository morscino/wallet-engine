package walletmodel

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Balance     int       `json:"balance" gorm:"type:bigint;not null"`
	Currency    string    `json:"currency"  gorm:"type:varchar(50);not null"`
	PhoneNumber string    `json:"phoneNumber"  gorm:"type:varchar(100);not null"`
	Disabled    bool      `json:"disabled" gorm:"type:bool;not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type WalletRegistrationData struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Currency    string `json:"currency" binding:"required"`
	CountryCode string `json:"countryCode" binding:"required" `
}

type WalletTransactionInput struct {
	FromWallet string `json:"fromWallet" binding:"required"`
	ToWallet   string `json:"toWallet" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

type WalletStatusInput struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}
