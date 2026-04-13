package gokvstore

import (
	"os"
	"sync"
	"testing"
)

var (
	testKey1   = "T1"
	testValue1 = "Test1"
	testKey2   = "T2"
	testValue2 = "Test2"
	saveFile   = "./examples/test_store.gob"
)

func Test_PutAndGet(t *testing.T) {
	kvs := NewKVStore[string, string]()

	if err := kvs.Put(testKey1, testValue1); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	val, err := kvs.Get(testKey1)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if val != testValue1 {
		t.Errorf("expected %v got %v", testValue1, val)
	}
}

func Test_Update(t *testing.T) {
	kvs := NewKVStore[string, string]()
	kvs.Put(testKey1, testValue1)

	newVal := "UpdatedValue"
	if err := kvs.Update(testKey1, newVal); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	val, _ := kvs.Get(testKey1)
	if val != newVal {
		t.Errorf("expected %v got %v", newVal, val)
	}
}

func Test_Delete(t *testing.T) {
	kvs := NewKVStore[string, string]()
	kvs.Put(testKey1, testValue1)

	_, err := kvs.Delete(testKey1)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = kvs.Get(testKey1)
	if err == nil {
		t.Errorf("expected error for deleted key")
	}
}

func Test_Clear(t *testing.T) {
	kvs := NewKVStore[string, string]()
	kvs.Put(testKey1, testValue1)
	kvs.Put(testKey2, testValue2)

	kvs.Clear()

	_, err1 := kvs.Get(testKey1)
	_, err2 := kvs.Get(testKey2)

	if err1 == nil || err2 == nil {
		t.Errorf("Clear did not remove all keys")
	}
}

func Test_SaveAndLoad(t *testing.T) {
	kvs := NewKVStore[string, string]()
	kvs.Put(testKey1, testValue1)
	kvs.Put(testKey2, testValue2)

	if err := kvs.Save(saveFile); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Modify store to ensure Load overwrites it
	kvs.Update(testKey2, "ModifiedValue")

	if err := kvs.Load(saveFile); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	val, _ := kvs.Get(testKey2)
	if val != testValue2 {
		t.Errorf("expected %v got %v", testValue2, val)
	}

	os.Remove(saveFile)
}

func Test_PutConcurrent(t *testing.T) {
	kvs := NewKVStore[int, int]()
	const count = 100000

	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(n int) {
			defer wg.Done()
			kvs.Put(n, n)
		}(i)
	}

	wg.Wait()

	for i := 0; i < count; i++ {
		val, err := kvs.Get(i)
		if err != nil {
			t.Fatalf("missing key %v: %v", i, err)
		}
		if val != i {
			t.Errorf("expected %v got %v", i, val)
		}
	}
}

// Fuzz Test
func Fuzz_PutGet(f *testing.F) {
	kvs := NewKVStore[string, string]()

	f.Add("hello", "world")
	f.Add("foo", "bar")

	f.Fuzz(func(t *testing.T, key, value string) {
		err := kvs.Put(key, value)
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		got, err := kvs.Get(key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if got != value {
			t.Fatalf("expected %v got %v", value, got)
		}
	})
}

func Benchmark_Put(b *testing.B) {
	kvs := NewKVStore[int, int]()
	for i := 0; i < b.N; i++ {
		kvs.Put(i, i)
	}
}

func Benchmark_Get(b *testing.B) {
	kvs := NewKVStore[int, int]()
	for i := 0; i < b.N; i++ {
		kvs.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kvs.Get(i)
	}
}

func Benchmark_ConcurrentPut(b *testing.B) {
	kvs := NewKVStore[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			kvs.Put(i, i)
			i++
		}
	})
}

// Shard Distribution Test
func Test_ShardDistribution(t *testing.T) {
	kvs := NewKVStore[int, int]()

	const keys = 10000
	counts := make([]int, defaultShards)

	for i := 0; i < keys; i++ {
		sh := kvs.getShard(i)
		idx := int(hashKey(i) & kvs.mask)
		if &kvs.shards[idx] != sh {
			t.Fatalf("shard mismatch for key %v", i)
		}
		counts[idx]++
	}

	// Check distribution is reasonably even
	for i, c := range counts {
		if c == 0 {
			t.Errorf("shard %d received no keys", i)
		}
	}
}
