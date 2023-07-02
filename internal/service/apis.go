package service

import "context"

type APIs struct {
	Email  EmailAPI
	Crypto CryptoAPI
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
