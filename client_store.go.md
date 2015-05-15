	<<#-->>
	package main

The client store stores client/user information, hashing user passwords that
can then be used to validate the provided credentials.

	type ClientStore interface {
	}

	func NewClientStore() (ClientStore, error) {
		return boltClientStore {
		}, nil
	}

The standard implementation stores (hashed) user credentials in a local Bolt DB.

	type boltClientStore struct {
	}
