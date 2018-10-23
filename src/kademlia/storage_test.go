package kademlia

import (
	"testing"
)

func TestStorage(t *testing.T) {
	storage := NewStorage("TEST")

	storage.Store("test.txt", []byte("bonjour"))
	storage.Store("test.txt", []byte("bonjour"))

	out := storage.Read("test.txt")
	if string(out) != "bonjour" {
		t.Error("Invalid content")
	}

	storage.deleteFile("test.txt")
	if storage.Exists("test.txt") {
		t.Error("Should not exists")
	}
}
