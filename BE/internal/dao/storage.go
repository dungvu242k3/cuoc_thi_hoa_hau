/*
Package dao đôi khi cũng chứa các File Storage Adapter.
Tác dụng của storage.go:
- Abstract hóa việc lưu trữ file (ảnh, video).
- Hiện tại đang 구현 (implement) lưu vào thư mục Local của server (`/uploads`).
- Tương lai có thể dễ dàng thay bằng struct S3Storage mà không phá vỡ kiến trúc.
*/
package dao

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
	if _, err := os.Stat(baseUploadPath); os.IsNotExist(err) {
		os.MkdirAll(baseUploadPath, 0755)
	}

	return &LocalStorage{
		BaseUploadPath: baseUploadPath,
		BasePublicURL:  basePublicURL,
	}
}

// Upload copy stream dữ liệu từ RAM xuống Disk ổ cứng, đảm bảo an toàn bộ nhớ khi file quá to.
func (s *LocalStorage) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	dstPath := filepath.Join(s.BaseUploadPath, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.BasePublicURL, filename), nil
}

func (s *LocalStorage) Delete(ctx context.Context, filename string) error {
	filePath := filepath.Join(s.BaseUploadPath, filename)
	return os.Remove(filePath)
}
