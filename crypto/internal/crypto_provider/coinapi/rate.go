package coinapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type getRateResponseBody struct {
	Rate float64 `json:"rate"`
}

func (c *coinAPI) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	logger := c.logger.
		Named("GetRate").
		WithContext(ctx).
		With("fromCurrency", fromCurrency).
		With("toCurrency", toCurrency)

	url := fmt.Sprintf("/exchangerate/%s/%s", strings.ToUpper(fromCurrency), strings.ToUpper(toCurrency))

	var respBody getRateResponseBody
	resp, err := c.client.R().
		SetResult(&respBody).
		Get(url)
	logger = logger.With("responseBody", resp.String()).With("statusCode", resp.StatusCode())

	if err != nil {
		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error("failed to get rate")
		return 0, fmt.Errorf("failed to get rate: status %s", resp.Status())
	}
	logger = logger.With("rate", respBody.Rate)

	logger.Info("successfully got rate")
	return respBody.Rate, nil
}
