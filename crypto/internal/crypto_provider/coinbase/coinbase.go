package coinbase

import (
	"github.com/go-resty/resty/v2"

	"github.com/vadimpk/gses-2023/crypto/pkg/logging"
)

type coinbaseAPI struct {
	client *resty.Client
	logger logging.Logger
}

type Options struct {
	Logger logging.Logger
}

func New(opts *Options) *coinbaseAPI {
	c := resty.New()

	c = c.SetBaseURL("https://api.coinbase.com/v2")

	return &coinbaseAPI{
		client: c,
		logger: opts.Logger.Named("CoinBaseAPI"),
	}
}
