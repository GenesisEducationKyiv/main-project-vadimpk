package coingecko

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type getRateResponseBody struct {
	Rate float64 `json:"rate"`
}

func (c *coinGeckoAPI) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	logger := c.logger.
		Named("GetRate").
		WithContext(ctx).
		With("fromCurrency", fromCurrency).
		With("toCurrency", toCurrency)

	var respBody getRateResponseBody
	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"ids":           strings.ToUpper(fromCurrency),
			"vs_currencies": strings.ToUpper(toCurrency),
		}).
		SetResult(&respBody).
		Get("/simple/price")

	if err != nil {
		logger.Error("failed to get rate", "err", err, "body", resp.String())
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error("failed to get rate", "status", resp.Status(), "body", resp.String())
		return 0, fmt.Errorf("failed to get rate: status %s", resp.Status())
	}
	logger = logger.With("rate", respBody.Rate)

	logger.Info("successfully got rate")
	return respBody.Rate, nil
}
