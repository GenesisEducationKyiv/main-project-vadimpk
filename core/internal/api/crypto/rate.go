package crypto

import (
	"context"
	"fmt"
	"net/http"
)

func (c *cryptoAPI) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	logger := c.logger.Named("GetRate").
		WithContext(ctx).
		With("fromCurrency", fromCurrency).
		With("toCurrency", toCurrency)

	var respBody float64
	res, err := c.client.R().
		SetQueryParams(map[string]string{
			"crypto_currency": fromCurrency,
			"fiat_currency":   toCurrency,
		}).
		SetResult(&respBody).
		Get("/api/rate")
	logger = logger.With("response", res.String()).
		With("status", res.Status()).
		With("respBody", respBody)

	if err != nil {
		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}

	return respBody, err
}
