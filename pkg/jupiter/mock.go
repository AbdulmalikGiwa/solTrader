package jupiter

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

// MockClient simulates the behavior of the Jupiter client for testing
type MockClient struct{}

// NewMockClient creates a new instance of MockClient
func NewMockClient() *MockClient {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	return &MockClient{}
}

// TradeToken simulates trading a specified amount of a token
func (m *MockClient) TradeToken(tokenMint string, amount float64, tradeType string) error {
	if int(amount)%2 != 0 {
		return errors.New("amount must be an even number")
	}
	fmt.Printf("Trade executed: %s %f of token %s\n", tradeType, amount, tokenMint)
	return nil
}

// GetCurrentPrice simulates fetching the current market price of a token
func (m *MockClient) GetCurrentPrice(tokenMint string) (decimal.Decimal, error) {
	price := decimal.NewFromFloat(float64(rand.Intn(5) + 1)) // Generate a random price between 1 and 5
	return price, nil
}
