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

	// TODO: validate opts

	rate, err := s.apis.Crypto.GetRate(ctx, opts)
	if err != nil {
		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate from api: %w", err)
	}

	logger.Info("successfully got rate")
	return rate, nil
}
