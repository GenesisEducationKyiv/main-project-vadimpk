package coinapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/vadimpk/gses-2023/internal/service"
)

type getRateResponseBody struct {
	Rate float64 `json:"rate"`
}

func (c *coinAPI) GetRate(ctx context.Context, opts *service.GetRateOptions) (float64, error) {
	logger := c.logger.
		Named("GetRate").
		WithContext(ctx).
		With("opts", opts)

	url := fmt.Sprintf("/exchangerate/%s/%s", strings.ToUpper(opts.CryptoCurrency), strings.ToUpper(opts.Currency))

	var respBody getRateResponseBody
	resp, err := c.client.R().
		SetResult(&respBody).
		Get(url)
	if err != nil {
		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error("failed to get rate", "status", resp.Status())
		return 0, fmt.Errorf("failed to get rate: status %s", resp.Status())
	}
	logger = logger.With("rate", respBody.Rate)

	logger.Info("successfully got rate")
	return respBody.Rate, nil
}
