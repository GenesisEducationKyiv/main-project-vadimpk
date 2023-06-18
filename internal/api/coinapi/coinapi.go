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
	ApiKey string
	Logger logging.Logger
}

func New(opts *Options) *coinAPI {
	c := resty.New()

	c = c.SetBaseURL("https://rest.coinapi.io/v1").
		SetHeader("X-CoinAPI-Key", opts.ApiKey)

	return &coinAPI{
		client: c,
		logger: opts.Logger.Named("CoinAPI"),
	}
}
