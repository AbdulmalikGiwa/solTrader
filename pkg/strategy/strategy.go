package strategy

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"log"
	"solTrader/internal/db"
)

// TradeStrategy is the interface for the trade strategy, we currently have just one simple strategy but can be
// extended to more complex strategies.
type TradeStrategy interface {
	ShouldBuy(DB *sql.DB, currentPrice decimal.Decimal, token string) (bool, float64)
	ShouldSell(DB *sql.DB, currentPrice decimal.Decimal, token string) (bool, float64)
}

type Strategy struct {
	BuyThreshold  float64 // e.g., 0.10 for 10%
	SellThreshold float64 // e.g., 0.20 for 20%
	Log           log.Logger
}

// ShouldBuy is simply custom strategy to determine if we should buy a token, can be customized, very simple for now
func (s *Strategy) ShouldBuy(DB *sql.DB, currentPrice decimal.Decimal, token string) (bool, float64) {
	balance, lastSellPrice, err := db.GetBalanceAndLastPrice(DB, token)
	if err != nil {
		s.Log.Printf("ShouldBuy: failed to get last price for token from db: %s, error: %s", token, err)
		return false, 0
	}
	s.Log.Printf("last sell price: %s", lastSellPrice)
	return currentPrice.LessThanOrEqual(lastSellPrice.Mul(decimal.NewFromFloat(1 - s.BuyThreshold))), balance
}

// ShouldSell is simply custom strategy to determine if we should sell a token, can be customized, very simple for now
func (s *Strategy) ShouldSell(DB *sql.DB, currentPrice decimal.Decimal, token string) (bool, float64) {
	balance, lastBuyPrice, err := db.GetBalanceAndLastPrice(DB, token)
	if err != nil {
		s.Log.Printf("ShouldSell: failed to get last price for token from db: %s, error: %s", token, err)
		return false, 0
	}
	s.Log.Printf("last buy price: %s", lastBuyPrice)
	return currentPrice.GreaterThanOrEqual(lastBuyPrice.Mul(decimal.NewFromFloat(1 + s.SellThreshold))), balance

}
