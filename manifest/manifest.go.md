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

Put returns true if the file was newly created, or false if it already existed
(and was updated).

	func (m Manifest) Put(fname string, data io.Reader) (bool, error) {
		id := uuid.NewUUID().String()
		path := filepath.Join(m.storeDir, id)

		log.Debug("Saving", fname, "as", id, "at", path)

		var isNew bool
		var err error

		m.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return err
			}
		
			isNew = b.Get([]byte(fname)) == nil

			err = b.Put([]byte(fname), []byte(id))
			return err
		})

		return isNew, err
	}
