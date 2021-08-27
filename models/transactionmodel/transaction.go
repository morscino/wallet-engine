package transactionmodel

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	WalletId     string    `gorm:"type:varchar(50);not null"`
	ResponseCode string    `gorm:"type:varchar(10);not null"`
	Amount       int       `gorm:"type:bigint;not null"`
	TransRef     string    `gorm:"type:varchar(100);not null"`
	TransType    string    `gorm:"type:varchar(10);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
