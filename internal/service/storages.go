package service

import "context"

type Storages struct {
	Email EmailStorage
}

// EmailStorage provides methods for storing emails that are used in EmailService.
//
//go:generate go run github.com/vektra/mockery/v2@v2.27.1 --dir . --name EmailStorage --output ../../internal/service/mocks
type EmailStorage interface {
	// Save saves email to storage.
	Save(ctx context.Context, email string) error
	// List returns list of emails from storage.
	List(ctx context.Context) ([]string, error)
	// Get returns email from storage.
	Get(ctx context.Context, email string) (string, error)
}
