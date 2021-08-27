package facade

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/morscino/wallet-engine/handlers"
	"github.com/morscino/wallet-engine/models/walletmodel"
	"github.com/morscino/wallet-engine/utility/message"
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

	phoneNumber := input.CountryCode + input.PhoneNumber

	//validate phone number

	//ensure PhoneNumber is unique
	if w.walletExists(phoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": message.WALLET_ALREADY_EXISTS.String()})
		return
	}

	//Do some check before creation

	newWallet, err := w.WalletHandler.Createwallet(phoneNumber, input.Currency)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": message.WALLET_NOT_SUCCESSFULLY_CREATED.String(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": message.WALLET_SUCCESSFULLY_CREATED.String(), "data": newWallet})

}

func (w WalletFacade) DebitCreditWallet(c *gin.Context) {
	var input walletmodel.WalletTransactionInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	//get debit and credit wallets
	debitWallet := w.WalletHandler.WalletService.GetWalletByPhoneNumber(input.FromWallet)
	creditWallet := w.WalletHandler.WalletService.GetWalletByPhoneNumber(input.ToWallet)

	if !w.walletExists(input.FromWallet) {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message.DEBIT_WALLET_DOES_NOT_EXIST.String()})
		return
	}

	if !w.walletExists(input.ToWallet) {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message.CREDIT_WALLET_DOES_NOT_EXIST.String()})
		return
	}

	//check if enabled
	if w.walletIsDisabled(input.FromWallet) {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusOK, "message": message.WALLET_IS_DISABLED.String()})
		return
	}

	tranxAmount, _ := strconv.Atoi(input.Amount)

	//check for negative amount
	if tranxAmount < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid amount"})
		return
	}

	tranxAmount *= 100
	//check if it debit wallet has sufficient fund
	if debitWallet.Balance < tranxAmount {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusOK, "message": message.INSUFFICIEND_FUND.String()})
		return
	}

	//debit-credit wallets
	err = w.WalletHandler.DebitCreditWallet(debitWallet, creditWallet, tranxAmount)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": message.TRANSACTION_DECLINED.String()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": message.TRANSACTION_APPROVED.String()})

}

func (w WalletFacade) EnableWallet(c *gin.Context) {

	var input walletmodel.WalletStatusInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	if !w.walletExists(input.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message.WALLET_NOT_EXIST.String()})
		return
	}

	if !w.walletIsDisabled(input.PhoneNumber) {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": message.WALLET_ALREADY_ENABLED.String()})
		return
	}

	err = w.WalletHandler.EnableWallet(input.PhoneNumber)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": message.WALLET_NOT_SUCCESSFULLY_ENABLED.String()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": message.WALLET_SUCCESSFULLY_ENABLED.String()})

}

func (w WalletFacade) DisableWallet(c *gin.Context) {

	var input walletmodel.WalletStatusInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	if !w.walletExists(input.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message.WALLET_NOT_EXIST.String()})
		return
	}

	if w.walletIsDisabled(input.PhoneNumber) {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": message.WALLET_ALREADY_DISABLED.String()})
		return
	}
	err = w.WalletHandler.DisableWallet(input.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": message.WALLET_NOT_SUCCESSFULLY_DISABLED.String()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": message.WALLET_SUCCESSFULLY_DISABLED.String()})

}

func (w WalletFacade) walletExists(phoneNumber string) bool {
	exists := true
	emptyWallet := walletmodel.Wallet{}

	wallet := w.WalletHandler.WalletService.GetWalletByPhoneNumber(phoneNumber)

	if wallet == emptyWallet {
		exists = false
	}

	return exists

}

func (w WalletFacade) walletIsDisabled(PhoneNumber string) bool {

	wallet := w.WalletHandler.WalletService.GetWalletByPhoneNumber(PhoneNumber)

	return wallet.Disabled
}
