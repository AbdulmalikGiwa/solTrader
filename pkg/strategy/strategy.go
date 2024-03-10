package strategy

type Strategy struct {
	BuyThreshold  float64 // e.g., 0.10 for 10%
	SellThreshold float64 // e.g., 0.20 for 20%
}

func (s *Strategy) ShouldBuy(currentPrice, lastSellPrice float64) bool {
	return currentPrice <= lastSellPrice*(1-s.BuyThreshold)
}

func (s *Strategy) ShouldSell(currentPrice, lastBuyPrice float64) bool {
	return currentPrice >= lastBuyPrice*(1+s.SellThreshold)
}
