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
	path := strings.ToLower(req.URL.Path)
	fmt.Println(path)
	if !rxValidPath.MatchString(path) {
		res.WriteHeader(400)
		return
	}

	path, err := filepath.Abs(filepath.Join("data", path))
	if err != nil {
		res.WriteHeader(500)
		return
	}

	switch(req.Method) {
	case "GET":
		// Retrieve an existing file.
		fs, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				res.WriteHeader(404)
			} else {
				res.WriteHeader(500)
			}
			return
		}

		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				res.WriteHeader(404)
			} else {
				res.WriteHeader(500)
			}
			return
		}

		res.Header()["Content-Length"] = []string{fmt.Sprint(fs.Size())}
		io.Copy(res, f)
	case "PUT":
		// Create or update a file.
	case "DELETE":
		// Delete a file.
	default:
		res.WriteHeader(400)
	}
}

func main() {
	port := fmt.Sprintf(":%s", GetEnvOr("PORT", "3000"))
	http.ListenAndServe(port, http.HandlerFunc(handleRequest))
}
