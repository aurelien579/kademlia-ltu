package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      Network
}

func NewKademlia(id string, ip string, port int) Kademlia {
	me := NewContact(NewKademliaID(id), ip+strconv.Itoa(port))
	kademlia := Kademlia{
		RoutingTable: NewRoutingTable(me),
		Network:      NewNetwork(me.ID, ip, port),
	}

	return kademlia
}

func (kademlia *Kademlia) Listen(ip string, port int) {
	udpAddr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))

	for {
		udpConn, err := net.ListenUDP("udp", udpAddr)

		if err != nil {
			fmt.Println(err)
			return
		}

		inputBytes := make([]byte, 1024)
		length, _ := udpConn.Read(inputBytes)
		buf := bytes.NewBuffer(inputBytes[:length])

		decoder := gob.NewDecoder(buf)
		var header Header
		decoder.Decode(&header)

		fmt.Printf("Header received: %v\n", header)
		fmt.Printf("Address: %s\n", IPToStr(header.SrcIP))

		switch header.SubType {

		case MSG_PING:
			kademlia.HandlePing(header)

		case MSG_FIND_NODES:

		case MSG_FIND_VALUE:

		case MSG_STORE:

		}

		udpConn.Close()
		// TODO: gérer le message dans un nouveau thread + AddContact
	}
}

func (kademlia *Kademlia) HandlePing(header Header) {

	kademlia.RoutingTable.AddContact(NewContact(&(header.SrcID), IPToStr(header.SrcIP)))

	if header.Type == MSG_REQUEST {

		addr, _ := net.ResolveUDPAddr("udp", IPToStr(header.SrcIP)+":"+strconv.Itoa(int(header.SrcPort)))
		conn, err := net.DialUDP("udp", nil, addr)

		if err != nil {
			fmt.Println(err)
			return
		}

		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		msg := Header{
			SrcID:   *kademlia.Network.ID,
			SrcIP:   kademlia.Network.IP,
			SrcPort: kademlia.Network.Port,
			Type:    MSG_RESPONSE,
			SubType: MSG_PING,
		}

		fmt.Printf("Ping recu de: %s, %d\n", IPToStr(header.SrcIP), header.SrcPort)

		enc.Encode(msg)

		time.Sleep(1 * time.Second)

		conn.Write(buffer.Bytes())

	}

}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
