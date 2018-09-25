package kademlia

import (
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
		udpConn, _ := net.ListenUDP("udp", udpAddr)

		header := ReceiveHeader(udpConn)

		fmt.Printf("Header received: %v\n", header)

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

func connectTo(header Header) (*net.UDPConn, error) {
	addr, _ := net.ResolveUDPAddr("udp", IPToStr(header.SrcIP)+":"+strconv.Itoa(int(header.SrcPort)))
	conn, err := net.DialUDP("udp", nil, addr)

	return conn, err
}

func (kademlia *Kademlia) HandlePing(header Header) {
	if header.Type == MSG_REQUEST {
		conn, err := connectTo(header)

		if err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(1 * time.Second)

		fmt.Printf("Ping recu de: %s, %d\n", IPToStr(header.SrcIP), header.SrcPort)

		EncodeAndSend(conn, NewHeader(&kademlia.Network, MSG_RESPONSE, MSG_PING))
	}
}

func (kademlia *Kademlia) SendContacts(header Header, contacts []Contact) {
	conn, err := connectTo(header)

	if err != nil {
		fmt.Println(err)
		return
	}

	EncodeAndSend(conn, NewHeader(&kademlia.Network, MSG_RESPONSE, MSG_FIND_NODES))

	time.Sleep(1 * time.Second)

	/* Convert contacts ID to string before sending them */
	var contactsResults []ContactResult
	for _, c := range contacts {
		contactsResults = append(contactsResults, ContactResult{
			ID:      c.ID.String(),
			Address: c.Address,
		})
	}

	fmt.Println("Sending: ", contactsResults)

	EncodeAndSend(conn, contactsResults)
}

func (kademlia *Kademlia) HandleFindNodes(header Header, udpConn *net.UDPConn) {
	switch header.Type {
	case MSG_REQUEST:
		var findArguments FindArguments
		ReceiveAndDecode(udpConn, &findArguments)

		fmt.Printf("Argument received: %v\n", findArguments)

		var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(&(findArguments.Key), int(findArguments.Count))
		kademlia.SendContacts(header, contacts)
	case MSG_RESPONSE:
		var contacts []ContactResult
		ReceiveAndDecode(udpConn, &contacts)

		//TODO: traiter les contacts

		fmt.Printf("Contacts received: %v\n", contacts)
	}
}

func (kademlia *Kademlia) SendFile(header Header, filename string) {
	conn, err := connectTo(header)

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := NewHeader(&kademlia.Network, MSG_RESPONSE, MSG_FIND_VALUE)
	EncodeAndSend(conn, msg)

	time.Sleep(1 * time.Second)

	conn.Write(kademlia.Storage.Read(filename))

	conn.Close()
}

func (kademlia *Kademlia) HandleFindValue(header Header, udpConn *net.UDPConn) {
	switch header.Type {
	case MSG_REQUEST:
		var findArguments FindArguments
		ReceiveAndDecode(udpConn, &findArguments)
		key := findArguments.Key

		fmt.Printf("Argument received: %v\n", findArguments)

		if kademlia.Storage.Exists(key.String()) {
			kademlia.SendFile(header, key.String())
		} else {
			var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(&(key), int(findArguments.Count))
			kademlia.SendContacts(header, contacts)
		}
	case MSG_RESPONSE:
		content := make([]byte, 1024)

		length, _ := udpConn.Read(content)
		fmt.Println("Received: ", string(content[:length]))
		//JE RECOIS UN FICHIER
	}
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
