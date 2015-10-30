	<<#-->>
	package main

	import (
		"bytes"
		"crypto/sha256"

		"github.com/tokenshift/log"
	)

The client store stores client/user information, hashing user passwords that
can then be used to validate the provided credentials.

	type ClientStore interface {
		GetUsers() ([]string, error)
		SaveUser(username, password string) (bool, error)
		DeleteUser(username string) (bool, error)
		VerifyUser(username, password string) (bool, error)
	}

	func NewClientStore() (ClientStore, error) {
		return &boltClientStore {
			hashes: make(map[string][]byte),
		}, nil
	}

The standard implementation stores (hashed) user credentials in a local Bolt DB.
(Or it will, this is a stub.)

	type boltClientStore struct {
		hashes map[string][]byte
	}

	func (store *boltClientStore) GetUsers() ([]string, error) {
		users := make([]string, 0, len(store.hashes))

		for user := range store.hashes {
			users = append(users, user)
		}

		return users, nil
	}

	func (store *boltClientStore) SaveUser(username, password string) (bool, error) {
		if _, ok := store.hashes[username]; ok {
			log.Debug("Updating client", username)
			store.hashes[username] = hashPassword(password)
			return false, nil
		} else {
			log.Debug("Adding client", username)
			store.hashes[username] = hashPassword(password)
			store.hashes[username] = hashPassword(password)
			return true, nil
		}
	}

	func (store *boltClientStore) DeleteUser(username string) (bool, error) {
		if _, ok := store.hashes[username]; ok {
			log.Debug("Removing client", username)
			delete(store.hashes, username)
			return true, nil
		} else {
			return false, nil
		}
	}

	func (store *boltClientStore) VerifyUser(username, password string) (bool, error) {
		if actualHash, ok := store.hashes[username]; ok {
			return bytes.Equal(hashPassword(password), actualHash), nil
		} else {
			return false, nil
		}
	}

	func hashPassword(password string) []byte {
		hash := sha256.Sum256([]byte(password))
		return hash[:]
	}
