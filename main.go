package main

import (
	"context"
	"fmt"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/morscino/wallet-engine/database/postgres"
	"github.com/morscino/wallet-engine/facade"
	"github.com/morscino/wallet-engine/handlers"
	"github.com/morscino/wallet-engine/models/walletmodel"
	"github.com/morscino/wallet-engine/routes"
	"github.com/morscino/wallet-engine/service/walletservice"
	"github.com/morscino/wallet-engine/utility/config"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env data: %v", err)
	}
	var config config.Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Could not load configuration data: %v", err)
	}

	db := postgres.DbConnect(config.DB)

	//migrations
	var WalletModel walletmodel.Wallet
	db.AutoMigrate(&WalletModel)

	server := gin.New()
	server.Use(gin.Logger())

	server.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			//c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	var ctx context.Context

	WalletRepo := walletservice.NewWalletervice(db)
	WalletHandler := handlers.NewWalletHandler(WalletRepo)
	WalletFacade := *facade.NewWalletFacade(WalletHandler, ctx)

	w := routes.NewWalletRoute(WalletFacade)
	w.WalletRoutes(server)

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if config.App.Port == "" {
		config.App.Port = "7000"
	}

	if err := server.Run(fmt.Sprintf(":%s", config.App.Port)); err != nil {

		//log.Error()
		log.Fatal("main run:", err)
	}

}
