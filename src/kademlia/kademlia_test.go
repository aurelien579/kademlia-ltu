package kademlia

import (
	"log"
	"testing"
	"time"
)

func otherNode(myId string, myPort int, contact *Contact) {
	node := NewKademlia(myId, "127.0.0.1", myPort)

	go node.Listen("127.0.0.1", myPort)

	time.Sleep(100 * time.Microsecond)

	if contact != nil {
		/*node.RoutingTable.AddContact(*contact)
		node.Network.SendPingMessage(contact)*/
		node.Bootstrap(*contact)
	}
}

func TestKademliaBootstrap(t *testing.T) {
	contacts := map[int]Contact{
		//8001: NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1:8001"),
		8002: NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1:8002"),
		/*8003: NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1:8003"),
		8004: NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "127.0.0.1:8004"),
		8005: NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "127.0.0.1:8005"),
		8006: NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "127.0.0.1:8006"),
		8007: NewContact(NewKademliaID("0000000000000000000000000000000000000007"), "127.0.0.1:8007"),
		8008: NewContact(NewKademliaID("0000000000000000000000000000000000000008"), "127.0.0.1:8008"),
		8009: NewContact(NewKademliaID("0000000000000000000000000000000000000009"), "127.0.0.1:8009"),
		8010: NewContact(NewKademliaID("000000000000000000000000000000000000000A"), "127.0.0.1:8010"),
		8011: NewContact(NewKademliaID("000000000000000000000000000000000000000B"), "127.0.0.1:8011"),
		8012: NewContact(NewKademliaID("000000000000000000000000000000000000000C"), "127.0.0.1:8012"),
		8013: NewContact(NewKademliaID("000000000000000000000000000000000000000D"), "127.0.0.1:8013"),
		8014: NewContact(NewKademliaID("000000000000000000000000000000000000000E"), "127.0.0.1:8014"),
		8015: NewContact(NewKademliaID("000000000000000000000000000000000000000F"), "127.0.0.1:8015"),
		8016: NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "127.0.0.1:8016"),
		8017: NewContact(NewKademliaID("0000000000000000000000000000000000000011"), "127.0.0.1:8017"),
		8018: NewContact(NewKademliaID("0000000000000000000000000000000000000012"), "127.0.0.1:8018"),*/
	}

	go otherNode("0000000000000000000000000000000000000010", 8000, nil)

	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "127.0.0.1:8000")
	for k, v := range contacts {
		go otherNode(v.ID.String(), k, &contact)
	}

	time.Sleep(100 * time.Microsecond)

	node := NewKademlia("0000000000000000000000000000000000000014", "127.0.0.1", 8020)

	go node.Listen("127.0.0.1", 8020)

	time.Sleep(5 * time.Second)

	node.Bootstrap(NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "127.0.0.1:8000"))
}

func TestKademliaPin(t *testing.T) {
	contacts := map[int]Contact{
		8001: NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1:8001"),
		/*8002: NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1:8002"),
		8003: NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1:8003"),
		8004: NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "127.0.0.1:8004"),
		8005: NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "127.0.0.1:8005"),
		8006: NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "127.0.0.1:8006"),
		8007: NewContact(NewKademliaID("0000000000000000000000000000000000000007"), "127.0.0.1:8007"),
		8008: NewContact(NewKademliaID("0000000000000000000000000000000000000008"), "127.0.0.1:8008"),
		8009: NewContact(NewKademliaID("0000000000000000000000000000000000000009"), "127.0.0.1:8009"),
		8010: NewContact(NewKademliaID("000000000000000000000000000000000000000A"), "127.0.0.1:8010"),
		8011: NewContact(NewKademliaID("000000000000000000000000000000000000000B"), "127.0.0.1:8011"),
		8012: NewContact(NewKademliaID("000000000000000000000000000000000000000C"), "127.0.0.1:8012"),
		8013: NewContact(NewKademliaID("000000000000000000000000000000000000000D"), "127.0.0.1:8013"),
		8014: NewContact(NewKademliaID("000000000000000000000000000000000000000E"), "127.0.0.1:8014"),
		8015: NewContact(NewKademliaID("000000000000000000000000000000000000000F"), "127.0.0.1:8015"),
		8016: NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "127.0.0.1:8016"),
		8017: NewContact(NewKademliaID("0000000000000000000000000000000000000011"), "127.0.0.1:8017"),
		8018: NewContact(NewKademliaID("0000000000000000000000000000000000000012"), "127.0.0.1:8018"),*/
	}

	node := NewKademlia("0000000000000000000000000000000000000014", "127.0.0.1", 8020)
	for k, v := range contacts {
		go otherNode(v.ID.String(), k, nil)
		node.RoutingTable.AddContact(v)
	}

	time.Sleep(2 * time.Second)

	go node.Listen("127.0.0.1", 8020)

	time.Sleep(100 * time.Microsecond)

	data := []byte("TEST")
	key := NewKademliaID(HashBytes(data))

	log.Printf("Key: %s\n", key.String())

	contact := contacts[8001]
	node.Network.SendPinMessage(&contact, key, data)
	node.Network.SendUnpinMessage(&contact, key)

	time.Sleep(60 * time.Second)
}
