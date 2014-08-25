# Blob - Request Handler

	<<#-->>
	package main

	import "fmt"
	import "io"
	import "net/http"

The request handler listens on the configured port (3103 by default) for file
requests and dispatches them to workers, after checking them against the
manifest and ensuring the requestor is authorized.

A request is represented as a struct containing the name of the file that was
requested, and a response channel to which file contents should be written.

	type request struct {
		fname string
		res chan response
	}

	type response struct {
		out io.Writer
	}

	var requestQueue = make(chan request, 100)

	func requestHandler() {
		logDebug("Starting request handler")
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), http.HandlerFunc(handleRequest))
		if err != nil {
			logFatal(err)
		}
	}

	func handleRequest(res http.ResponseWriter, req *http.Request) {
		method := req.Method
		path := req.URL.Path
		logDebug(method, path)
	}
