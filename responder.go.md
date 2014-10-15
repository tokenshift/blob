# Responder

	<<#-->>

	package main

	import (
		"fmt"
		"io"
		"math/rand"
		"net/http"
		"regexp"
		"strings"
	)

Responds to HTTP requests to the node.

A new Responder is created for each request, and acts as a request context,
tracking an ID that can be used to correlate log entries.

	type Responder struct {
		id int
	}

	func NewResponder() Responder {
		return Responder {
			id: rand.Int() % 0xffffffff,
		}
	}

	func (r Responder) ID() int {
		return r.id
	}

	func handleRequest(res http.ResponseWriter, req *http.Request) {
		NewResponder().Handle(res, req)
	}

Valid request methods include GET (retrieve a file), PUT (create/update a file)
and DELETE (remove a file).

	func (r Responder) Handle(res http.ResponseWriter, req *http.Request) {
		LogInfo(r, req.Method, req.URL)

		path := strings.ToLower(req.URL.Path)
		if !isValidPath(path) {
			LogInfo(r, "Invalid path:", path)
			res.WriteHeader(400)
			return
		}

		switch(req.Method) {
		case "GET":
			// Retrieve an existing file.
			r.handleGet(res, path)
		case "PUT":
			// Create or update a file.
			r.handlePut(res, req, path)
		case "DELETE":
			// Delete a file.
			r.handleDelete(res, path)
		default:
			LogInfo(r, "Unsupported request method:", req.Method)
			res.WriteHeader(400)
		}
	}

	func (r Responder) handleGet(res http.ResponseWriter, path string) {
		entry, err := manifest.Get(path)
		if err != nil {
			if _, ok := err.(NotFound); ok {
				LogInfo(r, err)
				res.WriteHeader(404)
			} else {
				LogError(r, err)
				res.WriteHeader(500)
			}
			return
		}

		res.Header()["Content-Length"] = []string{fmt.Sprint(entry.Size)}
		res.Header()["Content-Type"] = []string{entry.MimeType}
		data, err := entry.Open()
		if err != nil {
			LogError(r, err)
			res.WriteHeader(500)
		}
		io.Copy(res, data)
	}

	func (r Responder) handleDelete(res http.ResponseWriter, path string) {
		_, err := manifest.Get(path)
		if err == nil {
			err = manifest.Delete(path)
		}

		if err == nil {
			return
		} else if _, ok := err.(NotFound); ok {
			LogInfo(r, err)
			res.WriteHeader(404)
		} else {
			LogError(r, err)
			res.WriteHeader(500)
			res.Write([]byte(err.Error()))
		}
	}

	func (r Responder) handlePut(res http.ResponseWriter, req *http.Request, path string) {
		mimeType := ""
		mimeTypes := req.Header["Content-Type"]
		if len(mimeTypes) == 1 {
			mimeType = mimeTypes[0]
		}

		err := manifest.Put(path, mimeType, req.Body)
		if err != nil {
			LogError(r, err)

			if br, ok := err.(BadRequest); ok {
				res.WriteHeader(400)
				res.Write([]byte(br))
			} else {
				res.WriteHeader(500)
			}
		}
	}

Only a restricted subset of possible file paths is permitted to ensure safety
and platform-agnostic behavior.

	var rxValidPath = regexp.MustCompile(`^[\/0-9a-z \-_',]+(\.[a-z0-9]+)?$`)

	func isValidPath(path string) bool {
		return rxValidPath.MatchString(path)
	}
