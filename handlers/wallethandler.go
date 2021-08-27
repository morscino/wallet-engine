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

var e error

func (w WalletHandler) Createwallet(phoneNumber, currency string) (walletmodel.Wallet, error) {

	id := uuid.New()

	wallet := walletmodel.Wallet{
		ID:          id,
		PhoneNumber: phoneNumber,
		Currency:    currency,
		Balance:     500000,
		Disabled:    false,
		CreatedAt:   time.Now(),
	}

	wa, err := w.WalletService.CreateWallet(wallet)

	if err != nil {
		e = err
	}

	return wa, e

}

func (w WalletHandler) DebitCreditWallet(dr, cr walletmodel.Wallet, amount int) error {
	err := w.debitWallet(dr, amount)
	if err == nil {
		err = w.creditWallet(cr, amount)
	}

	return err

}

func (w WalletHandler) debitWallet(wallet walletmodel.Wallet, amount int) error {
	newWalletBalance := wallet.Balance - amount

	err := w.WalletService.UpdateWalletBalance(wallet.PhoneNumber, newWalletBalance)

	if err != nil {
		e = err
	}
	return e

	//send debit notification

}

func (w WalletHandler) creditWallet(wallet walletmodel.Wallet, amount int) error {
	newWalletBalance := wallet.Balance + amount

	err := w.WalletService.UpdateWalletBalance(wallet.PhoneNumber, newWalletBalance)

	if err != nil {
		e = err
	}
	return e

	//send credit notification

}

func (w WalletHandler) EnableWallet(phoneNumber string) error {

	err := w.WalletService.UpdateWalletStatus(phoneNumber, false)

	if err != nil {
		e = err
	}

	return e

}

func (w WalletHandler) DisableWallet(phoneNumber string) error {

	err := w.WalletService.UpdateWalletStatus(phoneNumber, true)

	if err != nil {
		e = err
	}

	return e

}
