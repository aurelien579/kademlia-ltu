package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type Network struct {
	ID   *KademliaID
	IP   uint32
	Port uint16
}

const MSG_REQUEST uint8 = 1
const MSG_RESPONSE uint8 = 2

const MSG_PING uint8 = 1
const MSG_FIND_NODES uint8 = 2
const MSG_FIND_VALUE uint8 = 3
const MSG_STORE uint8 = 4

type Header struct {
	SrcID   KademliaID
	SrcIP   uint32
	SrcPort uint16
	Type    uint8 /* Request/Response */
	SubType uint8
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
		Type:    MSG_REQUEST,
		SubType: MSG_PING,
	}
}

func sendHeader(connection *net.UDPConn, header Header) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(header)
	connection.Write(buffer.Bytes())
}

func connectAndSendHeader(contact *Contact, header Header) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	sendHeader(conn, header)

	conn.Close()
}

func (network *Network) SendPingMessage(contact *Contact) {
	msg := network.createHeader(MSG_REQUEST, MSG_PING)
	connectAndSendHeader(contact, msg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
