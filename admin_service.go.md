	<<#-->>
	package main

	import (
		"crypto/sha256"
		"fmt"
		"net/http"
		"sync"
		"time"

		"github.com/bmizerany/pat"
		"github.com/tokenshift/env"
		"github.com/tokenshift/log"
	)

The admin interface provides a mechanism for configuring a running Blob
instance.

	type AdminService interface {
		Start(*sync.WaitGroup)
	}

	func NewAdminService(clientStore ClientStore, siblings Siblings) (AdminService, error) {
		port     := env.MustGetInt("BLOB_ADMIN_SERVICE_PORT")
		username := env.MustGet("BLOB_ADMIN_SERVICE_USERNAME")
		passhash := env.MustGet("BLOB_ADMIN_SERVICE_PASSHASH")

		svc := httpAdminService {
			clientStore: clientStore,
			siblings: siblings,
			port: port,
			username: username,
			passhash: passhash,
		}

		svc.makeRoutes()

		return svc, nil
	}

The admin service will run as a REST service on a different port to the "main"
file service. All requests to this service are authenticated using HTTP basic
auth, with a username and password hash provided as environment variables.

	type httpAdminService struct {
		clientStore ClientStore
		siblings Siblings
		mux http.Handler
		port int
		username, passhash string
	}

	func (svc httpAdminService) Start(wait *sync.WaitGroup) {
		log.Info("Starting admin service on port", svc.port)
		http.ListenAndServe(fmt.Sprintf(":%d", svc.port), svc)
		log.Info("Stopping admin service.")
		wait.Done()
	}

	func (svc httpAdminService) ServeHTTP(res http.ResponseWriter, req *http.Request) {
		if username, password, ok := req.BasicAuth(); ok {
			if (username == svc.username && Hash(password) == svc.passhash) {
				svc.mux.ServeHTTP(res, req)
			} else {
				// 403; invalid username/password.
				res.WriteHeader(403)
				res.Write([]byte("Invalid username/password\n"))
			}
		} else {
			// 401; authentication required.
			res.WriteHeader(401)
			res.Write([]byte("Authentication required\n"))
		}
	}

For convenience, the admin service provides its own hash function (the same as
is exposed by the github.com/tokenshift/blob/hash command-line utility) to use
when veriying the service credentials.

	func Hash(password string) string {
		h := sha256.Sum256([]byte(password))
		return fmt.Sprintf("%x", h)
	}

Route definitions for the admin service. The service uses [Pat](https://github.com/bmizerany/pat)
for route multiplexing.

The admin interface supports dynamically adding client applications (/clients)
and siblings (/siblings) to the running node.

	func (svc *httpAdminService) makeRoutes() {
		mux := pat.New()

		mux.Get("/clients", http.HandlerFunc(svc.getClients))
		mux.Put("/clients/:username", http.HandlerFunc(svc.putClient))
		mux.Del("/clients/:username", http.HandlerFunc(svc.deleteClient))

		mux.Get("/siblings", http.HandlerFunc(svc.getSiblings))
		mux.Get("/siblings/:uri", http.HandlerFunc(svc.getSibling))
		mux.Post("/siblings", http.HandlerFunc(svc.addSibling))

		svc.mux = mux
	}

	func (svc httpAdminService) getClients(res http.ResponseWriter, req *http.Request) {
		users, err := svc.clientStore.GetUsers()
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown error occurred.\n"))
			return
		}

		res.WriteHeader(200)
		for _, user := range(users) {
			fmt.Fprintln(res, user)
		}
	}

	func (svc httpAdminService) putClient(res http.ResponseWriter, req *http.Request) {
		username := req.URL.Query().Get(":username")
		password := req.FormValue("password")

		if password == "" {
			res.WriteHeader(400)
			res.Write([]byte("Password is required\n"))
			return
		}

		created, err := svc.clientStore.SaveUser(username, password)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown error occurred.\n"))
			return
		}

		if created {
			res.WriteHeader(201)
		} else {
			res.WriteHeader(200)
		}
	}

	func (svc httpAdminService) deleteClient(res http.ResponseWriter, req *http.Request) {
		username := req.URL.Query().Get(":username")

		deleted, err := svc.clientStore.DeleteUser(username)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown error occurred.\n"))
			return
		}

		if deleted {
			res.WriteHeader(200)
		} else {
			res.WriteHeader(404)
			res.Write([]byte("User not found.\n"))
		}
	}

	func (svc httpAdminService) getSiblings(res http.ResponseWriter, req *http.Request) {
		siblings := svc.siblings.All()

		res.WriteHeader(200)
		for _, sibling := range(siblings) {
			fmt.Fprintln(res, sibling.URI)
			fmt.Fprintln(res, "Introduced:  ", sibling.Introduced.Format(time.RFC3339))
			fmt.Fprintln(res, "Last Contact:", sibling.LastContact.Format(time.RFC3339))
			fmt.Fprintln(res, "Estranged:   ", sibling.Estranged)
			fmt.Fprintln(res, "")
		}
	}

	func (svc httpAdminService) getSibling(res http.ResponseWriter, req *http.Request) {
		uri := req.URL.Query().Get(":uri")
		sibling, ok := svc.siblings.Status(uri)

		if ok {
			res.WriteHeader(200)
			fmt.Fprintln(res, sibling.URI)
			fmt.Fprintln(res, "Introduced:  ", sibling.Introduced.Format(time.RFC3339))
			fmt.Fprintln(res, "Last Contact:", sibling.LastContact.Format(time.RFC3339))
			fmt.Fprintln(res, "Estranged:   ", sibling.Estranged)
			fmt.Fprintln(res, "")
		} else {
			res.WriteHeader(404)
			fmt.Fprintln(res, "Sibling not found:", uri)
		}
	}

	func (svc httpAdminService) addSibling(res http.ResponseWriter, req *http.Request) {
		uri := req.FormValue("uri")

		if uri == "" {
			res.WriteHeader(400)
			res.Write([]byte("URI is required\n"))
			return
		}

		status, err := svc.siblings.Add(uri)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown error occurred.\n"))
			return
		}

		res.WriteHeader(200)
		fmt.Fprintln(res, status.URI)
		fmt.Fprintln(res, "Introduced:  ", status.Introduced.Format(time.RFC3339Nano))
		fmt.Fprintln(res, "Last Contact:", status.LastContact.Format(time.RFC3339Nano))
		fmt.Fprintln(res, "Estranged:   ", status.Estranged)
		fmt.Fprintln(res, "")
	}