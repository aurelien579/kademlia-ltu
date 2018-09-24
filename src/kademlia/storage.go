package kademlia

import "os"

type Storage struct {
	Root string
}

func NewStorage(root string) Storage {
	return Storage{
		Root: root,
	}
}

func (storage *Storage) getPath(filename string) string {
	return storage.Root + "/" + filename
}

func (storage *Storage) Exists(filename string) bool {
	if _, err := os.Stat(storage.getPath(filename)); !os.IsNotExist(err) {
		return true
	}

	return false
}
