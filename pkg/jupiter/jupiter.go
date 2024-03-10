package jupiter

import (
	"errors"
	"github.com/ilkamo/jupiter-go/jupiter"
	"log"
	"os"
	"solTrader/pkg/config"
	"time"
)

const (
	// BuyType is the type of trade to buy
	BuyType = "buy"
	// SellType is the type of trade to sell
	SellType = "sell"
)

// Client wraps the Jupiter client
type Client struct {
	JupiterClient *jupiter.ClientWithResponses
	log           *log.Logger
	config        *config.TradingConfig
}

// NewClient creates a new Jupiter client
func NewClient(config *config.TradingConfig) *Client {
	jupClient, err := jupiter.NewClientWithResponses(config.ClientRPC)
	if err != nil {
		log.Fatal("failed to create Jupiter client", err)
	}
	logger := log.New(os.Stdout, "jupiter", log.LstdFlags)
	return &Client{
		JupiterClient: jupClient,
		log:           logger,
		config:        config,
	}
}

// TradeToken buys the specified amount of a token, buying is with token set in config, default is USDC
func (c *Client) TradeToken(tokenMint string, amount float64, tradeType string) error {
	var buyToken string
	var sellToken string

	if tradeType == BuyType {
		buyToken = tokenMint
		sellToken = c.config.BaseMint
	}
	if tradeType == SellType {
		buyToken = c.config.BaseMint
		sellToken = tokenMint
	}
	quoteResponse, err := c.JupiterClient.GetQuoteWithResponse(c.config.Ctx, &jupiter.GetQuoteParams{
		InputMint:  buyToken,
		OutputMint: sellToken,
		Amount:     jupiter.AmountParameter(amount),
		// TODO: Set slippage
	})
	if err != nil {
		c.log.Println("BuyToken: failed to get quote", err)
		return err
	}
	if quoteResponse.JSON200 == nil {
		c.log.Println("BuyToken: failed to get quote, returned NON 200 code: ", quoteResponse.StatusCode())
		return err
	}
	dynamicComputeUnitLimit := true
	quote := quoteResponse.JSON200
	swapResponse, err := c.JupiterClient.PostSwapWithResponse(c.config.Ctx, jupiter.PostSwapJSONRequestBody{
		QuoteResponse:           *quote,
		UserPublicKey:           c.config.PubKey,
		DynamicComputeUnitLimit: &dynamicComputeUnitLimit,
	})
	if err != nil {
		c.log.Println("BuyToken: failed to swap", err)
		return err
	}
	if swapResponse.JSON200 == nil {
		c.log.Println("BuyToken: failed to swap, returned NON 200 code: ", swapResponse.StatusCode())
		return err
	}
	swap := swapResponse.JSON200
	tx, err := c.config.SolClient.SendTransactionOnChain(c.config.Ctx, swap.SwapTransaction)

	// Wait for tx, could be done better tbh. Just lifting sample from jup client library
	time.Sleep(20 * time.Second)
	confirmed, err := c.config.SolClient.CheckSignature(c.config.Ctx, tx)
	if err != nil {
		c.log.Println("BuyToken: failed to check signature", err)
		// will raise this error a level up to have a retry mechanism of some sort
	}
	if confirmed {
		c.log.Println("BuyToken: transaction confirmed")
	} else {
		c.log.Println("BuyToken: transaction not confirmed")
		err = errors.New("transaction not confirmed")
		return err
		// can be refactored to keep polling and identify failed txs
	}
	return nil
}

// GetCurrentPrice gets the current market price of a token, to USDC :)
func (c *Client) GetCurrentPrice(tokenMint string) (float64, error) {
	// TODO: Implement the logic to fetch current price using the Jupiter client

	return 0, nil
}
