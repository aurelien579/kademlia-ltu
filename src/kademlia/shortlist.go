package kademlia

import (
	"container/list"
	"sync"
)

type element struct {
	contact Contact
	done    bool
}

type shortlist struct {
	contacts *list.List
	maxSize  int
	mutex    sync.Mutex
}

func NewShortList(maxSize int) *shortlist {
	return &shortlist{
		contacts: list.New(),
		maxSize:  maxSize,
	}
}

func (l *shortlist) AddOrdered(c *Contact) {
	// Todo: abandonner si c est après maxSize éléments
}

func (l *shortlist) AddAllOrdered(contacts []Contact) {

}

/* Renvoie true si les k plus proches ont été fait */
func (l *shortlist) ClosestDone(k int) {

}

func (l *shortlist) GetClosestUndone() {

}
