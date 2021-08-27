package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/morscino/wallet-engine/database/postgres"
	"github.com/morscino/wallet-engine/facade"
	"github.com/morscino/wallet-engine/handlers"
	"github.com/morscino/wallet-engine/models/transactionmodel"
	"github.com/morscino/wallet-engine/models/walletmodel"
	"github.com/morscino/wallet-engine/routes"
	"github.com/morscino/wallet-engine/service/transactionservice"
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

	postgresDB := postgres.DbConnect(config.PostgresDB)

	//migrations
	var WalletModel walletmodel.Wallet
	var TransactionModel transactionmodel.Transaction

	postgresDB.AutoMigrate(&WalletModel)
	postgresDB.AutoMigrate(&TransactionModel)

	server := gin.New()
	server.Use(gin.Logger())

	server.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {

			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	TransactionRepo := transactionservice.NewTransactionService(postgresDB)
	TransactionHandler := handlers.NewTransactionHandler(TransactionRepo)

	WalletRepo := walletservice.NewWalletervice(postgresDB)
	WalletHandler := handlers.NewWalletHandler(WalletRepo, TransactionHandler)
	WalletFacade := *facade.NewWalletFacade(WalletHandler, TransactionHandler)

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
