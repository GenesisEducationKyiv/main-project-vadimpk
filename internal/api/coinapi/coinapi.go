package coinapi

import (
	"github.com/go-resty/resty/v2"

	"github.com/vadimpk/gses-2023/pkg/logging"
)

type coinAPI struct {
	client *resty.Client
	logger logging.Logger
}

type Options struct {
	Logger logging.Logger
	APIKey string
}

func New(opts *Options) *coinAPI {
	c := resty.New()

	c = c.SetBaseURL("https://rest.coinapi.io/v1").
		SetHeader("X-CoinAPI-Key", opts.APIKey)

	return &coinAPI{
		client: c,
		logger: opts.Logger.Named("CoinAPI"),
	}
}
