package kademlia

import (
	"fmt"
	"testing"
	"time"
)

func TestBucketCalcDistance(t *testing.T) {
	input := []Contact{
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000007"), "localhost:8001"),
	}

	expected := map[string]*KademliaID{
		"ffffffff00000000000000000000000000000000": NewKademliaID("FFFFFFFF00000000000000000000000000000003"),
		"0000000000000000000000000000000000000001": NewKademliaID("0000000000000000000000000000000000000002"),
		"0000000000000000000000000000000000000003": NewKademliaID("0000000000000000000000000000000000000000"),
		"0000000000000000000000000000000000000005": NewKademliaID("0000000000000000000000000000000000000006"),
		"0000000000000000000000000000000000000007": NewKademliaID("0000000000000000000000000000000000000004"),
	}

	b := newBucket()

	for _, c := range input {
		b.AddContact(c)
	}

	results := b.GetContactAndCalcDistance(NewKademliaID("0000000000000000000000000000000000000003"))

	for i := 0; i < len(results); i++ {
		if *results[i].distance != *expected[results[i].ID.String()] {
			t.Errorf("Wrong distance %v expected %v\n", results[i].distance, expected[results[i].ID.String()])
			t.Fail()
		}
	}
}

func bucketFullFakeNode() {
	node := NewKademlia("0000000000000000000000000000000000000001", "127.0.0.1", 5000)

	node.Listen("127.0.0.1", 5000)
}

func TestBucketFull(t *testing.T) {
	contacts := []Contact{
		NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:5000"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000007"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000008"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000009"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000011"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000012"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000013"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000014"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000015"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000016"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000017"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000018"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000019"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000020"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000021"), "localhost:8001"),
	}

	go bucketFullFakeNode()

	node := NewKademlia("0000000000000000000000000000000000000002", "127.0.0.1", 5001)
	go node.Listen("127.0.0.1", 5001)

	time.Sleep(100 * time.Microsecond)

	for _, c := range contacts {
		node.RoutingTable.Buckets[158].AddContact(c)
	}

	time.Sleep(100 * time.Microsecond)

	fmt.Println("Contacts added")

	//index := node.RoutingTable.getBucketIndex(contacts[0].ID)
	fmt.Println(node.RoutingTable.Buckets[158].Len())

	for e := node.RoutingTable.Buckets[158].list.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(Contact).ID)

		nodeID := e.Value.(Contact).ID
		if nodeID.Equals(contacts[0].ID) {
			t.Fail()
		}
	}
}
