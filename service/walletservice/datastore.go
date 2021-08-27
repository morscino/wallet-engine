package walletservice

import (
	"github.com/morscino/wallet-engine/models/walletmodel"
	"github.com/morscino/wallet-engine/utility/log"
	"github.com/morscino/wallet-engine/utility/message"
	"gorm.io/gorm"
)

type WalletService struct {
	db *gorm.DB
}

var e error

type WalletRepository interface {
	CreateWallet(wm walletmodel.Wallet) (walletmodel.Wallet, error)
	GetWalletByPhoneNumber(phoneNumber string) walletmodel.Wallet
	UpdateWalletBalance(phoneNumber string, amount int) error
	UpdateWalletStatus(phoneNumber string, status bool) error
}

func NewWalletervice(database *gorm.DB) WalletRepository {
	return &WalletService{db: database}
}

func (w WalletService) CreateWallet(wallet walletmodel.Wallet) (walletmodel.Wallet, error) {

	newWallet := w.db.Create(&wallet)

	if newWallet.Error != nil {

		log.Error(message.DB_ERROR_OCCURED.String(), newWallet.Error)
		e = newWallet.Error
	}

	return wallet, e

}

func (w WalletService) GetWalletByPhoneNumber(phoneNumber string) walletmodel.Wallet {

	var wallet walletmodel.Wallet
	wa := w.db.Where("phone_number=?", phoneNumber).Find(&wallet)
	if wa.Error != nil {
		log.Error(message.DB_ERROR_OCCURED.String(), wa.Error)

		panic(message.DATABASE_ERROR.String())
	}

	return wallet
}

func (w WalletService) UpdateWalletBalance(phoneNumber string, amount int) error {

	var wallet walletmodel.Wallet
	u := w.db.Model(&wallet).Where("phone_number=?", phoneNumber).Update("balance", amount)
	if u.Error != nil {
		log.Error(message.DB_ERROR_OCCURED.String(), u.Error)

		e = u.Error
		panic(message.DATABASE_ERROR.String())
	}

	return e

}

func (w WalletService) UpdateWalletStatus(phoneNumber string, status bool) error {
	var wallet walletmodel.Wallet
	u := w.db.Model(&wallet).Where("phone_number=?", phoneNumber).Update("disabled", status)

	if u.Error != nil {
		log.Error(message.DB_ERROR_OCCURED.String(), u.Error)
		e = u.Error
		panic(message.DATABASE_ERROR.String())
	}

	return e

}
