package main

// Stores client/user information.
type ClientStore interface {
}

func NewClientStore() (ClientStore, error) {
	return boltClientStore {
	}, nil
}

// Stores client info in a local BoltDB.
type boltClientStore struct {
}
