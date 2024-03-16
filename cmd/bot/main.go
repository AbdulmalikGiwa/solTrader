package main

import (
	"log"
	"solTrader/internal/db"
	"solTrader/pkg/config"
	"solTrader/pkg/jupiter"
	"solTrader/pkg/strategy"
	"solTrader/pkg/trader"
	"time"
)

func main() {
	database := db.OpenDB("trades.db")
	defer database.Close()
	db.CreateTable(database)
	tradingConfig, err := config.LoadConfig()
	if err != nil {
		log.Println("failed to load config", err)
		panic(err)
	}
	log.Println("loaded config", tradingConfig)
	jupiterClient, err := jupiter.NewClient(tradingConfig)
	if err != nil {
		log.Println("failed to create Jupiter client", err)
		panic(err)
	}
	tradingStrategy := &strategy.Strategy{
		BuyThreshold:  0.10,
		SellThreshold: 0.10,
		Log:           *log.Default(),
	}
	// TODO: Replace the hardcoded token with the token from the config
	mainTrader := trader.NewTrader(jupiterClient, tradingStrategy, "BFek4xVLbyW9w2cfcuFxh974f7TtAjWWjJq2kSrgthGL")
	for {
		if err := mainTrader.ExecuteTrade(database); err != nil {
			log.Printf("Error executing trade: %v", err)
		}
		time.Sleep(1 * time.Minute)
	}
}
