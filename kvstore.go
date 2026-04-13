package gokvstore

import (
	"encoding/gob"
	"fmt"
	"hash/fnv"
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

/* Each store is split into smaller maps. And each key is assigned to exactly one map based on its hash.
// Sharding spreads keys across multiple independent maps so:
// - Different keys often hit different shards
// - Operations on different shards dont intefere with each other
// - Better througput with concurrency
// */

// Number of shards. Power of two makes masking easy.
const defaultShards = 32

type shard[V any] struct {
	data sync.Map // map[any]V via interface{}, keep types consistent at the API
}

// Key value store.
type KVStore[K comparable, V any] struct {
	shards []shard[V]
	mask   uint64
}

// NewKVStore creates a sharded sync.Map-based KV store.
func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	s := &KVStore[K, V]{
		shards: make([]shard[V], defaultShards),
		mask:   uint64(defaultShards - 1),
	}
	return s
}

// Produce a uint64 hash from the keys string representation.
func hashKey[K comparable](key K) uint64 {
	// Create a new hash function.
	h := fnv.New64a()
	// Convert key into string and wirte to hash function.
	fmt.Fprintf(h, "%v", key) // simple, generic; can be optimized for specific K
	return h.Sum64()
}

/* mask = 32 -1 = 31 = 0b11111
h & 0b11111 extracts lowest number of bits of the hash (number between 0 - 31)
The number is the shard index.
Only works because shards is a power of two.*/

// Get associated shard for key.
func (s *KVStore[K, V]) getShard(key K) *shard[V] {
	// Get hash from previous function.
	h := hashKey(key)
	// Return pointer to selected shard.
	return &s.shards[h&s.mask]
}

// Insert key to store.
func (s *KVStore[K, V]) Put(key K, value V) error {
	// Find the shard. Than only operate on that shard.
	sh := s.getShard(key)
	sh.data.Store(key, value)
	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	sh := s.getShard(key)
	val, ok := sh.data.Load(key)
	if !ok {
		var zero V
		return zero, fmt.Errorf("key %v not found", key)
	}
	return val.(V), nil
}

// Update specified key with specified value.
func (s *KVStore[K, V]) Update(key K, value V) error {
	sh := s.getShard(key)
	_, ok := sh.data.Load(key)
	if !ok {
		return fmt.Errorf("key %v not found", key)
	}
	sh.data.Store(key, value)
	return nil
}

func (s *KVStore[K, V]) Delete(key K) (V, error) {
	sh := s.getShard(key)
	val, ok := sh.data.LoadAndDelete(key)
	if !ok {
		var zero V
		return zero, fmt.Errorf("key %v not found", key)
	}
	return val.(V), nil
}

// Removes all keys from store.
// Clear removes all keys from all shards.
func (s *KVStore[K, V]) Clear() {
	for i := range s.shards {
		sh := &s.shards[i]
		sh.data.Range(func(k, _ any) bool {
			sh.data.Delete(k)
			return true
		})
	}
}

// Print the current store.
// Print is mainly for debugging.
func (s *KVStore[K, V]) Print() {
	for i := range s.shards {
		sh := &s.shards[i]
		sh.data.Range(func(k, v any) bool {
			fmt.Printf("key: %v value: %v\n", k, v)
			return true
		})
	}
}

// Load store from file into current store.
//
// This will overwrite all keys with the ones from file.
// If they were deleted before, they will be recreated.
// Load replaces current contents with those from file.
func (s *KVStore[K, V]) Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("cannot open file %v: %w", file, err)
	}
	defer f.Close()

	tmp := make(map[K]V)
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&tmp); err != nil {
		return fmt.Errorf("cannot decode data from %v: %w", file, err)
	}

	// Clear and repopulate shards.
	s.Clear()
	for k, v := range tmp {
		_ = s.Put(k, v)
	}
	return nil
}

// Save keys to file.
// Save serializes all shards into a single map[K]V and writes it to file.
func (s *KVStore[K, V]) Save(file string) error {
	// Collect into a plain map for gob.
	tmp := make(map[K]V)

	for i := range s.shards {
		sh := &s.shards[i]
		sh.data.Range(func(k, v any) bool {
			tmp[k.(K)] = v.(V)
			return true
		})
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("cannot create file %v: %w", file, err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(tmp); err != nil {
		return fmt.Errorf("cannot encode data to %v: %w", file, err)
	}
	return nil
}
