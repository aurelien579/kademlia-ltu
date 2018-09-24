package main

import (
	"kademlia"
)

func main() {
	/*if len(os.Args) != 2 {
	      fmt.Printf("Usage: %s <port>\n", os.Args[0])
	      return
	  }

	  port, _ := strconv.Atoi(os.Args[1])*/

	kademlia := kademlia.NewKademlia("0000000000000000000000000000000000000001", "127.0.0.1", 3333)
	kademlia.Listen("127.0.0.1", 3333)
}
