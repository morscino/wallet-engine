package handlers

import (
	"fmt"
	"math/rand"

	"github.com/morscino/wallet-engine/models/transactionmodel"
	"github.com/morscino/wallet-engine/service/transactionservice"
)

type TransactionHandler struct {
	TransactionService transactionservice.TransactionRepository
}

func (t TransactionHandler) GenerateTransRef() string {

	s := fmt.Sprintf("%016d", rand.Int63n(1e16))

	return "OPAY-" + s

}

func NewTransactionHandler(t transactionservice.TransactionRepository) TransactionHandler {
	return TransactionHandler{TransactionService: t}
}

func (t TransactionHandler) CreateTransaction(transaction transactionmodel.Transaction) string {
	ref := t.TransactionService.CreateTransaction(transaction)
	return ref
}

func (t TransactionHandler) ConcludeTransaction(ref, walletId string) transactionmodel.Transaction {
	transaction := t.TransactionService.UpdateTransaction(ref, walletId)

	return transaction
}
