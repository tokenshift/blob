package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := fmt.Sprintf(":%s", GetEnvOr("PORT", "3000"))
	http.ListenAndServe(port, Responder{})
}
