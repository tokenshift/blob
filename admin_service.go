package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/tokenshift/blob/env"
	"github.com/tokenshift/blob/log"
)

// Administrative interface for managing the Blob instance.
type AdminService interface {
	Start(*sync.WaitGroup)
}

func NewAdminService(clientStore ClientStore) (AdminService, error) {
	port, ok, err := env.GetInt("BLOB_ADMIN_SERVICE_PORT")
	if !ok {
		return nil, fmt.Errorf("Missing $BLOB_ADMIN_SERVICE_PORT")
	}
	if err != nil {
		return nil, fmt.Errorf("Invalid $BLOB_ADMIN_SERVICE_PORT")
	}

	return httpAdminService {
		clientStore: clientStore,
		port: port,
	}, nil
}

// Provides a REST interface to administrative services.
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
