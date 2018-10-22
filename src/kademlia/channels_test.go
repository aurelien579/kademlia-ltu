package kademlia

import (
	"testing"
)

func TestChannels(t *testing.T) {
	channels := NewChannelList()

	contacts := []Contact{
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8000"),
		NewContact(NewKademliaID("FFFFFFFF10000000000000000000000000000000"), "127.0.0.1:8000"),
		NewContact(NewKademliaID("FFFFFFFF20000000000000000000000000000000"), "127.0.0.1:8000"),
	}

	channels.Add(&contacts[0], 5, make(chan Header))

	_, err := channels.Find(&contacts[0], 6)
	if err == nil {
		t.Error("Channel found")
	}

	channels.Delete(&contacts[0], 5)
	_, err = channels.Find(&contacts[0], 5)
	if err == nil {
		t.Error("Channel found")
	}

	channels.Add(&contacts[0], 5, make(chan Header))
	channels.Add(&contacts[2], 8, make(chan Header))
	channels.Add(&contacts[1], 9, make(chan Header))

	_, err = channels.Find(&contacts[2], 8)
	if err != nil {
		t.Error("Channel not found")
	}
}
