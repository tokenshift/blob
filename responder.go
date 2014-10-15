package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

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

var rxValidPath = regexp.MustCompile(`^[\/0-9a-z \-_',]+(\.[a-z0-9]+)?$`)

func isValidPath(path string) bool {
	return rxValidPath.MatchString(path)
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	NewResponder().Handle(res, req)
}

func (r Responder) Handle(res http.ResponseWriter, req *http.Request) {
	LogInfo(r, req.Method, req.URL)

	path := strings.ToLower(req.URL.Path)
	if !isValidPath(path) {
		LogInfo(r, "Invalid path:", path)
		res.WriteHeader(400)
		return
	}

	path, err := filepath.Abs(filepath.Join("data", path))
	if err != nil {
		LogError(r, "Failed to construct local file path:", err)
		res.WriteHeader(500)
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
	entry, ok := ManifestGet(path)
	if ok {
		res.Header()["Content-Length"] = []string{fmt.Sprint(entry.Size)}
		res.Header()["Content-Type"] = []string{entry.MimeType}
		data, err := entry.Open()
		if err != nil {
			LogError(r, err)
			res.WriteHeader(500)
		}
		io.Copy(res, data)
	} else {
		LogInfo(r, "File not found:", path)
		res.WriteHeader(404)
	}
}

func (r Responder) handleDelete(res http.ResponseWriter, path string) {
	if _, ok := ManifestGet(path); ok {
		ManifestDelete(path)
	} else {
		res.WriteHeader(404)
	}
}

func (r Responder) handlePut(res http.ResponseWriter, req *http.Request, path string) {
	mimeType := ""
	mimeTypes := req.Header["Content-Type"]
	if len(mimeTypes) == 1 {
		mimeType = mimeTypes[0]
	}

	_, err := ManifestPut(path, mimeType, req.Body)
	if err != nil {
		LogError(r, err)
		res.WriteHeader(500)
		return
	}
}
