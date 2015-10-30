	<<#-->>
	package main

	import (
		"time"
	)

Each Blob node keeps track of all known sibling nodes, including when they were
introduced, whether the sibling is estranged, and the last time they've been in
touch.

	type Siblings interface {
		Add(uri string) (SiblingStatus, error)
		All() []SiblingStatus
		Status(uri string) (SiblingStatus, bool)
	}

	type SiblingStatus struct {
		URI string
		Introduced time.Time
		LastContact time.Time
		Estranged bool
	}

	func NewSiblingStore() Siblings {
		return memSiblingStore {
			make(map[string]SiblingStatus),
		}
	}

Sibling information is kept in memory, since it is transient and always changing.

	type memSiblingStore struct {
		siblings map[string]SiblingStatus
	}

Siblings are considered "introduced" when they are first added. Adding a sibling
will immediately try to connect to the sibling upon introduction; if this
connection fails, the sibling will not be added, and an error will be returned.

	func (store memSiblingStore) Add(uri string) (SiblingStatus, error) {
		// TODO: check connection and add to bus.

		status := SiblingStatus {
			URI: uri,
			Introduced: time.Now(),
			LastContact: time.Now(),
			Estranged: false,
		}

		store.siblings[uri] = status
		return status, nil
	}

	func (store memSiblingStore) All() []SiblingStatus {
		ss := make([]SiblingStatus, 0, len(store.siblings))
		for _, status := range(store.siblings) {
			ss = append(ss, status)
		}
		return ss
	}

	func (store memSiblingStore) Status(uri string) (SiblingStatus, bool) {
		status, ok := store.siblings[uri]
		return status, ok
	}