package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/morscino/wallet-engine/facade"
)

type WalletRoute struct {
	WalletFacade facade.WalletFacade
}

func NewWalletRoute(wf facade.WalletFacade) *WalletRoute {
	return &WalletRoute{WalletFacade: wf}
}

func (c WalletRoute) WalletRoutes(router *gin.Engine) {

	WalletGroup := router.Group("/wallet")
	{
		WalletGroup.GET("/", c.WalletFacade.Test)
		WalletGroup.POST("/create-wallet", c.WalletFacade.CreateWallet)
		WalletGroup.POST("/debit-credit", c.WalletFacade.DebitCreditWallet)
		WalletGroup.PUT("/enable-wallet", c.WalletFacade.EnableWallet)
		WalletGroup.PUT("/disable-wallet", c.WalletFacade.DisableWallet)

	}

}
