package gokvstore

import (
	"fmt"
	"sync"
)

// This storer interface represents the methods to manage the key value store.
// Generic interface to use with any datatype.
type Storer[K comparable, V any] interface {
	// Insert to storage.
	Put(K, V) error
	// Get value for specified key.
	Get(K) (V, error)
	// Update value for specified key.
	Update(K, V) error
	// Delete key from storage.
	Delete(K) (V, error)
}

// Key value store.
type KVStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// Create key value store.
func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		data: make(map[K]V),
	}
}

// Check if key is present in store.
func (s *KVStore[K, V]) hasKey(key K) bool {
	_, ok := s.data[key]
	return ok
}

// Insert key to store.
// Implementing the Storer interface for the KVStore struct.
// This is a write method. It needs a write mutex to prevent changes while inserting values.
func (s *KVStore[K, V]) Put(key K, value V) error {
	// Lock store for writing.
	s.mu.Lock()
	// Unlock store.
	defer s.mu.Unlock()
	// Insert into store.
	s.data[key] = value
	return nil
}

// Retrieve value for specified key.
// Implementing the Storer interface for the KVStore struct.
// This is read method so it needs a read mutex to prevent chagnes while reading.
func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	if !exists {
		return value, fmt.Errorf("the key (%v) does not exist", key)
	}
	return value, nil
}

// Update specified key with specified value.
// Implementing the Storer interface for the KVStore struct.
// This is a write method. It needs a write mutex to prevent changes while inserting values.
func (s *KVStore[K, V]) Update(key K, value V) error {
	// Lock store for writing.
	s.mu.Lock()
	// Unlock store.
	defer s.mu.Unlock()

	// Check if key exists.
	if !s.hasKey(key) {
		return fmt.Errorf("the key (%v) does not exist", key)
	}
	// Insert into store if not exists.
	s.data[key] = value
	return nil
}

// Delete specified key. This will remove key value pair from store.
// Implementing the Storer interface for the KVStore struct.
// This is a read method so it needs a read mutex to prevent changes while reading.
func (s *KVStore[K, V]) Delete(key K) (V, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.data[key]
	if !exists {
		return value, fmt.Errorf("the key (%v) does not exist", key)
	}
	delete(s.data, key)
	// Show old value.
	return value, nil
}

// Print the current store.
func (s *KVStore[K, V]) Print() {
	for k, d := range s.data {
		fmt.Printf("key: %v value: %v\n", k, d)
	}
}
