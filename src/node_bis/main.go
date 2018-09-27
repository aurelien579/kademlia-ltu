package main

import (
	"fmt"
	"kademlia"
)

const MY_ID = "0000000000000000000000000000000000000001"
const MY_IP = "127.0.0.1"
const MY_PORT = 3334

func main() {
	node := kademlia.NewKademlia(MY_ID, MY_IP, MY_PORT)
	contact_id := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	contact := kademlia.NewContact(contact_id, "127.0.0.1:3336")

	//node.RoutingTable.AddContact(contact)

	/*node.RegisterHandler(&contact, kademlia.MSG_PING, func(contact *kademlia.Contact, val interface{}) {
		fmt.Println("Ping back!!!")
	})*/

	go node.Listen(MY_IP, MY_PORT)
	node.Bootstrap(contact)

	for _, c := range node.RoutingTable.FindClosestContacts(node.RoutingTable.Me.ID, 50) {
		fmt.Println(c)
	}

	//node.Network.SendFindContactMessage(&contact, kademlia.NewKademliaID("000000000000000000000000000000000000FFFF"))
	//node.Network.SendFindDataMessage("000000000000000000000000000000000000FFFF")
	//node.Network.SendPingMessage(&contact)
	//go node.Listen(MY_IP, MY_PORT)

	//node.LookupContact(&contact)
}
