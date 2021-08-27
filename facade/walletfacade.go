package facade

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morscino/wallet-engine/handlers"
	"github.com/morscino/wallet-engine/models/walletmodel"
)

type WalletFacade struct {
	ctx           context.Context
	WalletHandler handlers.WalletHandler
}

func NewWalletFacade(w handlers.WalletHandler, ctx context.Context) *WalletFacade {
	return &WalletFacade{
		ctx:           ctx,
		WalletHandler: w,
	}
}

func (w WalletFacade) Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong-pong-pong",
	})
}

func (w WalletFacade) CreateWallet(c *gin.Context) {

	var input walletmodel.WalletRegistrationData
	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	//ensure userid is unique
	if !w.userIdIsUnique(strings.ToLower(input.UserID)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet with phone number already exists"})
		return
	}

	//Do some check before creation

	newWallet := w.WalletHandler.Createwallet(input.UserID, input.Currency)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Wallet successfully created", "data": newWallet})

}

func (w WalletFacade) DebitCreditWallet(c *gin.Context) {
	var input walletmodel.WalletTransactionInput
	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	//get debit wallet
	debitWallet := w.WalletHandler.WalletService.GetWalletByUserId(input.FromWallet)
	creditWallet := w.WalletHandler.WalletService.GetWalletByUserId(input.ToWallet)

	//check if enabled
	if debitWallet.Disabled {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusOK, "message": "Wallet is disabled, please contact the administrator"})
		return
	}

	tranxAmount := float64(input.Amount) * 100

	//check if it has sufficient fund
	if debitWallet.Balance < (tranxAmount) {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusOK, "message": "Insufficient funds"})
		return
	}

	//debit from wallet
	w.WalletHandler.DebitWallet(debitWallet, tranxAmount)
	//credit to wallet
	w.WalletHandler.CreditWallet(creditWallet, tranxAmount)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "transaction done"})

}

func (w WalletFacade) userIdIsUnique(userId string) bool {
	exists := true

	wallet := w.WalletHandler.WalletService.GetWalletByUserId(userId)
	if wallet.UserID != "" {
		exists = false
	}

	return exists
}
