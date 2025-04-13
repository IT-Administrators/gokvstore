package gokvstore

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
)

// This storer interface represents the methods to manage the key value store.
// Generic interface to use with any datatype.
//
// Use [K int64 | float64 | string, V any] to restrict keys to types int, float and string.
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
	Data map[K]V
}

// Create key value store.
func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		Data: make(map[K]V),
	}
}

// Check if key is present in store.
func (s *KVStore[K, V]) hasKey(key K) bool {
	_, ok := s.Data[key]
	return ok
}

// Insert key to store.
func (s *KVStore[K, V]) Put(key K, value V) error {
	// This is a write method. It needs a write mutex to prevent changes while inserting values.
	s.mu.Lock()
	// Unlock store.
	defer s.mu.Unlock()
	// Insert into store.
	s.Data[key] = value
	return nil
}

// Retrieve value for specified key.
func (s *KVStore[K, V]) Get(key K) (V, error) {
	// This is read method so it needs a read mutex to prevent chagnes while reading.
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if key is in store.
	value, exists := s.Data[key]
	if !exists {
		return value, fmt.Errorf("the key (%v) does not exist", key)
	}
	return value, nil
}

// Update specified key with specified value.
func (s *KVStore[K, V]) Update(key K, value V) error {
	// This is a write method. It needs a write mutex to prevent changes while inserting values.
	// Lock store for writing.
	s.mu.Lock()
	// Unlock store.
	defer s.mu.Unlock()

	// Check if key exists.
	if !s.hasKey(key) {
		return fmt.Errorf("the key (%v) does not exist", key)
	}
	// Insert into store if not exists.
	s.Data[key] = value
	return nil
}

// Delete specified key. This will remove key value pair from store.
func (s *KVStore[K, V]) Delete(key K) (V, error) {
	// This is a read method so it needs a read mutex to prevent changes while reading.
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.Data[key]
	if !exists {
		return value, fmt.Errorf("the key (%v) does not exist", key)
	}
	delete(s.Data, key)
	// Show old value.
	return value, nil
}

// Removes all keys from store.
func (s *KVStore[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.Data {
		delete(s.Data, key)
	}
}

// Print the current store.
func (s *KVStore[K, V]) Print() {
	for k, d := range s.Data {
		fmt.Printf("key: %v value: %v\n", k, d)
	}
}

// Load store from file into current store.
//
// This will overwrite all keys with the ones from file.
// If they were deleted before, they will be recreated.
func (s *KVStore[K, V]) Load(file string) error {

	// Apply lock to stop writing while importin.
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Open file.
	loadFrom, err := os.Open(file)

	if err != nil {
		return fmt.Errorf("empty key/value store!, Error: %v", err)
	}
	defer loadFrom.Close()

	// Create new decoder and decode.
	decoder := gob.NewDecoder(loadFrom)
	decoder.Decode(&s)

	return nil
}

// Save keys to file.
func (s *KVStore[K, V]) Save(file string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove file if already exists.
	err := os.Remove(file)
	if err != nil {
		fmt.Println(err)
	}

	// Create file if not exists.
	saveTo, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("cannot create file %v with error %v", file, err)
	}
	defer saveTo.Close()
	// Create new encoder and encode.
	encoder := gob.NewEncoder(saveTo)
	err = encoder.Encode(&s)
	if err != nil {
		return fmt.Errorf("cannot save to file %v with error %v", file, err)
	}

	return nil
}
