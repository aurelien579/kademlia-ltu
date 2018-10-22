package kademlia

import (
	"testing"
)

func TestStorage(t *testing.T) {
	storage := NewStorage("TEST")
	storage.deleteFile("Unexisting file")

	storage.Store("test.txt", []byte("bonjour"))
	storage.Store("test.txt", []byte("bonjour2"))

	out := storage.Read("test.txt")
	if string(out) != "bonjour" {
		t.Fail()
	}

	storage.deleteFile("test.txt")
	if storage.Exists("test.txt") {
		t.Fail()
	}
}
