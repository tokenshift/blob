# Blob

	<<#-->>
	package main

	import (
		"os"

		"github.com/tokenshift/blob/env"
		"github.com/tokenshift/blob/log"

		. "github.com/tokenshift/blob/admin"
		. "github.com/tokenshift/blob/manifest"
		. "github.com/tokenshift/blob/rest"
	)

The main entry point for a **Blob** node. All settings are configured by
environment variables.

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

		adminUsername, ok := env.Get("ADMIN_USERNAME")
		if !ok {
			log.Fatal("Missing $ADMIN_USERNAME")
			os.Exit(1)
		}

		adminPasshash, ok := env.Get("ADMIN_PASSHASH")
		if !ok {
			log.Fatal("Missing $ADMIN_USERNAME")
			os.Exit(1)
		}

		restPort, ok, err := env.GetInt("REST_PORT")
		if !ok {
			log.Fatal("Missing $REST_PORT")
			os.Exit(1)
		}
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		adminPort, ok, err := env.GetInt("ADMIN_PORT")
		if !ok {
			log.Fatal("Missing $ADMIN_PORT")
			os.Exit(1)
		}
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

Once all settings have been validated, the individual components are
initialized and started.

		manifest, err := CreateManifest(dbFile, storeDir)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		defer manifest.Close()

		adminInterface := CreateAdminInterface(adminUsername, adminPasshash)

		restHandler := CreateRestHandler(manifest)

		adminWait := make(chan struct{})
		go func() {
			adminInterface.Serve(adminPort)
			close(adminWait)
		}()

		restWait := make(chan struct{})
		go func() {
			restHandler.Serve(restPort)
			close(restWait)
		}()

		<-adminWait
		<-restWait
	}
