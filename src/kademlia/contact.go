package kademlia

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
	Done     bool
}

// NewContact returns a new instance of a Contact
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil, false}
}

func ContactFromHeader(header *Header) Contact {
	return NewContact(&(header.SrcID), IPToStr(header.SrcIP)+":"+strconv.Itoa(int(header.SrcPort)))
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Closer(otherContact *Contact) bool {
	return contact.distance.Closer(otherContact.distance)
}

// String returns a simple string representation of a Contact
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

// ContactCandidates definition
// stores an array of Contacts
type ContactCandidates struct {
	contacts []Contact
	mutex    sync.Mutex
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.mutex.Lock()
	candidates.contacts = append(candidates.contacts, contacts...)
	candidates.mutex.Unlock()
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Closer(&candidates.contacts[j])
}

func (candidates *ContactCandidates) Finish(k int) bool {
	for i := 0; i < Min(k, len(candidates.contacts)); i++ {
		if !candidates.contacts[i].Done {
			return false
		}
	}
	return true
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}

}

func (candidates *ContactCandidates) GetClosestUnDone(k int) (Contact, error) {

	candidates.mutex.Lock()
	for i := 0; i < Min(k, len(candidates.contacts)); i++ {
		if !candidates.contacts[i].Done {
			candidates.contacts[i].Done = true
			candidates.mutex.Unlock()
			return candidates.contacts[i], nil
		}
	}
	candidates.mutex.Unlock()
	return Contact{}, errors.New("Can't find contact")
}
