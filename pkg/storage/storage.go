package storage

import (
	"context"
	"io"
)

// Storage represents under storage interface
type Storage interface {
	// Use a namespace to provide the readonly filesystem
	Use(ctx context.Context, namespace string) (ReadonlyFs, error)
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
