package coinbase

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type getRateResponseBody struct {
	Data struct {
		Currency string            `json:"currency"`
		Rates    map[string]string `json:"rates"`
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
	logger = logger.With("responseBody", resp.String()).With("statusCode", resp.StatusCode())

	if err != nil {
		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error("failed to get rate")
		return 0, fmt.Errorf("failed to get rate: status %s", resp.Status())
	}
	logger = logger.With("respBody", respBody)

	rate, ok := respBody.Data.Rates[strings.ToUpper(toCurrency)]
	if !ok {
		logger.Error("toCurrency not found in response")
		return 0, fmt.Errorf("currency %s not found in response", toCurrency)
	}

	logger.Info("successfully got rate")
	return strconv.ParseFloat(rate, 64)
}
