package crypto

import (
	"github.com/go-resty/resty/v2"
	"github.com/vadimpk/gses-2023/core/config"
	"github.com/vadimpk/gses-2023/core/pkg/logging"
)

type cryptoAPI struct {
	client *resty.Client
	logger logging.Logger
}

type Options struct {
	Logger logging.Logger
	Config *config.Config
}

func New(options *Options) *cryptoAPI {
	h := resty.New().
		SetBaseURL(options.Config.CryptoService.BaseURL)

	return &cryptoAPI{
		client: h,
		logger: options.Logger.Named("MailgunAPI"),
	}
}
