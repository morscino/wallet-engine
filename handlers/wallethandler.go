package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/morscino/wallet-engine/models/transactionmodel"
	"github.com/morscino/wallet-engine/models/walletmodel"
	"github.com/morscino/wallet-engine/service/walletservice"
)

type WalletHandler struct {
	WalletService walletservice.WalletRepository
	TransactionHandler
}

func NewWalletHandler(w walletservice.WalletRepository, t TransactionHandler) WalletHandler {
	return WalletHandler{WalletService: w, TransactionHandler: t}

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

func (w WalletHandler) DebitCreditWallet(dr, cr walletmodel.Wallet, amount int) (transactionmodel.Transaction, error) {
	txn, err := w.debitWallet(dr, amount)
	if err == nil {
		txn, err = w.creditWallet(cr, amount, txn)
	}

	return txn, err

}

func (w WalletHandler) debitWallet(wallet walletmodel.Wallet, amount int) (transactionmodel.Transaction, error) {
	newWalletBalance := wallet.Balance - amount

	id := uuid.New()
	//Initialize debit tranx
	transaction := transactionmodel.Transaction{
		ID:           id,
		WalletId:     wallet.PhoneNumber,
		ResponseCode: "039",
		Amount:       amount,
		TransRef:     w.TransactionHandler.GenerateTransRef(),
		TransType:    "DR",
		CreatedAt:    time.Now(),
	}

	//fmt.Println(transaction.ResponseCode)

	ref := w.TransactionHandler.CreateTransaction(transaction)

	//fmt.Println(ref)

	err := w.WalletService.UpdateWalletBalance(wallet.PhoneNumber, newWalletBalance)

	if err != nil {
		e = err
	}

	//conclude debit transaction
	txn := w.TransactionHandler.ConcludeTransaction(ref, wallet.PhoneNumber)
	return txn, e

	//send debit notification

}

func (w WalletHandler) creditWallet(wallet walletmodel.Wallet, amount int, txn transactionmodel.Transaction) (transactionmodel.Transaction, error) {
	newWalletBalance := wallet.Balance + amount

	id := uuid.New()
	//Initialize debit tranx
	transaction := transactionmodel.Transaction{
		ID:           id,
		WalletId:     wallet.PhoneNumber,
		ResponseCode: "039",
		Amount:       amount,
		TransRef:     txn.TransRef,
		TransType:    "CR",
		CreatedAt:    time.Now(),
	}

	fmt.Println("---------", txn.TransRef)

	ref := w.TransactionHandler.CreateTransaction(transaction)

	err := w.WalletService.UpdateWalletBalance(wallet.PhoneNumber, newWalletBalance)

	if err != nil {
		e = err
	}

	//conclude debit transaction
	trans := w.TransactionHandler.ConcludeTransaction(ref, wallet.PhoneNumber)
	return trans, e

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
