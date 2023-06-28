package database

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

type FileDB struct {
	baseFilePath string
}

func NewFileDB(baseFilePath string) *FileDB {
	return &FileDB{
		baseFilePath: baseFilePath,
	}
}

func (f *FileDB) Close() error {
	return nil
}

func (f *FileDB) Ping(ctx context.Context) error {
	_, err := os.Stat(f.baseFilePath)
	if os.IsNotExist(err) {
		return errors.New("base file path does not exist")
	}

	return err
}

var writePermissionCode = 0600

func (f *FileDB) Append(ctx context.Context, file string, data []byte) error {
	fullPath := filepath.Join(f.baseFilePath, file)

	// Open the file in append mode
	fh, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(writePermissionCode))
	if err != nil {
		return err
	}
	defer fh.Close()

	// Append new data onto the next line
	_, err = fh.Write(append(data, '\n'))

	return err
}

func (f *FileDB) Read(ctx context.Context, file string) ([]byte, error) {
	fullPath := filepath.Join(f.baseFilePath, file)

	// Check if file exists, if not return an error
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return os.ReadFile(fullPath)
}
