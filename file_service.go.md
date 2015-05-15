	<<#-->>
	package main

	import (
		"fmt"
		"net/http"
		"sync"

		"github.com/tokenshift/env"
		"github.com/tokenshift/log"
	)

Responds to HTTP requests to store, update and retrieve files.

	type FileService interface {
		Start(*sync.WaitGroup)
	}

	func NewFileService(fileStore FileStore, clientStore ClientStore) (FileService, error) {
		port := env.MustGetInt("BLOB_FILE_SERVICE_PORT")

		return httpFileService {
			fileStore: fileStore,
			port: port,
		}, nil
	}

The HTTP service handles GET, PUT, and DELETE requests to update and retrieve files.

	type httpFileService struct {
		fileStore FileStore
		port int
	}

	func (svc httpFileService) Start(wait *sync.WaitGroup) {
		log.Info("Starting file service on port", svc.port)
		http.ListenAndServe(fmt.Sprintf(":%d", svc.port), svc)
		log.Info("Stopping file service.")
		wait.Done()
	}

	func (svc httpFileService) ServeHTTP(res http.ResponseWriter, req *http.Request) {
		log.Info(req.Method, req.URL)

		switch req.Method {
		case "GET":
			svc.handleGet(res, req)
		case "HEAD":
			svc.handleHead(res, req)
		case "PUT":
			svc.handlePut(res, req)
		case "DELETE":
			svc.handleDelete(res, req)
		default:
			log.Debug("Invalid method:", req.Method)
			res.WriteHeader(405)
			res.Write([]byte("Method not allowed"))
		}
	}

GET requests retrieve an existing file.

	func (svc httpFileService) handleGet(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		exists, handle := svc.fileStore.GetFile(path)

		if exists {
			log.Debug("Retrieving", path)

			if handle.Size() > 0 {
				res.Header()["Content-Length"] = []string{fmt.Sprint(handle.Size())}
			}

			if handle.MimeType() != "" {
				res.Header()["Content-Type"] = []string{handle.MimeType()}
			}

			if handle.Hash() != nil {
				res.Header()["ETag"] = []string{fmt.Sprintf("%x", handle.Hash())}
			}

			err := handle.WriteTo(res)
			if err != nil {
				log.Error(err)
			}
		} else {
			log.Debug(path, "does not exist")
			res.WriteHeader(404)
		}
	}

HEAD requests retrieve metadata for an existing file.

	func (svc httpFileService) handleHead(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		exists, handle := svc.fileStore.GetFile(path)

		if exists {
			log.Debug("Retrieving metadata for", path)

			if handle.Size() > 0 {
				res.Header()["Content-Length"] = []string{fmt.Sprint(handle.Size())}
			}

			if handle.MimeType() != "" {
				res.Header()["Content-Type"] = []string{handle.MimeType()}
			}

			if handle.Hash() != nil {
				res.Header()["ETag"] = []string{fmt.Sprintf("%x", handle.Hash())}
			}

			res.WriteHeader(200)
		} else {
			log.Debug(path, "does not exist")
			res.WriteHeader(404)
		}
	}

PUT requests store a new file or update an existing file.

	func (svc httpFileService) handlePut(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		mimeType := ""
		mimeTypes := req.Header["Content-Type"]
		if len(mimeTypes) > 0 {
			mimeType = mimeTypes[0]
		}

		isNew, err := svc.fileStore.SaveFile(path, mimeType, req.Body)
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

	func (svc httpFileService) handleDelete(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		deleted, err := svc.fileStore.DeleteFile(path)
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
