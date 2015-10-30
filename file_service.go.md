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

	func NewFileService(fileStore FileStore, clientStore ClientStore, siblings Siblings) (FileService, error) {
		port := env.MustGetInt("BLOB_FILE_SERVICE_PORT")

		return httpFileService {
			clientStore: clientStore,
			siblings: siblings,
			fileStore: fileStore,
			port: port,
		}, nil
	}

The HTTP service handles GET, PUT, and DELETE requests to update and retrieve files.

	type httpFileService struct {
		clientStore ClientStore
		siblings Siblings
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

		username, ok := svc.authenticate(req)
		if !ok {
			res.WriteHeader(401)
			res.Write([]byte("Invalid username/password\n"))
			return
		}

		switch req.Method {
		case "GET":
			svc.handleGet(username, res, req)
		case "HEAD":
			svc.handleHead(username, res, req)
		case "PUT":
			svc.handlePut(username, res, req)
		case "DELETE":
			svc.handleDelete(username, res, req)
		default:
			log.Debug("Invalid method:", req.Method)
			res.WriteHeader(405)
			res.Write([]byte("Method not allowed\n"))
		}
	}

All requests are authenticated using HTTP basic auth, and scoped to the
specific client. See the admin service for details on adding/configuring API
clients.

	func (svc httpFileService) authenticate(req *http.Request) (string, bool) {
		if username, password, ok := req.BasicAuth(); ok {
			ok, err := svc.clientStore.VerifyUser(username, password)
			if err != nil || !ok {
				return "", false
			} else {
				return username, true
			}
		} else {
			return "", false
		}
	}

GET requests retrieve an existing file.

	func (svc httpFileService) handleGet(client string, res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		exists, handle, err := svc.fileStore.GetFile(client, path)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("Internal error\n"))
		}

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

	func (svc httpFileService) handleHead(client string, res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		exists, handle, err := svc.fileStore.GetFile(client, path)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("Internal error\n"))
		}

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

	func (svc httpFileService) handlePut(client string, res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		mimeType := ""
		mimeTypes := req.Header["Content-Type"]
		if len(mimeTypes) > 0 {
			mimeType = mimeTypes[0]
		}

		isNew, err := svc.fileStore.SaveFile(client, path, mimeType, req.Body)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown error occurred.\n"))
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

	func (svc httpFileService) handleDelete(client string, res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		deleted, err := svc.fileStore.DeleteFile(client, path)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown eror occurred.\n"))
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
