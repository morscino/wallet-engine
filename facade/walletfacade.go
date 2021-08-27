package facade

import (
	"context"
	"net/http"
	"strconv"
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
	wallet := walletmodel.Wallet{}

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	//get debit and credit wallets
	debitWallet := w.WalletHandler.WalletService.GetWalletByUserId(input.FromWallet)
	creditWallet := w.WalletHandler.WalletService.GetWalletByUserId(input.ToWallet)

	if wallet == debitWallet {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Debit wallet does not exist"})
		return
	}

	if wallet == creditWallet {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Credit wallet does not exist"})
		return
	}

	//check if enabled
	if debitWallet.Disabled {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusOK, "message": "Wallet is disabled, please contact the administrator"})
		return
	}

	tranxAmount, _ := strconv.ParseFloat(input.Amount, 64)

	//check for negative amount
	if tranxAmount < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid amount"})
		return
	}

	tranxAmount *= 100
	//check if it debit wallet has sufficient fund
	if debitWallet.Balance < tranxAmount {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusOK, "message": "Insufficient funds"})
		return
	}

	//debit-credit wallets
	w.WalletHandler.DebitCreditWallet(debitWallet, creditWallet, tranxAmount)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "transaction done"})

}

func (w WalletFacade) EnableWallet(c *gin.Context) {

	var input walletmodel.WalletStatusInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}
	enabled := w.WalletHandler.EnableWallet(input.UserID)
	_ = enabled

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Wallet Successfully Enabled"})

}

func (w WalletFacade) DisableWallet(c *gin.Context) {

	var input walletmodel.WalletStatusInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}
	enabled := w.WalletHandler.DisableWallet(input.UserID)
	_ = enabled

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Wallet Successfully Disabled"})

}

func (w WalletFacade) userIdIsUnique(userId string) bool {
	exists := true

	wallet := w.WalletHandler.WalletService.GetWalletByUserId(userId)
	if wallet.UserID != "" {
		exists = false
	}

	return exists
}
