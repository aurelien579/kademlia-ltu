package main

import (
	"fmt"
	"kademlia"
)

const MY_ID = "0000000000000000000000000000000000000002"
const MY_IP = "127.0.0.1"
const MY_PORT = 3334

func main() {
	node := kademlia.NewKademlia(MY_ID, MY_IP, MY_PORT)
	contact_id := kademlia.NewKademliaID("000000000000000000000000000000000000FFFF")
	contact := kademlia.NewContact(contact_id, "127.0.0.1:3333")

	node.RoutingTable.AddContact(contact)

	node.RegisterHandler(&contact, kademlia.MSG_PING, func(contact *kademlia.Contact, val interface{}) {
		fmt.Println("Ping back!!!")
	})

	//node.Network.SendFindContactMessage(&contact, kademlia.NewKademliaID("000000000000000000000000000000000000FFFF"))
	node.Network.SendFindDataMessage("000000000000000000000000000000000000FFFF")
	//node.Network.SendPingMessage(&contact)
	node.Listen(MY_IP, MY_PORT)
}
