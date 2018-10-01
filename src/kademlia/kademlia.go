package kademlia

import (
	"bytes"
	"container/list"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
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
	mutex            sync.Mutex
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

	kademlia.Storage.kademlia = &kademlia

	kademlia.ResponseHandlers = list.New()

	for i := 0; i < IDLength*8; i++ {
		kademlia.RoutingTable.Buckets[i].node = &kademlia
	}

	log.SetFlags(log.Ltime | log.Lmicroseconds)
	gob.Register(FindArguments{})
	gob.Register(StoreArguments{})
	gob.Register([]ContactResult{})

	return kademlia
}

func (kademlia *Kademlia) RegisterHandler(contact *Contact, reqType uint8, f ResponseHandlerFunc) {
	handler := ResponseHandler{
		Contact: contact,
		ReqType: reqType,
		F:       f,
	}

	kademlia.mutex.Lock()
	kademlia.ResponseHandlers.PushBack(handler)
	kademlia.mutex.Unlock()
}

func (kademlia *Kademlia) FindHandler(contact *Contact, reqType uint8) (ResponseHandler, error) {
	kademlia.mutex.Lock()

	for e := kademlia.ResponseHandlers.Front(); e != nil; e = e.Next() {
		handler := e.Value.(ResponseHandler)

		if handler.Contact.Address == contact.Address &&
			*handler.Contact.ID == *contact.ID &&
			handler.ReqType == reqType {
			kademlia.mutex.Unlock()
			return handler, nil
		}
	}

	kademlia.mutex.Unlock()
	return ResponseHandler{}, errors.New("Can't find handler")
}

func (kademlia *Kademlia) DeleteHandler(contact *Contact, reqType uint8) {
	kademlia.mutex.Lock()

	for e := kademlia.ResponseHandlers.Front(); e != nil; e = e.Next() {
		handler := e.Value.(ResponseHandler)

		if handler.Contact.Address == contact.Address &&
			*handler.Contact.ID == *contact.ID &&
			handler.ReqType == reqType {

			kademlia.ResponseHandlers.Remove(e)
		}
	}

	kademlia.mutex.Unlock()
}

