# Manifest

	<<#-->>
	package manifest

	import (
		"encoding/binary"
		"fmt"
		"io"
		"os"
		"path/filepath"

		"github.com/tokenshift/blob/log"

		"code.google.com/p/go-uuid/uuid"
		"github.com/boltdb/bolt"
		"github.com/spaolacci/murmur3"
	)

The manifest stores files and file metadata.

	type Manifest struct {
		dbFile, storeDir string
		db *bolt.DB
	}

	type File struct {
		id []byte
		path string

		Name string
		MimeType string
		Size int64
		Hash []byte
	}

	func CreateManifest(dbFile, storeDir string) (Manifest, error) {
		fi, err := os.Stat(dbFile)
		if err != nil && !os.IsNotExist(err) {
			return Manifest{}, err
		}
		if err == nil && fi.IsDir() {
			return Manifest{}, fmt.Errorf("%s is a directory", dbFile)
		}

		fi, err = os.Stat(storeDir)
		if err != nil {
			return Manifest{}, err
		}
		if !fi.IsDir() {
			return Manifest{}, fmt.Errorf("%s is not a directory.", dbFile)
		}

		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			return Manifest{}, err
		}

		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			return err
		})

		return Manifest {
			dbFile: dbFile,
			storeDir: storeDir,
			db: db,
		}, nil
	}

	func (m Manifest) Close() {
		m.db.Close()
	}

Filenames are mapped to UUIDs before they are stored on disk. This avoids any
collisions when a file is being updated at the same time as it is being read;
the update will write data to a new file (with a different UUID), with the
manifest only being updated once the write is complete.

File metadata is stored in a Bolt DB. The filename itself is is used as a key
where an internal ID for the file (a UUID) is saved. While the same filename
may refer to multiple versions of a file (a Get request will retrieve whatever
the latest version happens to be), the ID refers to a specific version. The ID
also becomes part of the file location on disk.

All other file metadata is stored at keys consisting of the ID with an
additional suffix of the name of the metadata field (e.g. `{id}size`).

	var bucketName = []byte("Files")

	func bucketKey(id []byte, key string) []byte {
		// Not sure why a copy is needed here; if I try to append to `id`
		// directly, I'm getting runtime panics from memmove_amd64.s.
		result := make([]byte, len(id), len(id) + len(key))
		copy(result, id)
		result = append(result, []byte(key)...)
		return result
	}

Get returns true if the file exists, and provides a handle to the file metadata
(that can also be used to access the file data).

	func (m Manifest) Get(fname string) (bool, File) {
		exists := false
		info := File{
			Name: fname,
		}

		m.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			id := b.Get([]byte(fname))

			exists = id != nil
			if !exists {
				return nil
			}

			info.id = id
			info.path = filepath.Join(m.storeDir, uuid.UUID(id).String())

			size := b.Get(bucketKey(id, "size"))
			if size != nil {
				info.Size, _ = binary.Varint(size)
			}

			mimeType := b.Get(bucketKey(id, "mime"))
			if mimeType != nil {
				info.MimeType = string(mimeType)
			}

			hash := b.Get(bucketKey(id, "hash"))
			if hash != nil {
				info.Hash = hash
			}

			return nil
		})

		return exists, info
	}

	func (f File) WriteTo(out io.Writer) error {
		in, err := os.Open(f.path)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		return err
	}

Put returns true if the file was newly created, or false if it already existed
(and was updated).

	func (m Manifest) Put(fname, mimeType string, data io.Reader) (bool, error) {
		id := uuid.NewRandom()
		path := filepath.Join(m.storeDir, id.String())

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

		var isNew bool

		m.db.Update(func(tx *bolt.Tx) error {
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

Delete returns true if the file existed and was removed.

	func (m Manifest) Delete(fname string) (bool, error) {
		var err error
		var exists bool

		m.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))

			exists = b.Get([]byte(fname)) != nil
			if exists {
				err = b.Delete([]byte(fname))
			}

			return err
		})

		return exists, err
	}
