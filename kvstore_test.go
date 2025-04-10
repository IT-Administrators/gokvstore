package gokvstore

import (
	"testing"
)

var testKey = "T1"
var testValue = "Test1"
var kvs = NewKVStore[string, any]()

func Test_Put(t *testing.T) {

	kvs.Put(testKey, testValue)
	kvs.Put("T2", "Test2")
	ok := kvs.data[testKey]
	if ok == nil {
		t.Errorf("could not insert value into store.")
	}
}

func Test_Get(t *testing.T) {

	if ok := kvs.data[testKey]; ok == nil {
		t.Errorf("key %v does not exist", testKey)
	}
}

func Test_Update(t *testing.T) {

	val := "This value was changed."
	kvs.Update(testKey, val)
	if val != kvs.data[testKey] {
		t.Errorf("got %v wanted %v", kvs.data[testKey], val)
	}
}

func Test_Delete(t *testing.T) {
	kvs.Delete(testKey)
	if ok := kvs.data[testKey]; ok != nil {
		t.Errorf("could not remove %v from store.", testKey)
	}
}
