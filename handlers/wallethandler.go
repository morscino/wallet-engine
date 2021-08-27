package handlers

import (
	"time"

	"github.com/google/uuid"
	"github.com/morscino/wallet-engine/models/walletmodel"
	"github.com/morscino/wallet-engine/service/walletservice"
)

type WalletHandler struct {
	WalletService walletservice.WalletRepository
}

func NewWalletHandler(w walletservice.WalletRepository) WalletHandler {
	return WalletHandler{WalletService: w}

}

func (w WalletHandler) Createwallet(userId, currency string) walletmodel.Wallet {

	id := uuid.New()

	wallet := walletmodel.Wallet{
		ID:        id,
		UserID:    userId,
		Currency:  currency,
		Balance:   0,
		Disabled:  false,
		CreatedAt: time.Now(),
	}

	w.WalletService.CreateWallet(wallet)

	return wallet

}

func (w WalletHandler) DebitWallet(wallet walletmodel.Wallet, amount float64) {
	newWalletBalance := wallet.Balance - amount

	w.WalletService.UpdateWalletBalance(wallet.UserID, newWalletBalance)

}

func (w WalletHandler) CreditWallet(wallet walletmodel.Wallet, amount float64) {
	newWalletBalance := wallet.Balance + amount

	w.WalletService.UpdateWalletBalance(wallet.UserID, newWalletBalance)

}
