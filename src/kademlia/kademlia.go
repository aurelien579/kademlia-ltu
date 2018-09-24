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
	Storage      Storage
}

func NewKademlia(id string, ip string, port int) Kademlia {
	me := NewContact(NewKademliaID(id), ip+strconv.Itoa(port))
	kademlia := Kademlia{
		RoutingTable: NewRoutingTable(me),
		Network:      NewNetwork(me.ID, ip, port),
		Storage:      NewStorage(id),
	}

	kademlia.Network.Kademlia = &kademlia

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

		kademlia.RoutingTable.AddContact(NewContact(&(header.SrcID), IPToStr(header.SrcIP)))

		switch header.SubType {

		case MSG_PING:

			kademlia.HandlePing(header)

		case MSG_FIND_NODES:

			kademlia.HandleFindNodes(header, udpConn)

		case MSG_FIND_VALUE:

			kademlia.HandleFindValue(header, udpConn)

		case MSG_STORE:

			kademlia.HandleStore(header)

		}

		udpConn.Close()
		// TODO: gérer le message dans un nouveau thread + AddContact
	}
}

func (kademlia *Kademlia) HandlePing(header Header) {

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

func (kademlia *Kademlia) HandleFindNodes(header Header, udpConn *net.UDPConn) {

	switch header.Type {

	case MSG_REQUEST:
		inputBytes := make([]byte, 1024)
		length, _ := udpConn.Read(inputBytes)
		buf := bytes.NewBuffer(inputBytes[:length])

		decoder := gob.NewDecoder(buf)
		var findArguments FindArguments
		decoder.Decode(&findArguments)

		fmt.Printf("Argument received: %v\n", findArguments)

		var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(&(findArguments.Key), int(findArguments.Count))

		kademlia.SendContact(header, contacts)

		udpConn.Close()

	case MSG_RESPONSE:
		inputBytes := make([]byte, 1024)
		length, _ := udpConn.Read(inputBytes)
		fmt.Println("REceived length: ", length)
		buf := bytes.NewBuffer(inputBytes[:length])

		decoder := gob.NewDecoder(buf)
		var contacts []ContactResult
		err := decoder.Decode(&contacts)

		//traiter les contacts

		fmt.Println("Error: ", err)
		fmt.Printf("Contacts received: %v\n", contacts)

		udpConn.Close()
	}

}

func (kademlia *Kademlia) HandleFindValue(header Header, udpConn *net.UDPConn) {

	switch header.Type {

	case MSG_REQUEST:
		inputBytes := make([]byte, 1024)
		length, _ := udpConn.Read(inputBytes)
		buf := bytes.NewBuffer(inputBytes[:length])

		decoder := gob.NewDecoder(buf)
		var findArguments FindArguments
		decoder.Decode(&findArguments)

		fmt.Printf("Argument received: %v\n", findArguments)

		fmt.Println(findArguments.Key.String(), " exists ? ", kademlia.Storage.Exists(findArguments.Key.String()))

		if kademlia.Storage.Exists(findArguments.Key.String()) {

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
				SubType: MSG_FIND_VALUE,
			}

			enc.Encode(msg)

			time.Sleep(1 * time.Second)

			conn.Write(buffer.Bytes())

			buffer.Reset()

			conn.Write(kademlia.Storage.Read(findArguments.Key.String()))

			udpConn.Close()

		} else {

			var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(&(findArguments.Key), int(findArguments.Count))

			kademlia.SendContact(header, contacts)

			udpConn.Close()

		}

	case MSG_RESPONSE:
		content := make([]byte, 1024)

		length, _ := udpConn.Read(content)
		fmt.Println("Received: ", string(content[:length]))
		//JE RECOIS UN FICHIER

	}

}

func (kademlia *Kademlia) SendContact(header Header, contacts []Contact) {

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
		SubType: MSG_FIND_NODES,
	}

	enc.Encode(msg)

	time.Sleep(1 * time.Second)

	conn.Write(buffer.Bytes())

	buffer.Reset()

	var contactsResults []ContactResult
	for _, c := range contacts {
		contactsResults = append(contactsResults, ContactResult{
			ID:      c.ID.String(),
			Address: c.Address,
		})
	}
	fmt.Println("Sending: ", contactsResults)

	enc.Encode(contactsResults)

	time.Sleep(1 * time.Second)

	conn.Write(buffer.Bytes())

}

func (kademlia *Kademlia) HandleStore(header Header) {

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
