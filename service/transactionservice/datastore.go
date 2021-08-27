package transactionservice

import (
	"github.com/morscino/wallet-engine/models/transactionmodel"
	"github.com/morscino/wallet-engine/utility/log"
	"github.com/morscino/wallet-engine/utility/message"
	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
}

type TransactionRepository interface {
	CreateTransaction(txn transactionmodel.Transaction) string
	UpdateTransaction(ref, walletId string) transactionmodel.Transaction
}

func NewTransactionService(database *gorm.DB) TransactionRepository {
	return &TransactionService{db: database}
}

func (t TransactionService) CreateTransaction(txn transactionmodel.Transaction) string {
	newTxn := t.db.Create(&txn)

	if newTxn.Error != nil {

		log.Error(message.DB_ERROR_OCCURED.String(), newTxn.Error)

	}

	return txn.TransRef

}

func (t TransactionService) UpdateTransaction(ref, walletId string) transactionmodel.Transaction {

	var transaction transactionmodel.Transaction
	c := t.db.Model(&transaction).Where("trans_ref=?", ref).Where("wallet_id=?", walletId).Update("response_code", "000")

	if c.Error != nil {

		log.Error(message.DB_ERROR_OCCURED.String(), c.Error)

	}
	transaction.TransRef = ref
	//fmt.Println("---+++----", transaction.Amount)
	return transaction
}
