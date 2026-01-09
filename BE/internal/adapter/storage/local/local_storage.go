package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BaseUploadPath string
	BasePublicURL  string
}

func NewLocalStorage(baseUploadPath, basePublicURL string) *LocalStorage {
	// Ensure upload directory exists on startup
	if _, err := os.Stat(baseUploadPath); os.IsNotExist(err) {
		os.MkdirAll(baseUploadPath, 0755)
	}

	return &LocalStorage{
		BaseUploadPath: baseUploadPath,
		BasePublicURL:  basePublicURL,
	}
}

func (s *LocalStorage) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	// 1. Create Destination Path
	dstPath := filepath.Join(s.BaseUploadPath, filename)

	// 2. Create File
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer dst.Close()

	// 3. Write Content
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	// 4. Return Public URL
	// e.g. /uploads/filename.jpg
	return fmt.Sprintf("%s/%s", s.BasePublicURL, filename), nil
}

func (s *LocalStorage) Delete(ctx context.Context, filename string) error {
	filePath := filepath.Join(s.BaseUploadPath, filename)
	return os.Remove(filePath)
}
