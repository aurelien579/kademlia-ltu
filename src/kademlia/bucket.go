package kademlia

import (
	"container/list"
	"fmt"
	"time"
)

//ceci est un putain de commentaire pour voir si ça marche

// bucket definition
// contains a List
type bucket struct {
	list *list.List
	node *Kademlia
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		} else {
			olderKnown := bucket.list.Back().Value.(Contact)
			fmt.Printf("OlderKnown : %v\n", olderKnown)

			channel := make(chan Header)

			bucket.node.Channels.Add(&olderKnown, MSG_PING, channel)

			bucket.node.Network.SendPingMessage(&olderKnown)

			timer2 := time.NewTimer(1 * time.Second)
			go func() {
				<-timer2.C

				bucket.node.Channels.Delete(&olderKnown, MSG_PING)
				bucket.list.Remove(bucket.list.Back())
				bucket.list.PushFront(contact)
			}()

			<-channel
			bucket.list.MoveToFront(bucket.list.Back())
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for e := bucket.list.Front(); e != nil; e = e.Next() {
		contact := e.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
