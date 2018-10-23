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

type Storage struct {
	Root      string
	kademlia  *Kademlia
	fileInfos []FileInfo
	mutex     sync.Mutex
}

type FileInfo struct {
	filename       string
	timerRepublish *time.Timer
	timerDelete    *time.Timer
}

func NewStorage(root string) Storage {
	root = "data/" + root
	os.Mkdir("data/", 0777)
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

func (storage *Storage) Store(filename string, data []byte, pin bool) {
	log.Println("Store: ", filename)
	log.Printf("%s %v\n", storage.kademlia.RoutingTable.Me.ID.String(), storage.fileInfos)

	ioutil.WriteFile(storage.getPath(filename), data, 0644)

	storage.mutex.Lock()

	var exist = storage.ExistsElement(filename)

	if exist {
		log.Println("the file exist: ", filename)

		fileInfo := storage.findFileInfo(filename)
		storage.updateFileInfo(fileInfo, filename, pin)
	} else {
		log.Println("the file doesn't exist: ", filename)

		fileInfo := storage.createFileInfo(filename, data, pin)

		log.Printf("Adding file: %s\n", fileInfo.filename)
		storage.fileInfos = append(storage.fileInfos, fileInfo)
		log.Printf("%s %v\n", storage.kademlia.RoutingTable.Me.ID.String(), storage.fileInfos)
	}

	storage.mutex.Unlock()
}

func (storage *Storage) createFileInfo(filename string, data []byte, pin bool) FileInfo {
	timerRepublish := storage.createTimerRepublish(data)
	var timerDelete *time.Timer = nil

	if !pin {
		timerDelete = storage.createTimerDelete(filename)
	}

	return FileInfo{filename, timerRepublish, timerDelete}
}

func (storage *Storage) updateFileInfo(fileInfo *FileInfo, filename string, pin bool) {
	fileInfo.timerRepublish.Reset(1 * REPUBLISH_TIME * time.Second)

	if pin == true {
		fileInfo.timerDelete = nil
	} else {
		if fileInfo.timerDelete == nil {
			fileInfo.timerDelete = storage.createTimerDelete(filename)
		} else {
			fileInfo.timerDelete.Reset(2 * REPUBLISH_TIME * time.Second)
		}
	}
}

func (storage *Storage) createTimerDelete(filename string) *time.Timer {
	return time.AfterFunc(2*REPUBLISH_TIME*time.Second, func() {
		log.Printf("DELETING FILE")
		storage.deleteFile(filename)
		storage.DeleteElement(filename)
	})
}

func (storage *Storage) createTimerRepublish(data []byte) *time.Timer {
	return time.AfterFunc(1*REPUBLISH_TIME*time.Second, func() {
		log.Printf("Republish id %s\n", storage.kademlia.RoutingTable.Me.String())
		storage.kademlia.Store(data)
	})
}

func (storage *Storage) findFileInfo(filename string) *FileInfo {
	for i := 0; i < len(storage.fileInfos); i++ {
		if storage.fileInfos[i].filename == filename {
			return &storage.fileInfos[i]
		}
	}

	return nil
}

func (storage *Storage) ExistsElement(filename string) bool {
	for i := 0; i < len(storage.fileInfos); i++ {
		if storage.fileInfos[i].filename == filename {
			return true
		}
	}

	return false
}

func (storage *Storage) DeleteElement(filename string) {
	log.Printf("%s is deleting %s\n", storage.kademlia.RoutingTable.Me.ID.String(), filename)

	for i := 0; i < len(storage.fileInfos); i++ {
		if storage.fileInfos[i].filename == filename {
			storage.fileInfos[i].timerRepublish.Stop()
			storage.fileInfos[i].timerDelete.Stop()
			storage.fileInfos = append(storage.fileInfos[:i], storage.fileInfos[i+1:]...)
		}
	}
}

func (storage *Storage) Unpin(filename string) {
	fileInfo := storage.findFileInfo(filename)
	if fileInfo.timerDelete == nil {
		fileInfo.timerDelete = storage.createTimerDelete(filename)
	}
}
