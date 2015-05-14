package main

import (
	"os"
	"sync"

	"github.com/tokenshift/blob/log"
)

func main() {
	var err error
	// Initialize and inject all components.

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

	fileService, err := NewFileService(fileStore, clientStore)
	if err != nil {
		log.Fatal("Failed to initialize file service")
		log.Fatal(err)
		os.Exit(1)
	}

	adminService, err := NewAdminService(clientStore)
	if err != nil {
		log.Fatal("Failed to initialize admin service")
		log.Fatal(err)
		os.Exit(1)
	}

	// Start the public interfaces.

	var wg sync.WaitGroup
	wg.Add(2)
	go fileService.Start(&wg)
	go adminService.Start(&wg)
	wg.Wait()
}
