package jupiter

import (
	"github.com/ilkamo/jupiter-go/jupiter"
	"log"
	"os"
	"solTrader/pkg/config"
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

// BuyToken buys the specified amount of a token, buying is with token set in config, default is USDC
func (c *Client) BuyToken(tokenMint string, amount float64) error {
	quoteResponse, err := c.JupiterClient.GetQuoteWithResponse(c.config.Ctx, &jupiter.GetQuoteParams{
		InputMint:  tokenMint,
		OutputMint: c.config.BaseMint,
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
	return nil
}

// SellToken sells the specified amount of a token, selling receives token set in config, default is USDC
func (c *Client) SellToken(tokenMint string, amount float64) error {
	// TODO: Implement the selling logic using the Jupiter client
	return nil
}

// GetCurrentPrice gets the current market price of a token, to USDC :)
func (c *Client) GetCurrentPrice(tokenMint string) (float64, error) {
	// TODO: Implement the logic to fetch current price using the Jupiter client

	return 0, nil
}
