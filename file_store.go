package main

import (
	"io"
)

// Stores static files with mime types and other metadata.
type FileStore interface {
	DeleteFile(path string) (bool, error)
	GetFile(path string) (bool, Handle)
	SaveFile(path, mimeType string, body io.Reader) (bool, error)
}

// A handle to retrieve a file and its metadata.
type Handle interface {
	Size() int
	MimeType() string
	Hash() []byte

	WriteTo(io.Writer) error
}

func NewFileStore() (FileStore, error) {
	return localFileStore {
	}, nil
}

// Local file storage that uses BoltDB for metadata and stores files in the
// local filesystem.
type localFileStore struct {
}

type localFileStoreEntry struct {
	size int
	mimeType string
	hash []byte

	path string
}

func (store localFileStore) DeleteFile(path string) (bool, error) {
	return false, nil
}

func (store localFileStore) GetFile(path string) (bool, Handle) {
	return false, nil
}

func (store localFileStore) SaveFile(path, mimeType string, body io.Reader) (bool, error) {
	return true, nil
}

func (h localFileStoreEntry) WriteTo(out io.Writer) error {
	return nil
}
