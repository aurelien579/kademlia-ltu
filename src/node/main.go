package main

import (
	"fmt"
	"kademlia"
	"os"
	"strconv"
)

const MY_ID = "000000000000000000000000000000000000FFFF"
const MY_IP = "127.0.0.1"
const MY_PORT = 3333

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <id> <port> <contactId> <contactPort>\n", os.Args[0])
		return
	}

	port, _ := strconv.Atoi(os.Args[2])
	node := kademlia.NewKademlia(os.Args[1], MY_IP, port)

	contactPort, _ := strconv.Atoi(os.Args[4])
	if contactPort != 0 {
		contact := kademlia.NewKademlia(os.Args[3], MY_IP, contactPort)

		node.RoutingTable.AddContact(kademlia.NewContact(contact.RoutingTable.Me.ID, "127.0.0.1:"+os.Args[4]))
	}

	//node.Storage.Store("000000000000000000000000000000000000FFFF", []byte("bonjour"))
	node.Listen(MY_IP, port)
}
