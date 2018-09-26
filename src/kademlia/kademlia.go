package kademlia

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

const ALPHA = 1

type Kademlia struct {
	RoutingTable     *RoutingTable
	Network          Network
	Storage          Storage
	ResponseHandlers *list.List
}

type ResponseHandler struct {
	Contact *Contact
	F       ResponseHandlerFunc
	ReqType uint8
}

type ResponseHandlerFunc func(*Contact, interface{})

func NewKademlia(id string, ip string, port int) Kademlia {
	me := NewContact(NewKademliaID(id), ip+":"+strconv.Itoa(port))
	kademlia := Kademlia{
		RoutingTable: NewRoutingTable(me),
		Network:      NewNetwork(me.ID, ip, port),
		Storage:      NewStorage(id),
	}

	kademlia.Network.Kademlia = &kademlia

	kademlia.ResponseHandlers = list.New()

	return kademlia
}

func (kademlia *Kademlia) RegisterHandler(contact *Contact, reqType uint8, f ResponseHandlerFunc) {
	handler := ResponseHandler{
		Contact: contact,
		ReqType: reqType,
		F:       f,
	}

	kademlia.ResponseHandlers.PushBack(handler)
}

func (kademlia *Kademlia) FindHandler(contact *Contact, reqType uint8) (ResponseHandler, error) {
	for e := kademlia.ResponseHandlers.Front(); e != nil; e = e.Next() {
		handler := e.Value.(ResponseHandler)

		if handler.Contact.Address == contact.Address &&
			*handler.Contact.ID == *contact.ID &&
			handler.ReqType == reqType {

			return handler, nil
		}
	}

	return ResponseHandler{}, errors.New("Can't find handler")
}

func (kademlia *Kademlia) CallHandler(header *Header, val interface{}) {
	contact := ContactFromHeader(header)
	handler, err := kademlia.FindHandler(&contact, header.SubType)
	if err != nil {
		fmt.Println(err)
		return
	}

	handler.F(&contact, val)

	// TODO: delete handler
}

func (kademlia *Kademlia) Listen(ip string, port int) {
	udpAddr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))

	for {
		udpConn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var header Header
		Decode(udpConn, &header)

		fmt.Printf("Header received: %v\n", header)

		kademlia.RoutingTable.AddContact(ContactFromHeader(&header))

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

		Encode(conn, NewHeader(&kademlia.Network, MSG_RESPONSE, MSG_PING))
	} else {
		kademlia.CallHandler(&header, nil)
	}
}

func (kademlia *Kademlia) SendContacts(header Header, contacts []Contact) {
	conn, err := connectTo(header)

	if err != nil {
		fmt.Println(err)
		return
	}

	header = NewHeader(&kademlia.Network, MSG_RESPONSE, MSG_FIND_NODES)
	fmt.Println("Sending header: ", header)
	Encode(conn, header)

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

	Encode(conn, contactsResults)
}

func (kademlia *Kademlia) HandleFindNodes(header Header, udpConn *net.UDPConn) {
	switch header.Type {
	case MSG_REQUEST:
		var findArguments FindArguments
		Decode(udpConn, &findArguments)

		fmt.Printf("Argument received: %v\n", findArguments)

		var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(&(findArguments.Key), int(findArguments.Count))
		kademlia.SendContacts(header, contacts)
	case MSG_RESPONSE:
		var contactResults []ContactResult
		Decode(udpConn, &contactResults)

		var contacts []Contact
		for _, r := range contactResults {
			contacts = append(contacts, NewContact(NewKademliaID(r.ID), r.Address))
		}

		fmt.Println("Calling handler")
		kademlia.CallHandler(&header, contacts)

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
	Encode(conn, msg)

	time.Sleep(1 * time.Second)

	conn.Write(kademlia.Storage.Read(filename))

	conn.Close()
}

func (kademlia *Kademlia) HandleFindValue(header Header, udpConn *net.UDPConn) {
	switch header.Type {
	case MSG_REQUEST:
		var findArguments FindArguments
		Decode(udpConn, &findArguments)
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

func (kademlia *Kademlia) lookupThread(i int, wg *sync.WaitGroup, l *ContactCandidates, key *KademliaID, done *bool) {
	fmt.Printf("Thread %d starting\n", i)

	channel := make(chan int)

	for !*done {
		if l.Len() == 0 {
			fmt.Println("No candidates")
			time.Sleep(1 * time.Second)
			continue
		}

		// TODO: check undone
		target, err := l.GetClosestUnTook(int(K))
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Println("Candidates: ", l)
		fmt.Printf("Thread %d target: %v\n", i, target)

		kademlia.RegisterHandler(target, MSG_FIND_NODES, func(c *Contact, val interface{}) {
			fmt.Println("HANDLER")

			results := val.([]Contact)

			for i := 0; i < len(results); i++ {
				results[i].CalcDistance(key)
			}

			target.State = DONE

			fmt.Printf("Thread %d results: [\n", i)

			for _, r := range results {
				fmt.Println(r, ",")
			}
			fmt.Println("]")

			l.Append(results)
			l.Sort()

			if l.Finish(int(K)) {
				*done = true
			}

			channel <- 1
		})

		kademlia.Network.SendFindContactMessage(target, key)

		<-channel
	}

	wg.Done()

	fmt.Printf("Thread %d done\n", i)
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	time.Sleep(1 * time.Second)

	var wg sync.WaitGroup
	var candidates ContactCandidates
	done := false
	closestKnown := kademlia.RoutingTable.FindClosestContacts(target.ID, int(K))

	me := NewContact(kademlia.RoutingTable.Me.ID, kademlia.RoutingTable.Me.Address)
	me.State = DONE

	closestKnown = append(closestKnown, me)
	candidates.Append(closestKnown)

	for i := 0; i < candidates.Len(); i++ {
		candidates.contacts[i].CalcDistance(target.ID)
	}
	candidates.Sort()

	fmt.Println("Initial Candidates: ", candidates)

	wg.Add(ALPHA)
	for i := 0; i < ALPHA; i++ {
		go kademlia.lookupThread(i, &wg, &candidates, target.ID, &done)
	}
	wg.Wait()

	fmt.Printf("Lookup done! Results: %v\n", &candidates)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
