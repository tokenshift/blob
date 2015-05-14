	<<#-->>
	package main

	import (
		"io"
	)

The FileStore stores files and file metadata.

	type FileStore interface {
		DeleteFile(path string) (bool, error)
		GetFile(path string) (bool, Handle)
		SaveFile(path, mimeType string, body io.Reader) (bool, error)
	}

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

The standard implementation stores files on the local filesystem, and file
metadata within a BoltDB instance. Filenames are mapped to UUIDs before they
are stored on disk. This avoids any collisions when a file is being updated at
the same time as it is being read; the update will write data to a new file
(with a different UUID), with the manifest only being updated once the write is
complete.

File metadata is stored in a Bolt DB. The filename itself is is used as a key
where an internal ID for the file (a UUID) is saved. While the same filename
may refer to multiple versions of a file (a GET request will retrieve whatever
the latest version happens to be), the ID refers to a specific version. The ID
also becomes part of the file location on disk.

All other file metadata is stored at keys consisting of the ID with an
additional suffix of the name of the metadata field (e.g. {id}size).

	type localFileStore struct {
	}

	type localFileStoreEntry struct {
		size int
		mimeType string
		hash []byte

		path string
	}

GetFile returns true if the file exists, and provides a handle to the file
metadata (that can also be used to access the file data).

	func (store localFileStore) GetFile(path string) (bool, Handle) {
		return false, nil
	}

	func (h localFileStoreEntry) WriteTo(out io.Writer) error {
		return nil
	}

SaveFile returns true if the file was newly created, or false if it already
existed (and was updated).

	func (store localFileStore) SaveFile(path, mimeType string, body io.Reader) (bool, error) {
		// The input data is tee'd to a rolling hash and the file on disk. The hash and the file size are both recorded as metadata.
		return true, nil
	}

DeleteFile returns true if the file existed and was removed.

	func (store localFileStore) DeleteFile(path string) (bool, error) {
		return false, nil
	}
