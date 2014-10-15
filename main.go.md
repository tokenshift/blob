# Blob

	<<#-->>

	package main

	import (
		"fmt"
		"net/http"
	)

Starts a **Blob** node on the configured port.

	func main() {
		port := fmt.Sprintf(":%s", GetEnvOr("PORT", "3000"))
		http.ListenAndServe(port, http.HandlerFunc(handleRequest))
	}
