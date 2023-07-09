package crypto

import (
	"context"
	"errors"
	"fmt"

	"github.com/vadimpk/gses-2023/config"
	"github.com/vadimpk/gses-2023/crypto/internal/crypto_provider"
	"github.com/vadimpk/gses-2023/crypto/internal/entity"
	"github.com/vadimpk/gses-2023/crypto/pkg/logging"
)

type Service interface {
	// GetRate returns current rate for crypto currency.
	GetRate(ctx context.Context, opts *GetRateOptions) (float64, error)
}

var (
	ErrGetRateInvalidCryptoCurrency = errors.New("invalid crypto currency")
	ErrGetRateInvalidFiatCurrency   = errors.New("invalid fiat currency")
)

type GetRateOptions struct {
	Crypto entity.CryptoCurrency
	Fiat   entity.FiatCurrency
}

func (o *GetRateOptions) Validate() error {
	if !o.Crypto.IsValid() {
		return ErrGetRateInvalidCryptoCurrency
	}

	if !o.Fiat.IsValid() {
		return ErrGetRateInvalidFiatCurrency
	}

	return nil
}

type Options struct {
	Providers crypto_provider.CryptoAPIProviders
	Logger    logging.Logger
	Config    *config.Config
}

type cryptoService struct {
	logger    logging.Logger
	cfg       *config.Config
	providers crypto_provider.CryptoAPIProviders
}

func NewCryptoService(opts Options) *cryptoService {
	return &cryptoService{
		providers: opts.Providers,
		logger:    opts.Logger.Named("Crypto"),
		cfg:       opts.Config,
	}
}

func (s *cryptoService) GetRate(ctx context.Context, opts *GetRateOptions) (float64, error) {
	logger := s.logger.Named("GetRate").
		WithContext(ctx).
		With("opts", opts)

	if err := opts.Validate(); err != nil {
		logger.Info(err.Error())
		return 0, err
	}

	chain, err := crypto_provider.NewCryptoAPIChain(s.providers)
	if err != nil {
		logger.Error("failed to create crypto api chain", "err", err)
		return 0, fmt.Errorf("failed to create crypto api chain: %w", err)
	}

	rate, err := chain.GetRate(ctx, opts.Crypto.String(), opts.Fiat.String())
	if err != nil {
		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate from api: %w", err)
	}

	logger.Info("successfully got rate")
	return rate, nil
}
