	<<#-->>
	package main

	import (
		"os"
		"sync"

		"github.com/tokenshift/log"
	)

The main entry point for a Blob node. All settings are configured by
environment variables.

	func main() {
		var err error

		defer func() {
			if r := recover(); r != nil {
				log.Fatal(r)
				os.Exit(1)
			}
		}()

Components are initialized one at a time in dependency order and then injected.

		fileStore, err := NewFileStore()
		if err != nil {
			log.Fatal("Failed to initialize file store")
			log.Fatal(err)
			os.Exit(1)
		}

		clientStore, err := NewClientStore()
		if err != nil {
			log.Fatal("Failed to initialize client store")
			log.Fatal(err)
			os.Exit(1)
		}

		siblings := NewSiblingStore()

		fileService, err := NewFileService(fileStore, clientStore, siblings)
		if err != nil {
			log.Fatal("Failed to initialize file service")
			log.Fatal(err)
			os.Exit(1)
		}

		adminService, err := NewAdminService(clientStore, siblings)
		if err != nil {
			log.Fatal("Failed to initialize admin service")
			log.Fatal(err)
			os.Exit(1)
		}

Once all settings have been validated, the individual components are
initialized and started.

		var wg sync.WaitGroup
		wg.Add(2)
		go fileService.Start(&wg)
		go adminService.Start(&wg)
		wg.Wait()
	}
