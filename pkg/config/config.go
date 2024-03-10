package config

import (
	"context"
	"fmt"
	"github.com/ilkamo/jupiter-go/solana"
	"os"
)

type TradingConfig struct {
	ClientRPC string
	BaseMint  string // defaults to USDC, every buy and sell is in USDC. Can be changed via config
	Ctx       context.Context
	PubKey    string
	SolClient *solana.Client
}

func LoadConfig() (*TradingConfig, error) {
	ctx := context.Background()
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("RPC_URL environment variable is not set")
	}
	mintAddress := os.Getenv("BASE_MINT")
	if mintAddress == "" {
		mintAddress = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	}
	publicKey := os.Getenv("PUBLIC_KEY")
	if publicKey == "" {
		return nil, fmt.Errorf("PUBLIC_KEY environment variable is not set")
	}
	privateKey := os.Getenv("PRIVATE_KEY")
	// We wont pass privateKey to the config, but we will check if it is set
	if privateKey == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable is not set")
	}

	wallet, err := solana.NewWalletFromPrivateKeyBase58(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet from private key: %w", err)
	}

	client, err := solana.NewClient(wallet, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create solana client: %w", err)
	}

	return &TradingConfig{
		ClientRPC: rpcURL,
		BaseMint:  mintAddress,
		Ctx:       ctx,
		PubKey:    publicKey,
		SolClient: &client,
	}, nil
}
