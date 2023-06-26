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
	// Exist checks if email exists in storage.
	Exist(ctx context.Context, email string) (bool, error)
}
