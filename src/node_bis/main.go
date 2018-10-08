package main

import (
	"bufio"
	"fmt"
	"kademlia"
	"os"
)

const MY_ID = "0000000000000000000000000000000000000002"
const MY_IP = "127.0.0.1"
const MY_PORT = 4100

func main() {
	node := kademlia.NewKademlia(MY_ID, MY_IP, MY_PORT)
	contact_id := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	contact := kademlia.NewContact(contact_id, "127.0.0.1:4000")

	node.RoutingTable.AddContact(contact)

	go node.Listen(MY_IP, MY_PORT)

	//bytes := []byte("BONJOURRR")
	//fmt.Printf("Data: %v\n ", bytes)
	//node.Store(bytes)

	node.Bootstrap(contact)

	for _, c := range node.RoutingTable.FindClosestContacts(node.RoutingTable.Me.ID, 50) {
		fmt.Println(c)
	}

	fmt.Printf("\n\n\n\n")

	//hash := "087245a9db4528d496941904ec716fdf3d8a6c3d"

	//fmt.Println(node.LookupData(hash))

	//node.Network.SendFindContactMessage(&contact, kademlia.NewKademliaID("000000000000000000000000000000000000FFFF"))
	//node.Network.SendFindDataMessage("000000000000000000000000000000000000FFFF")
	//node.Network.SendPingMessage(&contact)
	//go node.Listen(MY_IP, MY_PORT)

	//node.LookupContact(&contact)

	bufio.NewReader(os.Stdin).ReadString('\n')
}
