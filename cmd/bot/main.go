package main

import "solTrader/internal/db"

func main() {
	database := db.OpenDB("trades.db")
	defer database.Close()
	db.CreateTable(database)
}
