# REST Handler

	<<#-->>
	package rest

	import (
		"fmt"
		"net/http"

		"github.com/tokenshift/blob/log"
		. "github.com/tokenshift/blob/manifest"
	)

Responds to HTTP request to store, update and retrieve files.

	type Handler struct {
		http.Handler

		manifest Manifest
	}

	func CreateRestHandler(manifest Manifest) Handler {
		return Handler {
			manifest: manifest,
		}
	}

	func (h Handler) Serve(port int) {
		log.Info("Starting Blob node on port", port)

		portString := fmt.Sprintf(":%d", port)
		http.ListenAndServe(portString, h)
	}

	func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
		log.Info(req.Method, req.URL)
		res.WriteHeader(200)
	}
