# Admin Interface

	<<#-->>
	package admin

	import (
		"fmt"
		. "net/http"

		"github.com/tokenshift/blob/log"

		"github.com/gorilla/mux"
		"golang.org/x/crypto/sha3"
	)

The Admin Interface is a REST service intended to be exposed on a non-public
port that can be used by administrators of a **Node** cluster to configure
various settings and access internal data. It used HTTP Basic Auth to
authenticate the client and has a single admin user/password, configured using
environment variables.

	type AdminInterface interface {
		Serve(port int)
	}

	type server struct {
		mux.Router
		username string
		passhash string
	}

	func CreateAdminInterface(username, passhash string) AdminInterface {
		s := server{
			username: username,
			passhash: passhash,
		}

		s.HandleFunc("/apps", getAppsHandler).Methods("GET")
		s.HandleFunc("/apps/{name}", putAppHandler).Methods("PUT")

		return &s
	}

Call `Serve` to begin serving the admin interface on the specified port.

	func (s *server) Serve(port int) {
		log.Info("Starting admin interface on port", port)
		ListenAndServe(fmt.Sprintf(":%d", port), s)
	}

Password hashing is provided as a public function so that consumers can create
hashes matching the internal implementation. This is also made available
through the github.com/tokenshift/blob/admin/hash executable.

	func Hash(input string) string {
		hash := sha3.New256()
		hash.Write([]byte(input))
		sum := hash.Sum(nil)
		return fmt.Sprintf("%x", sum)
	}

Client applications can be set up with their own usernames and passwords at the
/apps endpoint.

	func getAppsHandler(res ResponseWriter, req *Request) {
	}

	func putAppHandler(res ResponseWriter, req *Request) {
	}
