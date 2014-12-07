# Manifest

	<<#-->>
	package manifest

	import (
		"fmt"
		"io"
		"os"
		"path/filepath"

		"github.com/tokenshift/blob/log"

		"code.google.com/p/go-uuid/uuid"
		"github.com/boltdb/bolt"
	)

Stores files and file metadata.

	type Manifest struct {
		dbFile, storeDir string
		db *bolt.DB
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

File metadata is stored in a Bolt DB.

	const bucketName = "Files"

Get returns true if the file exists, and writes the file to the specified
output stream.

	func (m Manifest) Get(fname string, out io.Writer) (bool, error) {
		path := ""

		m.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			id := b.Get([]byte(fname))

			if id != nil {
				path = filepath.Join(m.storeDir, string(id))
			}

			return nil
		})

		if path == "" {
			return false, nil
		}

		file, err := os.Open(path)
		if err != nil {
			return true, err
		}
		defer file.Close()

		_, err = io.Copy(out, file)
		return true, err
	}

Put returns true if the file was newly created, or false if it already existed
(and was updated).

	func (m Manifest) Put(fname string, data io.Reader) (bool, error) {
		id := uuid.NewUUID().String()
		path := filepath.Join(m.storeDir, id)

		log.Debug("Saving", fname, "as", id, "at", path)

		file, err := os.Create(path)
		if err != nil {
			return false, err
		}
		defer file.Close()

		_, err = io.Copy(file, data)
		if err != nil {
			return false, err
		}

		var isNew bool

		m.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			isNew = b.Get([]byte(fname)) == nil
			err = b.Put([]byte(fname), []byte(id))
			return err
		})

		return isNew, err
	}
