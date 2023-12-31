package localstorage

import (
	"context"
	"strings"

	"github.com/vadimpk/gses-2023/core/pkg/database"
)

type emailStorage struct {
	db       *database.FileDB
	filename string
}

func NewEmailStorage(db *database.FileDB, filename string) *emailStorage {
	return &emailStorage{
		db:       db,
		filename: filename,
	}
}

func (s *emailStorage) Save(ctx context.Context, email string) error {
	return s.db.Append(ctx, s.filename, []byte(email))
}

func (s *emailStorage) List(ctx context.Context) ([]string, error) {
	data, err := s.db.Read(ctx, s.filename)
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

func (s *emailStorage) Exist(ctx context.Context, email string) (bool, error) {
	emails, err := s.List(ctx)
	if err != nil {
		return false, err
	}

	if len(emails) == 0 {
		return false, nil
	}

	for _, e := range emails {
		if e == email {
			return true, nil
		}
	}

	return false, nil
}
