package service

import "context"

type Storages struct {
	Email EmailStorage
}

// EmailStorage provides methods for storing emails that are used in EmailService.
type EmailStorage interface {
	// Save saves email to storage.
	Save(ctx context.Context, email string) error
	// List returns list of emails from storage.
	List(ctx context.Context) ([]string, error)
	// Get returns email from storage.
	Get(ctx context.Context, email string) (string, error)
}
