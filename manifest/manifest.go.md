# Manifest

	<<#-->>
	package manifest

	import (
		//"github.com/boltdb/bolt"
		//"github.com/peterbourgon/diskv"
	)

Stores files and file metadata.

	type Manifest struct {
		dbFile, storeDir string
	}

	func CreateManifest(dbFile, storeDir string) Manifest {
		return Manifest {
			dbFile: dbFile,
			storeDir: storeDir,
		}
	}
