package kademlia

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Storage struct {
	Root string
}

func NewStorage(root string) Storage {
	os.Mkdir(root, 0777)
	return Storage{
		Root: root,
	}
}

func (storage *Storage) getPath(filename string) string {
	return storage.Root + "/" + strings.ToLower(filename)
}

func (storage *Storage) Exists(filename string) bool {
	_, err := os.Stat(storage.getPath(filename))

	fmt.Println("File: ", storage.getPath(filename))

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (storage *Storage) Read(filename string) []byte {
	bytes, err := ioutil.ReadFile(storage.getPath(filename))

	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil
	}

	return bytes
}

func (storage *Storage) Store(filename string, data []byte) {
	ioutil.WriteFile(storage.getPath(filename), data, 0644)
}
