package port

import (
	"context"
	"io"
)

// FileStorage defines the behavior for saving and retrieving files.
// This allows us to switch between Local Disk, AWS S3, Google Cloud, etc.
type FileStorage interface {
	// Upload saves the given content and returns the public URL.
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)

	// Delete removes the file (optional for now, but good for cleanup).
	Delete(ctx context.Context, filename string) error
}
