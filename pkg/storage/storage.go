package storage

import "context"

// Storage represents under storage interface
type Storage interface {
	// Use a namespace for storage read
	Use(ctx context.Context, namespace string) (err error)

	// Read content from the file path
	Read(path string) (content string, err error)
}
