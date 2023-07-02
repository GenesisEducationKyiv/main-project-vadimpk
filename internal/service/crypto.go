package service

import (
	"context"
	"fmt"
)

type cryptoService struct {
	serviceContext
}

func NewCryptoService(opts *Options) *cryptoService {
	return &cryptoService{
		serviceContext: serviceContext{
			storages: opts.Storages,
			apis:     opts.APIs,
			logger:   opts.Logger.Named("CryptoService"),
			cfg:      opts.Cfg,
		},
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

	chain, err := NewCryptoAPIChain(s.apis.CryptoProviders)
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
