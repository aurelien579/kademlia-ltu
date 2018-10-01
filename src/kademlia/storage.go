package kademlia

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Storage struct {
	Root          string
	kademlia      *Kademlia
	filenameTimer []Element2
}

type Element2 struct {
	filename string
	timer    time.Timer
}

func NewStorage(root string) Storage {
	root = "data/" + root
	os.Mkdir(root, 0777)

	return Storage{
		Root: root,
	}
}

func (storage *Storage) deleteFile(filename string) {
	fmt.Println("deleteFile")
	err := os.Remove(storage.getPath(filename))
	if err != nil {
		fmt.Print("error deleting file", err)
	}
}

func (storage *Storage) getPath(filename string) string {
	return storage.Root + "/" + strings.ToLower(filename)
}

func (storage *Storage) Exists(filename string) bool {
	_, err := os.Stat(storage.getPath(filename))
	return !os.IsNotExist(err)
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

	timer2 := time.NewTimer(10 * time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Republishing")
		storage.kademlia.Store(data)
	}()

	var exist = false

	for i := 0; i < len(storage.filenameTimer); i++ {
		if storage.filenameTimer[i].filename == filename {
			storage.filenameTimer[i].timer.Stop()
			storage.filenameTimer[i].timer = *time.AfterFunc(20*time.Second, func() {
				storage.deleteFile(filename)
			})
			exist = true
		}
	}

	if !exist {
		elem := Element2{filename, *time.AfterFunc(20*time.Second, func() {
			storage.deleteFile(filename)
		})}
		storage.filenameTimer = append(storage.filenameTimer, elem)

	}

}
