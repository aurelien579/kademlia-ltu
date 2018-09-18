package main

import (
    "kademlia"
	"fmt"
	"net"
	"encoding/gob"
	"bytes"
	"strconv"
)

var node kademlia.Node

const MSG_REQUEST uint8 = 1
const MSG_RESPONSE uint8 = 2

const MSG_PING uint8 = 1
const MSG_FIND_NODES uint8 = 2
const MSG_FIND_VALUE uint8 = 3
const MSG_STORE uint8 = 4

type Header struct {
    SrcIP uint32
    SrcPort uint8
    Type uint8      /* Request/Response */
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
        udpAddr, _ := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(port))
        udpConn, err := net.ListenUDP("udp", udpAddr)
                
        if err == nil {        
            udpConn.Close()
            return port
        }
    }

}

func main() {
    port := getFreePort(5555)
    udpAddr, _ := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(port))
    
    fmt.Printf("port found: %d\n", port)
    
    udpConn, err := net.ListenUDP("udp", udpAddr)

    if err != nil {
        fmt.Println(err)
        return
    }
    
    handleConn(udpConn)
}