# Blob - main.go

Parses and validates options given as environment variables or on the command
line, and starts running a **Blob** node.

	<<#-->>
	package main

	import "fmt"
	import "os"

	import "github.com/tokenshift/env"

## Settings

Almost all **Blob** node configuration is done through environment variables.
These variables--and their keys--are listed here, and made available globally
to the running node.

	const envDataFolderKey = "BLOB_DATA_FOLDER"
	var envDataFolder string

## Entry Point

	func main() {
		envDataFolder, ok := env.Get(envDataFolderKey)
		if !ok {
			fatalError(envDataFolderKey, "is required.")
		}

		fi, err := os.Stat(envDataFolder)
		if os.IsNotExist(err) {
			fatalError("Could not locate data folder", envDataFolder)
		}

		if !fi.IsDir() {
			fatalError(envDataFolder, "is not a folder.")
		}

The logger is started first:

		go runLogger()
	}

	<<#-->>
	// Writes an error to STDERR and terminates the process.
	func fatalError(message...interface{}) {
		fmt.Fprintln(os.Stderr, message...)
		os.Exit(1)
	}
