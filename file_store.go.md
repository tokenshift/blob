	<<#-->>
	package main

	import (
		"encoding/binary"
		"fmt"
		"io"
		"os"
		"path/filepath"

		"code.google.com/p/go-uuid/uuid"
		"github.com/boltdb/bolt"
		"github.com/spaolacci/murmur3"

		"github.com/tokenshift/env"
		"github.com/tokenshift/log"
	)

The FileStore stores files and file metadata.

	type FileStore interface {
		DeleteFile(path string) (bool, error)
		GetFile(path string) (bool, Handle)
		SaveFile(path, mimeType string, body io.Reader) (bool, error)
	}

	func NewFileStore() (FileStore, error) {
		dbFile := env.MustGet("BLOB_FILE_STORE_DB")
		storeDir := env.MustGet("BLOB_FILE_STORE_DIR")

		fi, err := os.Stat(dbFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if err == nil && fi.IsDir() {
			return nil, fmt.Errorf("%s is a directory", dbFile)
		}

		fi, err = os.Stat(storeDir)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("%s is not a directory.", storeDir)
		}

		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			return nil, err
		}

		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			return err
		})

		return localFileStore {
			dbFile: dbFile,
			storeDir: storeDir,
			db: db,
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
		dbFile, storeDir string
		db *bolt.DB
	}

	var bucketName = []byte("Files")

	func bucketKey(id []byte, key string) []byte {
		// Not sure why a copy is needed here; if I try to append to `id`
		// directly, I'm getting runtime panics from memmove_amd64.s.
		result := make([]byte, len(id), len(id) + len(key))
		copy(result, id)
		result = append(result, []byte(key)...)
		return result
	}

GetFile returns true if the file exists, and provides a handle to the file
metadata (that can also be used to access the file data).

	func (store localFileStore) GetFile(fname string) (bool, Handle) {
		exists := false
		handle := localFileStoreEntry{}

		store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			id := b.Get([]byte(fname))

			exists = id != nil
			if !exists {
				return nil
			}

			handle.id = id
			handle.path = filepath.Join(store.storeDir, uuid.UUID(id).String())

			size := b.Get(bucketKey(id, "size"))
			if size != nil {
				handle.size, _ = binary.Varint(size)
			}

			mimeType := b.Get(bucketKey(id, "mime"))
			if mimeType != nil {
				handle.mimeType = string(mimeType)
			}

			hash := b.Get(bucketKey(id, "hash"))
			if hash != nil {
				handle.hash = hash
			}

			return nil
		})

		return exists, handle
	}

	type Handle interface {
		Size() int64
		MimeType() string
		Hash() []byte

		WriteTo(io.Writer) error
	}

	type localFileStoreEntry struct {
		size int64
		mimeType string
		hash []byte

		id []byte
		path string
	}

	func (h localFileStoreEntry) Size() int64 {
		return h.size
	}

	func (h localFileStoreEntry) MimeType() string {
		return h.mimeType
	}

	func (h localFileStoreEntry) Hash() []byte {
		return h.hash
	}

	func (h localFileStoreEntry) WriteTo(out io.Writer) error {
		in, err := os.Open(h.path)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		return err
	}

SaveFile returns true if the file was newly created, or false if it already
existed (and was updated).

	func (store localFileStore) SaveFile(fname, mimeType string, data io.Reader) (bool, error) {
		id := uuid.NewRandom()
		path := filepath.Join(store.storeDir, id.String())

		log.Debug("Saving", fname, "as", id, "at", path)

		file, err := os.Create(path)
		if err != nil {
			return false, err
		}
		defer file.Close()

The input data is tee'd to a rolling hash and the file on disk. The hash and the file size are both recorded as metadata.

		hash := murmur3.New128()
		input := io.TeeReader(data, hash)

		size, err := io.Copy(file, input)
		if err != nil {
			return false, err
		}

		var isNew bool

		store.db.Update(func(tx *bolt.Tx) error {
			sizeBuf := make([]byte, 8)
			binary.PutVarint(sizeBuf, size)

			b := tx.Bucket([]byte(bucketName))
			isNew = b.Get([]byte(fname)) == nil

			err = b.Put([]byte(fname), id)
			if err != nil {
				return err
			}

			err = b.Put(bucketKey(id, "size"), sizeBuf)
			if err != nil {
				return err
			}

			err = b.Put(bucketKey(id, "hash"), hash.Sum(nil))
			if err != nil {
				return err
			}

			err = b.Put(bucketKey(id, "mime"), []byte(mimeType))
			return err
		})

		return isNew, err
	}

DeleteFile returns true if the file existed and was removed.

	func (store localFileStore) DeleteFile(fname string) (bool, error) {
		var err error
		var exists bool

		store.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))

			exists = b.Get([]byte(fname)) != nil
			if exists {
				err = b.Delete([]byte(fname))
			}

			return err
		})

		return exists, err
	}
