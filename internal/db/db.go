package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shopspring/decimal"
	"log"
)

func OpenDB(databasePath string) *sql.DB {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Println("failed to open database", err)
		log.Fatal(err)
	}
	return db
}

func CreateTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS trades (
		"token_address" TEXT NOT NULL PRIMARY KEY,
		"last_price" REAL,
		"holding" BOOLEAN,
		"balance" REAL
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}
}

func UpsertTokenTrade(db *sql.DB, address string, lastPrice decimal.Decimal, holding bool) {
	lastPriceString := lastPrice.String()
	upsertSQL := `INSERT INTO trades (token_address, last_price, holding) VALUES (?, ?, ?)
				ON CONFLICT(token_address) DO UPDATE SET last_price = excluded.last_price, holding = excluded.holding;`
	holdingInt := 0
	if holding {
		holdingInt = 1
	}
	statement, err := db.Prepare(upsertSQL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(address, lastPriceString, holdingInt)
	if err != nil {
		log.Fatal(err)
	}
}

// GetLastPrice returns the last price of the token for the given address
func GetLastPrice(db *sql.DB, address string) (decimal.Decimal, error) {
	var lastPrice float64

	query := `SELECT last_price FROM trades WHERE token_address = ?`
	err := db.QueryRow(query, address).Scan(&lastPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			// No result for the given address
			return decimal.Zero, fmt.Errorf("no token found for address: %s", address)
		}
		// Some other error occurred
		return decimal.Zero, err
	}

	// Convert the float64 to decimal.Decimal
	lastPriceDec := decimal.NewFromFloat(lastPrice)
	return lastPriceDec, nil
}

func GetBalance(db *sql.DB, address string) (decimal.Decimal, error) {
	var balance float64

	query := `SELECT balance FROM trades WHERE token_address = ?`
	err := db.QueryRow(query, address).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			// No result for the given address
			return decimal.Zero, fmt.Errorf("no token found for address: %s", address)
		}
		// Some other error occurred
		return decimal.Zero, err
	}

	// Convert the float64 to decimal.Decimal
	balanceDec := decimal.NewFromFloat(balance)
	return balanceDec, nil
}

func GetBalanceAndLastPrice(db *sql.DB, address string) (float64, decimal.Decimal, error) {
	var balance float64
	var lastPrice float64

	query := `SELECT balance, last_price FROM trades WHERE token_address = ?`
	err := db.QueryRow(query, address).Scan(&balance, &lastPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, decimal.Zero, fmt.Errorf("no token found for address: %s", address)
		}
		return 0, decimal.Zero, err
	}

	lastPriceDec := decimal.NewFromFloat(lastPrice)
	return balance, lastPriceDec, nil
}
