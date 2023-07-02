package coinbase

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type getRateResponseBody struct {
	Data struct {
		Currency string             `json:"currency"`
		Rates    map[string]float64 `json:"rates"`
	} `json:"data"`
}

func (c *coinbaseAPI) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	logger := c.logger.
		Named("GetRate").
		WithContext(ctx).
		With("fromCurrency", fromCurrency).
		With("toCurrency", toCurrency)

	var respBody getRateResponseBody
	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"currency": strings.ToUpper(fromCurrency),
		}).
		SetResult(&respBody).
		Get("/exchange-rates")

	if err != nil {
		logger.Error("failed to get rate", "err", err, "body", resp.String())
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error("failed to get rate", "status", resp.Status(), "body", resp.String())
		return 0, fmt.Errorf("failed to get rate: status %s", resp.Status())
	}
	logger = logger.With("response", respBody)

	rate, ok := respBody.Data.Rates[strings.ToUpper(toCurrency)]
	if !ok {
		logger.Error("Currency not found in response", "currency", toCurrency)
		return 0, fmt.Errorf("currency %s not found in response", toCurrency)
	}

	logger.Info("successfully got rate")
	return rate, nil
}
