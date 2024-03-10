package config

import (
	"context"
	"fmt"
	"os"
)

type TradingConfig struct {
	ClientRPC string
	BaseMint  string // defaults to USDC, every buy and sell is in USDC. Can be changed via config
	Ctx       context.Context
	PubKey    string
}

func LoadConfig() (*TradingConfig, error) {
	rpcURL := os.Getenv("RPC_URL")
	mintAddress := os.Getenv("BASE_MINT")
	publicKey := os.Getenv("PUBLIC_KEY")
	ctx := context.Background()
	if mintAddress == "" {
		mintAddress = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	}
	if rpcURL == "" {
		return nil, fmt.Errorf("RPC_URL environment variable is not set")
	}
	if publicKey == "" {
		return nil, fmt.Errorf("PUBLIC_KEY environment variable is not set")
	}
	return &TradingConfig{
		ClientRPC: rpcURL,
		BaseMint:  mintAddress,
		Ctx:       ctx,
		PubKey:    publicKey,
	}, nil
}
