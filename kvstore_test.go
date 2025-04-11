package gokvstore

import (
	"testing"
)

var testKey1 = "T1"
var testValue1 = "Test1"
var testKey2 = "T2"
var testValue2 = "Test2"
var saveFile = "./examples/store.gob"
var kvs = NewKVStore[string, any]()

func Test_Put(t *testing.T) {

	kvs.Put(testKey1, testValue1)
	kvs.Put(testKey2, testValue2)
	ok := kvs.Data[testKey1]
	if ok == nil {
		t.Errorf("could not insert value into store.")
	}
}

func Test_Get(t *testing.T) {

	if ok := kvs.Data[testKey1]; ok == nil {
		t.Errorf("key %v does not exist", testKey1)
	}
}

func Test_Update(t *testing.T) {

	val := "This value was changed."
	kvs.Update(testKey1, val)
	if val != kvs.Data[testKey1] {
		t.Errorf("got %v wanted %v", kvs.Data[testKey1], val)
	}
}

func Test_Delete(t *testing.T) {
	kvs.Delete(testKey1)
	if ok := kvs.Data[testKey1]; ok != nil {
		t.Errorf("could not remove %v from store.", testKey1)
	}
}

func Test_SaveAndLoad(t *testing.T) {
	kvs.Save(saveFile)
	// Change key to test load.
	kvs.Update(testKey2, "This key was changed")
	kvs.Load(saveFile)
	if val, _ := kvs.Get("T2"); val != "Test2" {
		t.Errorf("got: %v expected: %v", val, testValue2)
	}
}

func Test_Clear(t *testing.T) {
	kvs.Clear()
	if ok := kvs.Data[testKey1]; ok != nil {
		t.Errorf("removal not succesffull; value: %v", kvs.Data[testKey1])
	}
}
