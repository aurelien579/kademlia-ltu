package main

import (
	"kademlia"
)

const MY_ID = "0000000000000000000000000000000000000002"
const MY_IP = "127.0.0.1"
const MY_PORT = 3334

func main() {
	node := kademlia.NewKademlia(MY_ID, MY_IP, MY_PORT)
	contact := kademlia.NewContact(kademlia.NewKademliaID("000000000000000000000000000000000000FFFF"), "localhost:3333")

	node.RoutingTable.AddContact(contact)

	//node.Network.SendFindContactMessage(&contact, kademlia.NewKademliaID("000000000000000000000000000000000000FFFF"))
	node.Network.SendFindDataMessage("000000000000000000000000000000000000FFFF")
	//node.Network.SendPingMessage(&contact)
	node.Listen(MY_IP, MY_PORT)
}
