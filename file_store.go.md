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
		DeleteFile(client, path string) (bool, error)
		GetFile(client, path string) (bool, Handle, error)
		SaveFile(client, path, mimeType string, body io.Reader) (bool, error)
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

File metadata is stored in a BoltDB instance, with keys based on the filename.
The keys used are:

* `{filename}:id`  
  The UUID for the file on disk.
  This is kept instead of the full file path, so that the file store can be
  moved transparently to the Blob node.
* `{filename}:hash`  
  The hash of the file. This is computed as the file is saved to disk and saved
  as metadata.
* `{filename}:size`  
  The size (in bytes) of the file. Also computed as the file is saved to disk.
* `{filename}:mimeType`  
  The MIME type of the file, if known. This will usually be the Content-Type
  header from the request that provided the file.

Files for each client are stored in their own BoltDB bucket.

	type localFileStore struct {
		dbFile, storeDir string
		db *bolt.DB
	}

	func bucketKey(filename string, field string) []byte {
		return []byte(filename + ":" + field)
	}

GetFile returns true if the file exists, and provides a handle to the file
metadata (that can also be used to access the file data).

	func (store localFileStore) GetFile(client, fname string) (bool, Handle, error) {
		exists := false
		handle := localFileStoreEntry{}

		err := store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(client))
			if b == nil {
				return nil
			}

			handle.id = b.Get(bucketKey(fname, "id"))
			if handle.id == nil {
				return nil
			}

			exists = true
			handle.path     = filepath.Join(store.storeDir, uuid.UUID(handle.id).String())

			handle.hash     = b.Get(bucketKey(fname, "hash"))
			handle.mimeType = string(b.Get(bucketKey(fname, "mimeType")))

			size := b.Get(bucketKey(fname, "size"))
			handle.size, _ = binary.Varint(size)

			return nil
		})

		return exists, handle, err
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

	func (store localFileStore) SaveFile(client, fname, mimeType string, data io.Reader) (bool, error) {
		id := uuid.NewRandom()
		path := filepath.Join(store.storeDir, id.String())

		log.Debug("Saving", fname, "as", id, "at", path)

		file, err := os.Create(path)
		if err != nil {
			return false, err
		}
		defer file.Close()

The input data is tee'd to a rolling hash and the file on disk. The hash and
the file size are both recorded as metadata.

		hash := murmur3.New128()
		input := io.TeeReader(data, hash)

		size, err := io.Copy(file, input)
		if err != nil {
			return false, err
		}

		log.Debug("Saved", fname, "as", id, "at", path, "-", size, "bytes")

		var isNew bool

		store.db.Update(func(tx *bolt.Tx) error {
			sizeBuf := make([]byte, 8)
			binary.PutVarint(sizeBuf, size)

			b, err := tx.CreateBucketIfNotExists([]byte(client))
			if err != nil {
				return err
			}

			isNew = b.Get(bucketKey(fname, "id")) == nil

			err = b.Put(bucketKey(fname, "id"), id)
			if err != nil {
				return err
			}

			err = b.Put(bucketKey(fname, "size"), sizeBuf)
			if err != nil {
				return err
			}

			// BUG: github.com/spaolacci/murmur3 claims to have a block size of
			// 1, but returns an all-zero sum if provided any less than 16
			// bytes. Need to replace with a different hash or better
			// implementation, or pad the input to match the "real" block size.
			err = b.Put(bucketKey(fname, "hash"), hash.Sum(nil))
			if err != nil {
				return err
			}

			err = b.Put(bucketKey(fname, "mimeType"), []byte(mimeType))
			return err
		})

		return isNew, err
	}

DeleteFile returns true if the file existed and was removed.

	func (store localFileStore) DeleteFile(client, fname string) (bool, error) {
		exists := false

		err := store.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(client))
			if b == nil {
				return nil
			}

			exists = b.Get(bucketKey(fname, "id")) != nil
			if exists {
				return b.Delete([]byte(fname))
			} else {
				return nil
			}
		})

		return exists, err
	}
