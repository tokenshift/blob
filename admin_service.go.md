	<<#-->>
	package main

	import (
		"fmt"
		"net/http"
		"sync"

		"github.com/bmizerany/pat"
		"github.com/tokenshift/env"
		"github.com/tokenshift/log"
	)

The admin interface provides a mechanism for configuring a running Blob
instance.

	type AdminService interface {
		Start(*sync.WaitGroup)
	}

	func NewAdminService(clientStore ClientStore) (AdminService, error) {
		port := env.MustGetInt("BLOB_ADMIN_SERVICE_PORT")

		svc := httpAdminService {
			clientStore: clientStore,
			port: port,
		}

		svc.makeRoutes()

		return svc, nil
	}

The admin service will run as a REST service on a different port to the "main"
file service.

	type httpAdminService struct {
		clientStore ClientStore
		mux http.Handler
		port int
	}

	func (svc httpAdminService) Start(wait *sync.WaitGroup) {
		log.Info("Starting admin service on port", svc.port)
		http.ListenAndServe(fmt.Sprintf(":%d", svc.port), svc.mux)
		log.Info("Stopping admin service.")
		wait.Done()
	}

Route definitions for the admin service. The service uses [Pat](https://github.com/bmizerany/pat)
for route multiplexing. Currently, the only routes supported are to add/update
or delete clients (PUT and DELETE).

	func (svc *httpAdminService) makeRoutes() {
		mux := pat.New()

		mux.Put("/clients/:username", http.HandlerFunc(svc.putClient))
		mux.Del("/clients/:username", http.HandlerFunc(svc.deleteClient))

		svc.mux = mux
	}

	func (svc httpAdminService) putClient(res http.ResponseWriter, req *http.Request) {
		username := req.URL.Query().Get(":username")
		password := req.FormValue("password")

		if password == "" {
			res.WriteHeader(400)
			res.Write([]byte("Password is required"))
			return
		}

		created, err := svc.clientStore.SaveUser(username, password)
		if err != nil {
			log.Error(err)
			res.WriteHeader(500)
			res.Write([]byte("An unknown error occurred."))
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
			res.Write([]byte("An unknown error occurred."))
			return
		}

		if deleted {
			res.WriteHeader(200)
		} else {
			res.WriteHeader(404)
			res.Write([]byte("User not found."))
		}
	}
