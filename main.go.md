# Blob

	<<#-->>
	package main

	import (
		"os"

		"github.com/tokenshift/blob/env"
		"github.com/tokenshift/blob/log"
		. "github.com/tokenshift/blob/manifest"
		. "github.com/tokenshift/blob/rest"
	)

Main entry point for a **Blob** node. Loads configuration settings and starts
the node.

	func main() {
		dbFile, ok := env.Get("MANIFEST_DB_FILE")
		if !ok {
			log.Fatal("Missing $MANIFEST_DB_FILE")
			os.Exit(1)
		}

		storeDir, ok := env.Get("MANIFEST_STORE_DIR")
		if !ok {
			log.Fatal("Missing $MANIFEST_STORE_DIR")
			os.Exit(1)
		}

		manifest := CreateManifest(dbFile, storeDir)

		restHandler := CreateRestHandler(manifest)

		port, ok, err := env.GetInt("REST_PORT")
		if !ok {
			log.Fatal("Missing $REST_PORT")
			os.Exit(1)
		}
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		restHandler.Serve(port)
	}
