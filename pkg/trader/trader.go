package trader

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"solTrader/pkg/jupiter"
	"solTrader/pkg/strategy"
)

// Trader encapsulates trading logic
type Trader struct {
	JupiterClient *jupiter.Client
	Strategy      strategy.TradeStrategy
	LastBuyPrice  decimal.Decimal
	LastSellPrice decimal.Decimal
	Token         string
}

// NewTrader creates a new trader instance
func NewTrader(client *jupiter.Client, strategy *strategy.Strategy, token string) *Trader {
	return &Trader{
		JupiterClient: client,
		Strategy:      strategy,
		Token:         token,
	}
}

// ExecuteTrade executes the trading strategy
func (t *Trader) ExecuteTrade(DB *sql.DB) error {
	currentPrice, err := t.JupiterClient.GetCurrentPrice(t.Token)
	if err != nil {
		return err
	}
	shouldBuy, buyBalance := t.Strategy.ShouldBuy(DB, currentPrice, t.Token)
	shouldSell, sellBalance := t.Strategy.ShouldSell(DB, currentPrice, t.Token)
	if shouldBuy {
		// Execute buy logic
		if err := t.JupiterClient.TradeToken(t.Token, buyBalance, jupiter.BuyType); err != nil {
			return err
		}
		t.LastBuyPrice = currentPrice
	} else if shouldSell {
		// Execute sell logic
		if err := t.JupiterClient.TradeToken(t.Token, sellBalance, jupiter.SellType); err != nil {
			return err
		}
		t.LastSellPrice = currentPrice
	}

	return nil
}
