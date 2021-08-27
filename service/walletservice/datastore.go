package walletservice

import (
	"log"

	"github.com/morscino/wallet-engine/models/walletmodel"
	"gorm.io/gorm"
)

type WalletService struct {
	db *gorm.DB
}

type WalletRepository interface {
	CreateWallet(wm walletmodel.Wallet) walletmodel.Wallet
	GetWalletByUserId(userId string) walletmodel.Wallet
	UpdateWalletBalance(userId string, amount float64)
}

func NewWalletervice(database *gorm.DB) WalletRepository {
	return &WalletService{db: database}
}

func (w WalletService) CreateWallet(wallet walletmodel.Wallet) walletmodel.Wallet {

	newWallet := w.db.Create(&wallet)

	if newWallet.Error != nil {

		log.Fatal("there was an error")
	}

	return wallet

}

func (w WalletService) GetWalletByUserId(userId string) walletmodel.Wallet {

	var wallet walletmodel.Wallet
	w.db.Where("user_id=?", userId).Find(&wallet)

	return wallet
}

func (w WalletService) UpdateWalletBalance(userId string, amount float64) {
	var wallet walletmodel.Wallet
	w.db.Model(&wallet).Where("user_id=?", userId).Update("balance", amount)
	//db.Model(&user).Update("name", "hello")

}
