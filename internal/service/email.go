package service

import (
	"context"
	"fmt"
)

type emailService struct {
	serviceContext
	cryptoService CryptoService
}

func NewEmailService(opts *Options, cryptoService CryptoService) *emailService {
	return &emailService{
		serviceContext: serviceContext{
			storages: opts.Storages,
			apis:     opts.APIs,
			logger:   opts.Logger.Named("EmailService"),
			cfg:      opts.Cfg,
		},
		cryptoService: cryptoService,
	}
}

func (s *emailService) Subscribe(ctx context.Context, email string) error {
	logger := s.logger.Named("Subscribe").
		WithContext(ctx).
		With("email", email)

	existingEmail, err := s.storages.Email.Get(ctx, email)
	if err != nil {
		logger.Error("failed to get email from storage", "err", err)
		return fmt.Errorf("failed to get email from storage: %w", err)
	}
	if existingEmail != "" {
		logger.Info("email already exists")
		return ErrSubscribeAlreadySubscribed
	}

	err = s.storages.Email.Save(ctx, email)
	if err != nil {
		logger.Error("failed to save email", "err", err)
		return fmt.Errorf("failed to save email to storage: %w", err)
	}

	logger.Info("successfully subscribed")
	return nil
}

func (s *emailService) SendRateInfo(ctx context.Context) (*SendRateInfoOutput, error) {
	logger := s.logger.Named("SendRateInfo").
		WithContext(ctx)

	emails, err := s.storages.Email.List(ctx)
	if err != nil {
		logger.Error("failed to get emails from storage", "err", err)
		return nil, fmt.Errorf("failed to get emails from storage: %w", err)
	}

	rate, err := s.cryptoService.GetRate(ctx, &GetRateOptions{
		CryptoCurrency: "BTC",
		Currency:       "UAH",
	})
	if err != nil {
		logger.Error("failed to get rate", "err", err)
		return nil, fmt.Errorf("failed to get rate: %w", err)
	}

	var failedEmails []string
	for _, email := range emails {
		err = s.apis.Email.Send(ctx, &SendOptions{
			To:      email,
			Subject: "Rate info",
			Body:    fmt.Sprintf("Current rate is %f", rate),
		})
		if err != nil {
			logger.Error(fmt.Sprintf("failed to send email to: %s", email), "err", err)
			failedEmails = append(failedEmails, email)
		}
	}

	if len(failedEmails) == len(emails) {
		return &SendRateInfoOutput{
			FailedEmails: failedEmails,
		}, ErrSendRateInfoFailedToSendToAllEmails
	}

	logger.Info("successfully sent rate info")
	return &SendRateInfoOutput{
		FailedEmails: failedEmails,
	}, nil
}
