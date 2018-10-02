package main

import (
	"fmt"
	"kademlia"
	"net"
	"strconv"
	"time"
)

const MY_ID = "000000000000000000000000000000000000FFFF"
const MY_IP = "127.0.0.1"
const MY_PORT = 3333

func getFreePort(start int) int {
	for {
		if start >= 65555 {
			return 0
		}

		addr, _ := net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(start))
		conn, err := net.ListenUDP("udp", addr)

		if err != nil {
			start++
		} else {
			conn.Close()
			return start
		}
	}
}

func main() {
	var node kademlia.Kademlia
	port := getFreePort(4000)

	fmt.Println(port)

	if port == 4000 {
		node = kademlia.NewKademlia("0000000000000000000000000000000000000001", "127.0.0.1", port)
	} else {
		node = kademlia.NewKademlia("000000000000000000000000000000000000"+strconv.Itoa(port), "127.0.0.1", port)
	}

	go node.Listen("127.0.0.1", port)

	if port != 4000 {
		node.Bootstrap(kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1:4000"))
	}

	if port == 4001 {
		time.Sleep(10 * time.Second)
		node.Store([]byte("Toto"))
	}

	for {

	}
}
