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

func testBucketFull(t *testing.T, port int, otherNode bool) {
	contacts := []Contact{
		NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:5000"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000025"), "localhost:5001"),
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

	if otherNode {
		go bucketFullFakeNode()
	}

	node := NewKademlia("0000000000000000000000000000000000000002", "127.0.0.1", port)
	go node.Listen("127.0.0.1", port)

	for _, c := range contacts {
		node.RoutingTable.Buckets[158].AddContact(c)
	}

	if otherNode {
		time.Sleep(2 * time.Second)
	} else {
		time.Sleep(10 * time.Second)
	}

	fmt.Println("Contacts added")

	if node.RoutingTable.Buckets[158].Len() > 20 {
		t.Fatal("Bucket over capacity")
	}

	for e := node.RoutingTable.Buckets[158].list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if otherNode {
			if nodeID.Equals(contacts[len(contacts)-1].ID) {
				t.Error("21st contact in the bucket")
				t.Fail()
			}
		} else {
			if nodeID.Equals(contacts[0].ID) {
				t.Error("1st contact in the bucket")
				t.Fail()
			}
		}
	}
}

func TestBucketFull1(t *testing.T) {
	testBucketFull(t, 5002, false)
	testBucketFull(t, 5001, true)
}
