package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Network struct {
	ID       *KademliaID
	IP       uint32
	Port     uint16
	Kademlia *Kademlia
}

type Header struct {
	SrcID   KademliaID
	SrcIP   uint32
	SrcPort uint16
	Type    uint8 /* Request/Response */
	SubType uint8

	Arg interface{}
}

type FindArguments struct {
	Key   KademliaID
	Count uint8
}

type StoreArguments struct {
	Key    KademliaID
	Length int
	Data   []byte
}

type ContactResult struct {
	ID      string
	Address string
}

func NewNetwork(id *KademliaID, ip string, port int) Network {
	return Network{
		ID:   id,
		IP:   IPToLong(ip),
		Port: uint16(port),
	}
}

func NewHeader(network *Network, typeId uint8, subTypeId uint8) Header {
	return Header{
		SrcID:   *network.ID,
		SrcIP:   network.IP,
		SrcPort: network.Port,
		Type:    typeId,
		SubType: subTypeId,
	}
}

func Encode(c *net.UDPConn, value interface{}) {
	log.Printf("Encoding: %v (%s -> %s)\n", value, c.LocalAddr(), c.RemoteAddr())

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(value)

	if err != nil {
		log.Fatalf("[ERROR] Encoding: %v\n", err)
	}

	c.Write(buffer.Bytes())
}

func Decode(c *net.UDPConn, value *Header) error {
	inputBytes := make([]byte, 1024)

	log.Printf("Listening...\n")
	length, addr, err := c.ReadFromUDP(inputBytes)
	if err != nil {
		log.Fatalf("[ERROR] Reading: %v\n", err)
	}

	log.Printf("Received from : %v\n", addr)

	buf := bytes.NewBuffer(inputBytes[:length])

	decoder := gob.NewDecoder(buf)

	err = decoder.Decode(value)

	value.SrcIP = IPToLong(addr.IP.String())

	log.Printf("Decoded: %v\n", value)

	if err != nil {
		log.Fatalf("Error decoding: %v\n", err)
	}

	return err
}

func (network *Network) SendPingMessage(contact *Contact) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	Encode(conn, NewHeader(network, MSG_REQUEST, MSG_PING))

	conn.Close()
}

/* Common function to SendFindDataMessage and SendFindContactMessage */
func (network *Network) sendFindMessage(contact *Contact, key *KademliaID, findType uint8) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := NewHeader(network, MSG_REQUEST, findType)
	msg.Arg = FindArguments{
		Count: K,
		Key:   *key,
	}
	Encode(conn, msg)

	conn.Close()
}

func (network *Network) SendFindContactMessage(contact *Contact, key *KademliaID) {
	network.sendFindMessage(contact, key, MSG_FIND_NODES)
}

func (network *Network) SendFindDataMessage(contact *Contact, key *KademliaID) {
	network.sendFindMessage(contact, key, MSG_FIND_VALUE)
}

func (network *Network) SendStoreMessage(contact *Contact, key *KademliaID, data []byte) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := NewHeader(network, MSG_REQUEST, MSG_STORE)
	msg.Arg = StoreArguments{
		Key:    *key,
		Length: len(data),
		Data:   data,
	}
	Encode(conn, msg)

	conn.Close()
}
