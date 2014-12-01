# REST Handler

	<<#-->>
	package rest

	import (
		"github.com/tokenshift/blob/log"
		. "github.com/tokenshift/blob/manifest"
	)

Responds to HTTP request to store, update and retrieve files.

	type Handler struct {
		manifest Manifest
	}

	func CreateRestHandler(manifest Manifest) Handler {
		return Handler {
			manifest: manifest,
		}
	}

	func (h Handler) Serve(port int) {
		log.Info("Starting Blob node on port", port)

		//portString := fmt.Sprintf(":%d", port)
	}
