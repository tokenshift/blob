package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rxValidPath = regexp.MustCompile(`^[\/0-9a-z \-_',]+(\.[a-z0-9]+)?$`)

func isValidPath(path string) {
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	logInfo(req.Method, req.URL)

	path := strings.ToLower(req.URL.Path)
	if !rxValidPath.MatchString(path) {
		logInfo("Invalid path:", path)
		res.WriteHeader(400)
		return
	}

	path, err := filepath.Abs(filepath.Join("data", path))
	if err != nil {
		logError("Failed to construct local file path:", err)
		res.WriteHeader(500)
		return
	}

	switch(req.Method) {
	case "GET":
		// Retrieve an existing file.
		fs, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				logInfo("File not found:", path)
				res.WriteHeader(404)
			} else {
				logError("Failed to open file:", path, err)
				res.WriteHeader(500)
			}
			return
		}

		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				logInfo("File not found:", path)
				res.WriteHeader(404)
			} else {
				logError("Failed to open file:", path, err)
				res.WriteHeader(500)
			}
			return
		}

		logInfo("Content-Length:", fs.Size())
		res.Header()["Content-Length"] = []string{fmt.Sprint(fs.Size())}
		io.Copy(res, f)
	case "PUT":
		// Create or update a file.
	case "DELETE":
		// Delete a file.
	default:
		logInfo("Unsupported request method:", req.Method)
		res.WriteHeader(400)
	}
}

func main() {
	port := fmt.Sprintf(":%s", GetEnvOr("PORT", "3000"))
	http.ListenAndServe(port, http.HandlerFunc(handleRequest))
}
