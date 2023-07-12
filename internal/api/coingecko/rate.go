package coingecko

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type getRateResponseBody map[string]map[string]float64

func (c *coinGeckoAPI) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	logger := c.logger.
		Named("GetRate").
		WithContext(ctx).
		With("fromCurrency", fromCurrency).
		With("toCurrency", toCurrency)

	fromCurrency = parseCurrency(fromCurrency)
	toCurrency = parseCurrency(toCurrency)

	var respBody getRateResponseBody
	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"ids":           strings.ToUpper(fromCurrency),
			"vs_currencies": strings.ToUpper(toCurrency),
		}).
		SetResult(&respBody).
		Get("/simple/price")
	logger = logger.With("responseBody", resp.String()).With("statusCode", resp.StatusCode())

	if err != nil {
		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Error("failed to get rate")
		return 0, fmt.Errorf("failed to get rate: status %s", resp.Status())
	}
	logger = logger.With("response", respBody)

	if _, ok := respBody[fromCurrency][toCurrency]; !ok {
		logger.Error("failed to get rate", "err", "no such currency")
		return 0, fmt.Errorf("failed to get rate: no such currency")
	}

	logger.Info("successfully got rate")
	return respBody[fromCurrency][toCurrency], nil
}

func parseCurrency(c string) string {
	switch c {
	case "BTC":
		return "bitcoin"
	case "ETH":
		return "ethereum"
	case "USD":
		return "usd"
	case "UAH":
		return "uah"
	}
	return ""
}
