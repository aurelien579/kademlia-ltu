package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type Network struct {
	ID       *KademliaID
	IP       uint32
	Port     uint16
	Kademlia *Kademlia
}

const K uint8 = 20

const MSG_REQUEST uint8 = 1
const MSG_RESPONSE uint8 = 2

const MSG_PING uint8 = 1
const MSG_FIND_NODES uint8 = 45
const MSG_FIND_VALUE uint8 = 3
const MSG_STORE uint8 = 4

type Header struct {
	SrcID   KademliaID
	SrcIP   uint32
	SrcPort uint16
	Type    uint8 /* Request/Response */
	SubType uint8
}

type FindArguments struct {
	Key   KademliaID
	Count uint8
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

func (network *Network) createHeader(typeId uint8, subTypeId uint8) Header {
	return Header{
		SrcID:   *network.ID,
		SrcIP:   network.IP,
		SrcPort: network.Port,
		Type:    typeId,
		SubType: subTypeId,
	}
}

func encodeAndSend(c *net.UDPConn, value interface{}) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(value)
	c.Write(buffer.Bytes())
}

func connectAndSendHeader(contact *Contact, header Header) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	encodeAndSend(conn, header)

	conn.Close()
}

func (network *Network) SendPingMessage(contact *Contact) {
	msg := network.createHeader(MSG_REQUEST, MSG_PING)
	connectAndSendHeader(contact, msg)
}

/* Common function to SendFindDataMessage and SendFindContactMessage */
func (network *Network) sendFindMessage(contact *Contact, key *KademliaID, findType uint8) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	header := network.createHeader(MSG_REQUEST, findType)
	fmt.Println("Sending: ", header)
	encodeAndSend(conn, header)

	findMessage := FindArguments{
		Count: K,
		Key:   *key,
	}
	encodeAndSend(conn, findMessage)

	conn.Close()
}

func (network *Network) SendFindContactMessage(contact *Contact, key *KademliaID) {
	network.sendFindMessage(contact, key, MSG_FIND_NODES)
}

func (network *Network) SendFindDataMessage(hash string) {
	key := NewKademliaID(hash)

	closest := network.Kademlia.RoutingTable.FindClosestContacts(key, 1)

	fmt.Println("closest: ", closest)
	network.sendFindMessage(&closest[0], key, MSG_FIND_VALUE)
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
