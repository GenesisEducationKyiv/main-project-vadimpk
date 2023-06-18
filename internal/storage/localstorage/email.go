package localstorage

import (
	"context"
	"strings"

	"github.com/vadimpk/gses-2023/pkg/database"
)

type emailStorage struct {
	db *database.FileDB
}

func NewEmailStorage(db *database.FileDB) *emailStorage {
	return &emailStorage{
		db: db,
	}
}

const (
	emailStorageFileName = "emails.txt"
)

func (s *emailStorage) Save(ctx context.Context, email string) error {
	return s.db.Append(ctx, emailStorageFileName, []byte(email))
}

func (s *emailStorage) List(ctx context.Context) ([]string, error) {
	data, err := s.db.Read(ctx, emailStorageFileName)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	emails := strings.Split(string(data), "\n")

	// Filter out any empty strings that may occur due to trailing new lines
	filteredEmails := emails[:0]
	for _, email := range emails {
		if email != "" {
			filteredEmails = append(filteredEmails, email)
		}
	}

	return filteredEmails, nil
}

func (s *emailStorage) Get(ctx context.Context, email string) (string, error) {
	emails, err := s.List(ctx)
	if err != nil {
		return "", err
	}
	if len(emails) == 0 {
		return "", nil
	}

	for _, e := range emails {
		if e == email {
			return e, nil
		}
	}

	return "", nil
}
