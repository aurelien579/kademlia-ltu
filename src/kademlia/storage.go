package kademlia

import (
	"fmt"
	"io/ioutil"
	"os"
)

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

func (storage *Storage) Read(filename string) []byte {
	bytes, err := ioutil.ReadFile(storage.getPath(filename))

	if err == nil {
		fmt.Println("ERROR: ", err)
		return nil
	}

	return bytes
}
