package mailgun

import (
	"context"
	"fmt"
	"github.com/vadimpk/gses-2023/internal/service"
)

func (m *mailgunAPI) Send(ctx context.Context, opts *service.SendOptions) error {
	logger := m.logger.Named("Send").
		WithContext(ctx).
		With("opts", opts)

	_, id, err := m.client.Send(ctx, m.client.NewMessage(m.from, opts.Subject, opts.Body, opts.To))
	if err != nil {
		logger.Error("failed to send email", "err", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	logger = logger.With("id", id)

	logger.Info("successfully sent email")
	return err
}
