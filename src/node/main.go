package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"kademlia"
	"net"
	"strconv"
)

const MSG_REQUEST uint8 = 1
const MSG_RESPONSE uint8 = 2

const MSG_PING uint8 = 1
const MSG_FIND_NODES uint8 = 2
const MSG_FIND_VALUE uint8 = 3
const MSG_STORE uint8 = 4

type Header struct {
	SrcIP   uint32
	SrcPort uint8
	Type    uint8 /* Request/Response */
	SubType uint8
}

func handleConn(c *net.UDPConn) {
	inputBytes := make([]byte, 1024)
	length, _ := c.Read(inputBytes)
	buf := bytes.NewBuffer(inputBytes[:length])

	decoder := gob.NewDecoder(buf)
	var header Header
	decoder.Decode(&header)

	fmt.Println(header)
}

func getFreePort(start int) int {
	for port := start; ; port++ {
		udpAddr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
		udpConn, err := net.ListenUDP("udp", udpAddr)

		if err == nil {
			udpConn.Close()
			return port
		}
	}

}

func main() {
	/*if len(os.Args) != 2 {
	      fmt.Printf("Usage: %s <port>\n", os.Args[0])
	      return
	  }

	  port, _ := strconv.Atoi(os.Args[1])*/

	kademlia := kademlia.NewKademlia("0000000000000000000000000000000000000001", "127.0.0.1", 3333)
	kademlia.Listen("127.0.0.1", 3333)
}
