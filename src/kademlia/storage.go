package kademlia

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const REPUBLISH_TIME = 30

type Storage struct {
	Root          string
	kademlia      *Kademlia
	filenameTimer []Element2
	mutex         sync.Mutex
}

type Element2 struct {
	filename       string
	timerRepublish *time.Timer
	timerDelete    *time.Timer
}

func NewStorage(root string) Storage {
	root = "data/" + root
	os.Mkdir(root, 0777)

	return Storage{
		Root: root,
	}
}

func (storage *Storage) deleteFile(filename string) {
	log.Println("deleteFile")
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
		log.Println("ERROR: ", err)
		return nil
	}

	return bytes
}

func (storage *Storage) Store(filename string, data []byte) {

	log.Println("Store: ", filename)

	ioutil.WriteFile(storage.getPath(filename), data, 0644)

	storage.mutex.Lock()

	var exist = storage.Exist(filename)

	if exist {


		log.Println("the file exist: ",filename)

		for i := 0; i < len(storage.filenameTimer); i++ {
			if storage.filenameTimer[i].filename == filename {
				storage.filenameTimer[i].timerRepublish.Reset(1 * REPUBLISH_TIME * time.Second)
				storage.filenameTimer[i].timerDelete.Reset(2 * REPUBLISH_TIME * time.Second)
			}
		}
	} else {


		log.Println("the file doesn't exist: ",filename)

		
		timerRepublish := time.AfterFunc(1*REPUBLISH_TIME*time.Second, func() {
			storage.kademlia.Store(data)
		})

		timerDelete := time.AfterFunc(2*REPUBLISH_TIME*time.Second, func() {
			storage.deleteFile(filename)
			storage.DeleteElement(filename)
		})

		elem := Element2{filename, timerRepublish, timerDelete}
		storage.filenameTimer = append(storage.filenameTimer, elem)
	}

	storage.mutex.Unlock()

}

func (storage *Storage) Exist(filename string) bool {

	for i := 0; i < len(storage.filenameTimer); i++ {
		if storage.filenameTimer[i].filename == filename {

			log.Println("Store: ", filename)
			return true
		}
	}
	return false

}

func (storage *Storage) DeleteElement(filename string) {
	for i := 0; i < len(storage.filenameTimer); i++ {
		if storage.filenameTimer[i].filename == filename {
			storage.filenameTimer[i].timerRepublish.Stop()
			storage.filenameTimer[i].timerDelete.Stop()
			storage.filenameTimer = append(storage.filenameTimer[:i], storage.filenameTimer[i+1:]...)
		}
	}

}
