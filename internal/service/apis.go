package service

import (
	"context"
	"fmt"
)

type APIs struct {
	Email           EmailAPI
	CryptoProviders CryptoAPIProviders
}

// EmailAPI provides methods for sending emails that are used in EmailService and
// implemented in external packages.
//
//go:generate go run github.com/vektra/mockery/v2@v2.27.1 --dir . --name EmailAPI --output ../../internal/service/mocks
type EmailAPI interface {
	Send(ctx context.Context, opts *SendOptions) error
}

type SendOptions struct {
	To      string
	Subject string
	Body    string
}

// CryptoAPI provides methods for getting crypto rates that are used in CryptoService and
// implemented in external packages.
type CryptoAPI interface {
	GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error)
}

// CryptoAPIProviderType represents type of 3rd party provider of crypto API.
type CryptoAPIProviderType string

const (
	CryptoAPIProviderCoinbase  CryptoAPIProviderType = "coinbase"
	CryptoAPIProviderCoinAPI   CryptoAPIProviderType = "coinapi"
	CryptoAPIProviderCoinGecko CryptoAPIProviderType = "coingecko"
)

// CryptoAPIProviders is used to store all available crypto APIs.
type CryptoAPIProviders map[CryptoAPIProviderType]CryptoAPI

// CryptoAPIChain is used to create chain of responsibility for crypto APIs.
type CryptoAPIChain struct {
	api  CryptoAPI
	next *CryptoAPIChain
}

// NewCryptoAPIChain creates chain of responsibility for crypto APIs.
// It uses order to determine which API to use first. If no order is provided, default order is used.
func NewCryptoAPIChain(providers CryptoAPIProviders, order ...CryptoAPIProviderType) (*CryptoAPIChain, error) {
	var lastProvider *CryptoAPIChain

	// default order
	if len(order) == 0 {
		order = []CryptoAPIProviderType{
			CryptoAPIProviderCoinAPI,
			CryptoAPIProviderCoinGecko,
			CryptoAPIProviderCoinbase,
		}
	}

	for i := len(order) - 1; i >= 0; i-- {
		provider := order[i]
		api, ok := providers[provider]
		if !ok {
			return nil, fmt.Errorf("unknown provider: %s", provider)
		}
		lastProvider = &CryptoAPIChain{
			api:  api,
			next: lastProvider,
		}
	}

	if lastProvider == nil {
		return nil, fmt.Errorf("no providers")
	}

	return lastProvider, nil
}

func (chain *CryptoAPIChain) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	rate, err := chain.api.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil && chain.next != nil {
		return chain.next.GetRate(ctx, fromCurrency, toCurrency)
	}
	return rate, err
}
