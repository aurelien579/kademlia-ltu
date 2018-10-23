package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Network struct {
	ID       *KademliaID
	IP       uint32
	Port     uint16
	Kademlia *Kademlia
}

type Header struct {
	SrcID   string
	SrcIP   uint32
	SrcPort uint16
	Type    uint8 /* Request/Response */
	SubType uint8

	Arg interface{}
}

type FindArguments struct {
	Key   string
	Count uint8
}

type StoreArguments struct {
	Key    string
	Length int
	Data   []byte
}

func (header Header) String() string {
	str := "header("

	str += header.SrcID + ", " + IPToStr(header.SrcIP) + ":" + strconv.Itoa(int(header.SrcPort)) + ", "

	if header.Type == MSG_REQUEST {
		str += "REQUEST, "
	} else {
		str += "RESPONSE, "
	}

	if header.SubType == MSG_PING {
		str += "PING, "
	} else if header.SubType == MSG_STORE {
		str += "STORE, "
	} else if header.SubType == MSG_FIND_NODES {
		str += "FIND_NODES, "
	} else if header.SubType == MSG_PIN {
		str += "PIN, "
	} else if header.SubType == MSG_UNPIN {
		str += "UNPIN, "
	} else {
		str += "FIND_VALUE, "
	}

	str += fmt.Sprintf("%v)", header.Arg)

	return str
}

func NewNetwork(id *KademliaID, ip string, port int) *Network {
	gob.Register(FindArguments{})
	gob.Register(StoreArguments{})
	gob.Register([]ContactResult{})

	return &Network{
		ID:   id,
		IP:   IPToLong(ip),
		Port: uint16(port),
	}
}

func NewHeader(network *Network, typeId uint8, subTypeId uint8) Header {
	return Header{
		SrcID:   network.ID.String(),
		SrcIP:   network.IP,
		SrcPort: network.Port,
		Type:    typeId,
		SubType: subTypeId,
	}
}

func (header *Header) createConnection() (*net.UDPConn, error) {
	addr, _ := net.ResolveUDPAddr("udp", IPToStr(header.SrcIP)+":"+strconv.Itoa(int(header.SrcPort)))
	conn, err := net.DialUDP("udp", nil, addr)

	return conn, err
}

func Encode(c *net.UDPConn, value Header) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(value)

	if err != nil {
		log.Fatalf("[ERROR] Encoding: %v\n", err)
		return
	}

	log.Printf("Sending %v \n\t(%s -> %s)\n", value, c.LocalAddr(), c.RemoteAddr())

	c.Write(buffer.Bytes())
}

func Decode(c *net.UDPConn, value *Header) error {
	inputBytes := make([]byte, 1024)

	log.Printf("Listening...\n")
	length, addr, err := c.ReadFromUDP(inputBytes)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(inputBytes[:length])

	decoder := gob.NewDecoder(buf)

	err = decoder.Decode(value)

	//value.SrcIP = IPToLong(addr.IP.String())

	log.Printf("Received (%v): %v\n", addr.String(), value)

	if err != nil {
		return err
	}

	return nil
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
		Key:   key.String(),
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
		Key:    key.String(),
		Length: len(data),
		Data:   data,
	}
	Encode(conn, msg)

	conn.Close()
}

func (network *Network) SendPinMessage(contact *Contact, key *KademliaID, data []byte) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := NewHeader(network, MSG_REQUEST, MSG_PIN)
	msg.Arg = StoreArguments{
		Key:    key.String(),
		Length: len(data),
		Data:   data,
	}
	Encode(conn, msg)

	conn.Close()
}

func (network *Network) SendUnpinMessage(contact *Contact, key *KademliaID) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := NewHeader(network, MSG_REQUEST, MSG_UNPIN)
	msg.Arg = StoreArguments{
		Key: key.String(),
	}
	Encode(conn, msg)

	conn.Close()
}
