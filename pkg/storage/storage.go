package storage

import (
	"context"
	"io"
)

// Storage represents under storage interface
type Storage interface {
	// Get readonly filesystem from storage
	Get(ctx context.Context) (ReadonlyFs, error)

	// Use readonly file with a namespace
	Use(ctx context.Context, namespace string) (ReadonlyFile, error)

	// Env returns storage environment
	Env(ctx context.Context) string
}

// ReadonlyFs is a readonly filesystem
type ReadonlyFs interface {
	io.Closer

	// Open a readonly file for reading, the associated file descriptor has O_RDONLY
	Open(filename string) (ReadonlyFile, error)
}

// ReadonlyFile is a readonly file
type ReadonlyFile interface {
	io.ReadCloser
}
