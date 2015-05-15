	<<#-->>
	package main

	import (
		"fmt"
		"net/http"
		"sync"

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

		return httpAdminService {
			clientStore: clientStore,
			port: port,
		}, nil
	}

The admin service will run as a REST service on a different port to the "main"
file service.

	type httpAdminService struct {
		clientStore ClientStore
		port int
	}

	func (svc httpAdminService) Start(wait *sync.WaitGroup) {
		log.Info("Starting admin service on port", svc.port)
		http.ListenAndServe(fmt.Sprintf(":%d", svc.port), svc)
		log.Info("Stopping admin service.")
		wait.Done()
	}

	func (svc httpAdminService) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	}
