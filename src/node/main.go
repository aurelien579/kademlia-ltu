package main

import (
	"kademlia"
)

const MY_ID = "000000000000000000000000000000000000FFFF"
const MY_IP = "127.0.0.1"
const MY_PORT = 3333

func main() {
	/*if len(os.Args) != 2 {
	      fmt.Printf("Usage: %s <port>\n", os.Args[0])
	      return
	  }

	  port, _ := strconv.Atoi(os.Args[1])*/

	node := kademlia.NewKademlia(MY_ID, MY_IP, MY_PORT)
	node.Storage.Store("000000000000000000000000000000000000FFFF", []byte("bonjour"))
	node.Listen(MY_IP, MY_PORT)
}
