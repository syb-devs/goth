package kv

import (
	"sync"
)

// Store is just a simple, concurrency-safe, in-memory key/value store
type Store struct {
	sync.RWMutex
	values map[interface{}]interface{}
}

// New allocates and returns Store object
func New() *Store {
	return &Store{
		values: make(map[interface{}]interface{}),
	}
}

// Set sets a key in the store with the given value
func (s *Store) Set(key, value interface{}) {
	s.Lock()
	s.values[key] = value
	s.Unlock()
}

// Get returns the value for a given key
func (s *Store) Get(key interface{}) interface{} {
	s.RLock()
	defer s.RUnlock()

	return s.values[key]
}
