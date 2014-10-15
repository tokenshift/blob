package main

import (
	"path/filepath"
	"io"
	"os"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Entry struct {
	ID string
	MimeType string
	Added time.Time
	Path string
	Size int64
}

func (e Entry) Open() (io.Reader, error) {
	f, err := os.Open(filepath.Join("data", e.ID))
	return f, err
}

// This will be replaced by SQLite or similar.
var db = make(map[string]Entry)

func ManifestGet(path string) (Entry, bool) {
	if e, ok := db[path]; ok {
		return e, true
	} else {
		return Entry{}, false
	}
}

func ManifestPut(path string, mimeType string, data io.Reader) (Entry, error) {
	entry := Entry {
		ID: uuid.New(),
		MimeType: mimeType,
		Added: time.Now(),
		Path: path,
	}

	f, err := os.Create(filepath.Join("data", entry.ID))
	if err != nil {
		return entry, err
	}

	n, err := io.Copy(f, data)
	entry.Size = n

	db[path] = entry

	return entry, err
}

func ManifestDelete(path string) {
	delete(db, path)
}
