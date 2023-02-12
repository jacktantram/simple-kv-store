package kvstore

import (
	"errors"
	"strings"
)

// KVStore represents the base store for a basic in memory kv store.
// It holds information about the root kv in memory and keeps track of the head transaction.
// The head is the most recently created transaction.
type KVStore struct {
	dataStore dataStore
	head      *Transaction
}

// NewStore instantiates a new kv store.
func NewStore() KVStore {
	return KVStore{dataStore: make(dataStore)}
}

// Get is responsible for fetching a value from the store
// that matches the key. If the value is not found ok will be false.
// If the store currently has a running transaction it will fetch it from the value from the latest transaction.
func (s KVStore) Get(key string) (val string, ok bool) {
	if s.head != nil {
		return s.head.dataStore.Get(key)
	}
	return s.dataStore.Get(key)
}

// Set will set a value in the store with the specified key and value.
// The store does not allow empty key,values to be stored. Keys are case-insensitive.
// If the store currently has a running transaction it will set the value in the latest transaction and not the
// base store.
func (s *KVStore) Set(key, value string) error {
	if s.head != nil {
		return s.head.dataStore.Set(key, value)
	}
	return s.dataStore.Set(key, value)
}

// Delete will delete a value in the store with the specified key.
// If the store currently has a running transaction it will delete the value in the latest transaction and not the
// base store.
// Ok is returned indicating if deletion was successful. If the key did not exist 0 will be returned, 1 for success.
func (s *KVStore) Delete(key string) (ok bool) {
	if s.head != nil {
		return s.head.dataStore.Delete(key)
	}
	return s.dataStore.Delete(key)
}

// dataStore represents the data structure of how information
// is stored in the store.
type dataStore map[string]string

// Get returns a value from the data store.
// It returns a boolean if the value does not exist.
func (d dataStore) Get(key string) (val string, ok bool) {
	val, ok = d[key]
	return
}

// Set writes a value given a key to the data store.
// If the key or value are empty an error is returned
func (d dataStore) Set(key, value string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("key cannot be empty")
	}
	if strings.TrimSpace(value) == "" {
		return errors.New("value cannot be empty")
	}
	d[key] = value
	return nil
}

// Delete removes a value given a key from the data store.
// If the value does not exist False is returned
func (d dataStore) Delete(key string) bool {
	if _, ok := d.Get(key); !ok {
		return false
	}
	delete(d, key)
	return true
}

// Copy copies the contents of the internal structure to a new dataStore object.
func (d dataStore) Copy() dataStore {
	targetMap := make(dataStore)
	for key, value := range d {
		targetMap[key] = value
	}
	return targetMap
}
