package main

import (
	"kademlia"
)

func main() {
	/*if len(os.Args) != 2 {
	    fmt.Printf("Usage: %s <port>\n", os.Args[0])
	    return
	}*/

	//port, _ := strconv.Atoi(os.Args[1])

	contact := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), "localhost:3333")

	kademlia := kademlia.NewKademlia("0000000000000000000000000000000000000002", "127.0.0.1", 3334)
	//kademlia.Network.SendPingMessage(&contact)
	kademlia.Network.SendFindContactMessage(&contact, contact.ID)
	kademlia.Listen("127.0.0.1", 3334)
}