func (kademlia *Kademlia) CallHandler(header *Header, val interface{}) {
	contact := ContactFromHeader(header)
	handler, err := kademlia.FindHandler(&contact, header.SubType)
	if err != nil {
		fmt.Println(err)
		return
	}

	kademlia.mutex.Lock()
	handler.F(&contact, val)
	kademlia.mutex.Unlock()

	kademlia.DeleteHandler(&contact, header.SubType)
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

		kademlia.mutex.Lock()
		kademlia.RoutingTable.AddContact(ContactFromHeader(&header))
		kademlia.mutex.Unlock()

		switch header.SubType {
		case MSG_PING:
			kademlia.HandlePing(header)
		case MSG_FIND_NODES:
			kademlia.HandleFindNodes(header, udpConn)
		case MSG_FIND_VALUE:
			kademlia.HandleFindValue(header, udpConn)
		case MSG_STORE:
			kademlia.HandleStore(header, udpConn)
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

	header = NewHeader(&kademlia.Network, MSG_RESPONSE, header.SubType)

	/* Convert contacts ID to string before sending them */
	var contactsResults []ContactResult
	for _, c := range contacts {
		contactsResults = append(contactsResults, ContactResult{
			ID:      c.ID.String(),
			Address: c.Address,
		})
	}

	header.Arg = contactsResults

	Encode(conn, header)
}

func (kademlia *Kademlia) HandleFindNodes(header Header, udpConn *net.UDPConn) {
	switch header.Type {
	case MSG_REQUEST:
		findArguments := header.Arg.(FindArguments)

		var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(&(findArguments.Key), int(findArguments.Count))
		kademlia.SendContacts(header, contacts)
	case MSG_RESPONSE:
		contactResults := header.Arg.([]ContactResult)

		var contacts []Contact
		for _, r := range contactResults {
			contacts = append(contacts, NewContact(NewKademliaID(r.ID), r.Address))
		}

		kademlia.CallHandler(&header, contacts)
	}
}

func (kademlia *Kademlia) SendFile(header Header, filename string) {
	conn, err := connectTo(header)

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := NewHeader(&kademlia.Network, MSG_RESPONSE, MSG_FIND_VALUE)

	data := kademlia.Storage.Read(filename)

	msg.Arg = StoreArguments{
		Key:    *NewKademliaID(filename),
		Length: len(data),
		Data:   data,
	}

	Encode(conn, msg)

	conn.Close()
}

func (kademlia *Kademlia) HandleFindValue(header Header, udpConn *net.UDPConn) {
	switch header.Type {
	case MSG_REQUEST:
		findArguments := header.Arg.(FindArguments)
		key := findArguments.Key

		if kademlia.Storage.Exists(key.String()) {
			kademlia.SendFile(header, key.String())
		} else {
			var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(&(key), int(findArguments.Count))
			kademlia.SendContacts(header, contacts)
		}
	case MSG_RESPONSE:
		isContacts := true
		if reflect.TypeOf(header.Arg).String() != "[]kademlia.ContactResult" {
			isContacts = false
		}

		if !isContacts { /* Process file */
			kademlia.CallHandler(&header, header.Arg.(StoreArguments).Data)
		} else { /* Process contacts */
			var contacts []Contact
			for _, r := range header.Arg.([]ContactResult) {
				contacts = append(contacts, NewContact(NewKademliaID(r.ID), r.Address))
			}

			kademlia.CallHandler(&header, contacts)
		}
	}
}

func (kademlia *Kademlia) HandleStore(header Header, conn *net.UDPConn) {
	fmt.Println("HANDLER STORE: ", header)

	if header.Type != MSG_REQUEST {
		return
	}

	args := header.Arg.(StoreArguments)

	fmt.Println("Store received: ", args.Length, args.Data)

	kademlia.Storage.Store(args.Key.String(), args.Data)
}

func (kademlia *Kademlia) lookupThread(i int, wg *sync.WaitGroup, l *ContactCandidates, key *KademliaID, done *bool, lookupData bool, buffer *bytes.Buffer) {
	log.Printf("\n\n\n\n============ Thread %d starting ============\n\n", i)
	log.Println("Initial candidates: ", l)

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

		log.Printf("Thread %d target: %v\n", i, target)

		var subType = MSG_FIND_NODES
		if lookupData {
			subType = MSG_FIND_VALUE
		}
		kademlia.RegisterHandler(target, subType, func(c *Contact, val interface{}) {
			log.Printf("\n\n======= Handler called =======\n\n")

			isContacts := true
			if reflect.TypeOf(val).String() != "[]kademlia.Contact" {
				isContacts = false
			}

			target.State = DONE

			if isContacts {
				contacts := val.([]Contact)

				for i := 0; i < len(contacts); i++ {
					contacts[i].CalcDistance(key)
				}

				log.Printf("Thread %d results: [\n", i)
				for _, r := range contacts {
					log.Println("\t", r)
				}
				log.Println("]")

				l.Append(contacts)
				l.Sort()
			} else {
				*done = true
				buffer.Write(val.([]byte))
			}

			if !lookupData {
				if l.Finish(int(K)) {
					*done = true
				}
			}

			channel <- 1
		})

		if lookupData {
			kademlia.Network.SendFindDataMessage(target, key)
		} else {
			kademlia.Network.SendFindContactMessage(target, key)
		}

		<-channel
	}

	wg.Done()

	log.Printf("============ Thread %d done ============\n", i)

	close(channel)
}

func (kademlia *Kademlia) Lookup(key *KademliaID, lookupData bool) interface{} {
	time.Sleep(1 * time.Second)

	var wg sync.WaitGroup
	var candidates ContactCandidates
	var buffer bytes.Buffer
	done := false
	closestKnown := kademlia.RoutingTable.FindClosestContacts(key, int(K))

	me := NewContact(kademlia.RoutingTable.Me.ID, kademlia.RoutingTable.Me.Address)
	me.State = DONE

	closestKnown = append(closestKnown, me)
	candidates.Append(closestKnown)

	for i := 0; i < candidates.Len(); i++ {
		candidates.contacts[i].CalcDistance(key)
	}
	candidates.Sort()

	wg.Add(ALPHA)
	for i := 0; i < ALPHA; i++ {
		go kademlia.lookupThread(i, &wg, &candidates, key, &done, lookupData, &buffer)
	}
	wg.Wait()

	if lookupData {
		return buffer.Bytes()
	} else {
		return candidates
	}
}

func (kademlia *Kademlia) LookupContact(key *KademliaID) ContactCandidates {
	return kademlia.Lookup(key, false).(ContactCandidates)
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	key := NewKademliaID(hash)
	return kademlia.Lookup(key, true).([]byte)
}

func (kademlia *Kademlia) Store(data []byte) {
	h := sha1.New()
	io.WriteString(h, string(data))
	key := NewKademliaID(hex.EncodeToString(h.Sum(nil)))

	candidates := kademlia.LookupContact(key)

	for i := 0; i < Min(candidates.Len(), int(K)); i++ {
		if candidates.contacts[i].ID != kademlia.RoutingTable.Me.ID {
			kademlia.Network.SendStoreMessage(&candidates.contacts[i], key, data)
		} else {
			kademlia.Storage.Store(key.String(), data)
		}
	}
}

func (node *Kademlia) Bootstrap(contact Contact) {
	node.RoutingTable.AddContact(contact)
	node.LookupContact(node.RoutingTable.Me.ID)
}
