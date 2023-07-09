package service

import (
	"context"
	"errors"

	"github.com/vadimpk/gses-2023/config"
	"github.com/vadimpk/gses-2023/core/pkg/logging"
)

type Services struct {
	Email EmailService
}

type Options struct {
	Storages Storages
	APIs     APIs
	Logger   logging.Logger
	Cfg      *config.Config
}

type serviceContext struct {
	storages Storages
	apis     APIs
	logger   logging.Logger
	cfg      *config.Config
}

// EmailService provides business logic for email service.
type EmailService interface {
	// Subscribe subscribes email to newsletter.
	Subscribe(ctx context.Context, email string) error
	// SendRateInfo sends emails to all subscribers about current rate info.
	SendRateInfo(ctx context.Context) (*SendRateInfoOutput, error)
}

var (
	// ErrSubscribeAlreadySubscribed is returned when email is already subscribed.
	ErrSubscribeAlreadySubscribed = errors.New("already subscribed")

	// ErrSendRateInfoFailedToSendToAllEmails is returned when failed to send rate info to all emails.
	ErrSendRateInfoFailedToSendToAllEmails = errors.New("failed to send rate info to all emails")
)

type SendRateInfoOutput struct {
	FailedEmails []string
}
