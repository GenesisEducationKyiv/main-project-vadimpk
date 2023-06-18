package mailgun

import (
	"github.com/mailgun/mailgun-go/v4"

	"github.com/vadimpk/gses-2023/pkg/logging"
)

type mailgunAPI struct {
	client *mailgun.MailgunImpl
	from   string
	logger logging.Logger
}

type Options struct {
	APIKey string
	Domain string
	From   string

	Logger logging.Logger
}

func New(options *Options) *mailgunAPI {
	return &mailgunAPI{
		client: mailgun.NewMailgun(options.Domain, options.APIKey),
		from:   options.From,
		logger: options.Logger.Named("MailgunAPI"),
	}
}
