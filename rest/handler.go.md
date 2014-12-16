# REST Handler

	<<#-->>
	package rest

	import (
		"fmt"
		"net/http"
		"regexp"

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

The REST service handles GET, PUT, and DELETE requests to update and retrieve
files.

	func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
		log.Info(req.Method, req.URL)

		path := req.URL.Path
		if !isValidFilename(path) {
			log.Debug("Invalid filename:", path)
			res.WriteHeader(400)
			res.Write([]byte("Invalid filename"))
			return
		}

		switch req.Method {
		case "GET":
			h.handleGet(res, req)
		case "PUT":
			h.handlePut(res, req)
		case "DELETE":
			h.handleDelete(res, req)
		default:
			log.Debug("Invalid method:", req.Method)
			res.WriteHeader(405)
			res.Write([]byte("Method not allowed"))
		}
	}

Only a restricted subset of possible filenames is supported.

	var rxValidFilename = regexp.MustCompile(`(?i)[0-9a-z\-_\.]+`)

	func isValidFilename(fname string) bool {
		return rxValidFilename.MatchString(fname)
	}

GET requests retrieve an existing file.

	func (h Handler) handleGet(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		exists, info := h.manifest.Get(path)

		if exists {
			log.Debug("Retrieving", path)

			if info.Size > 0 {
				res.Header()["Content-Length"] = []string{fmt.Sprint(info.Size)}
			}

			if info.MimeType != "" {
				res.Header()["Content-Type"] = []string{info.MimeType}
			}

			if info.Hash != nil {
				res.Header()["ETag"] = []string{fmt.Sprintf("%x", info.Hash)}
			}

			err := info.WriteTo(res)
			if err != nil {
				log.Error(err)
			}
		} else {
			log.Debug(path, "does not exist")
			res.WriteHeader(404)
		}
	}

PUT requests store a new file or update an existing file.

	func (h Handler) handlePut(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		mimeType := ""
		mimeTypes := req.Header["Content-Type"]
		if len(mimeTypes) > 0 {
			mimeType = mimeTypes[0]
		}

		isNew, err := h.manifest.Put(path, mimeType, req.Body)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown error occurred."))
			return
		}

		if isNew {
			log.Debug("Created new file", path)
			res.WriteHeader(201)
		} else {
			log.Debug("Updated existing file", path)
			res.WriteHeader(200)
		}
	}

DELETE requests remove a file from the node.

	func (h Handler) handleDelete(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		deleted, err := h.manifest.Delete(path)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown eror occurred."))
			return
		}

		if deleted {
			log.Debug(path, "was deleted")
			res.WriteHeader(200)
		} else {
			log.Debug(path, "was not found")
			res.WriteHeader(404)
		}
	}
