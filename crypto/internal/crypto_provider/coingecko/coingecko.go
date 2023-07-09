package coingecko

import (
	"github.com/go-resty/resty/v2"

	"github.com/vadimpk/gses-2023/crypto/pkg/logging"
)

type coinGeckoAPI struct {
	client *resty.Client
	logger logging.Logger
}

type Options struct {
	Logger logging.Logger
}

func New(opts *Options) *coinGeckoAPI {
	c := resty.New()

	c = c.SetBaseURL("https://api.coingecko.com/api/v3")

	return &coinGeckoAPI{
		client: c,
		logger: opts.Logger.Named("CoinGeckoAPI"),
	}
}
