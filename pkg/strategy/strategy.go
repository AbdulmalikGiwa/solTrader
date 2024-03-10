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
	ShouldBuy(DB *sql.DB, currentPrice decimal.Decimal, token string) bool
	ShouldSell(currentPrice, lastBuyPrice float64) bool
}

type Strategy struct {
	BuyThreshold  float64 // e.g., 0.10 for 10%
	SellThreshold float64 // e.g., 0.20 for 20%
	Log           log.Logger
}

// ShouldBuy is simply custom strategy to determine if we should buy a token, can be customized, very simple for now
func (s *Strategy) ShouldBuy(DB *sql.DB, currentPrice decimal.Decimal, token string) bool {
	lastSellPrice, err := db.GetLastPrice(DB, token)
	if err != nil {
		s.Log.Printf("failed to get last price for token from db: %s, error: %s", token, err)
	}
	return currentPrice.LessThanOrEqual(lastSellPrice.Mul(decimal.NewFromFloat(1 - s.BuyThreshold)))
}

// ShouldSell is simply custom strategy to determine if we should sell a token, can be customized, very simple for now
func (s *Strategy) ShouldSell(DB *sql.DB, currentPrice decimal.Decimal, token string) bool {
	lastBuyPrice, err := db.GetLastPrice(DB, token)
	if err != nil {
		s.Log.Printf("failed to get last price for token from db: %s, error: %s", token, err)
	}
	return currentPrice.GreaterThanOrEqual(lastBuyPrice.Mul(decimal.NewFromFloat(1 + s.SellThreshold)))

}
