package kademlia

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      *Network
	Storage      *Storage
	Channels     *ChannelList
}

func NewKademlia(id string, ip string, port int) Kademlia {
	me := NewContact(NewKademliaID(id), ip+":"+strconv.Itoa(port))
	kademlia := Kademlia{
		RoutingTable: NewRoutingTable(me),
		Network:      NewNetwork(me.ID, ip, port),
		Storage:      NewStorage(id),
		Channels:     NewChannelList(),
	}

	log.SetFlags(0)
	//log.SetFlags(log.Ltime | log.Lmicroseconds)

	log.Println("MyID: ", id)

	kademlia.Network.Kademlia = &kademlia
	kademlia.Storage.kademlia = &kademlia

	for i := 0; i < IDLength*8; i++ {
		kademlia.RoutingTable.Buckets[i].node = &kademlia
	}

	return kademlia
}

func (kademlia *Kademlia) Listen(ip string, port int) {
	addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println(conn.LocalAddr())

	for {
		var header Header
		Decode(conn, &header)

		kademlia.RoutingTable.AddContact(ContactFromHeader(&header))

		switch header.SubType {
		case MSG_PING:
			go kademlia.HandlePing(header)
		case MSG_FIND_NODES:
			go kademlia.HandleFindNodes(header)
		case MSG_FIND_VALUE:
			go kademlia.HandleFindValue(header)
		case MSG_STORE:
			go kademlia.HandleStore(header)
		case MSG_PIN:
			go kademlia.HandlePin(header)
		case MSG_UNPIN:
			go kademlia.HandleUnpin(header)
		}
	}
}

func (kademlia *Kademlia) HandlePing(header Header) {
	if header.Type == MSG_REQUEST {
		conn, err := header.createConnection()

		if err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(1 * time.Second)

		Encode(conn, NewHeader(kademlia.Network, MSG_RESPONSE, MSG_PING))
	} else {
		kademlia.Channels.Send(&header)
	}
}

func (kademlia *Kademlia) SendContacts(header Header, contacts []Contact) {
	conn, err := header.createConnection()

	if err != nil {
		fmt.Println(err)
		return
	}

	header = NewHeader(kademlia.Network, MSG_RESPONSE, header.SubType)
	header.Arg = ContactsToContactResults(contacts)

	Encode(conn, header)
}

func (kademlia *Kademlia) HandleFindNodes(header Header) {
	switch header.Type {
	case MSG_REQUEST:
		findArguments := header.Arg.(FindArguments)

		contacts := kademlia.RoutingTable.FindClosestContacts(NewKademliaID(findArguments.Key), int(findArguments.Count))
		kademlia.SendContacts(header, contacts)
	case MSG_RESPONSE:
		kademlia.Channels.Send(&header)
	}
}

func (kademlia *Kademlia) SendFile(header Header, filename string) {
	conn, err := header.createConnection()

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := NewHeader(kademlia.Network, MSG_RESPONSE, MSG_FIND_VALUE)

	data := kademlia.Storage.Read(filename)

	msg.Arg = StoreArguments{
		Key:    filename,
		Length: len(data),
		Data:   data,
	}

	Encode(conn, msg)

	conn.Close()
}

func (kademlia *Kademlia) HandleFindValue(header Header) {
	switch header.Type {
	case MSG_REQUEST:
		findArguments := header.Arg.(FindArguments)
		key := findArguments.Key

		if kademlia.Storage.Exists(key) {
			kademlia.SendFile(header, key)
		} else {
			var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(NewKademliaID(key), int(findArguments.Count))
			kademlia.SendContacts(header, contacts)
		}

	case MSG_RESPONSE:
		kademlia.Channels.Send(&header)
	}
}

func (kademlia *Kademlia) HandleStore(header Header) {
	log.Println("HANDLER STORE: ", header)

	if header.Type != MSG_REQUEST {
		return
	}

	args := header.Arg.(StoreArguments)

	kademlia.Storage.Store(args.Key, args.Data, false)
}

func (kademlia *Kademlia) HandlePin(header Header) {
	log.Println("HANDLER PIN: ", header)

	if header.Type != MSG_REQUEST {
		return
	}

	args := header.Arg.(StoreArguments)

	kademlia.Storage.Store(args.Key, args.Data, true)
}

func (kademlia *Kademlia) HandleUnpin(header Header) {
	log.Println("HANDLER UNPIN: ", header)

	if header.Type != MSG_REQUEST {
		return
	}

	args := header.Arg.(StoreArguments)

	kademlia.Storage.Unpin(args.Key)
}

type ThreadContext struct {
	wg         sync.WaitGroup
	candidates *ContactCandidates
	target     *KademliaID
	done       bool
	lookupData bool
	buffer     bytes.Buffer
	mutex      sync.Mutex
}

func (kademlia *Kademlia) lookupThread(i int, context *ThreadContext) {
	log.Printf("=====================  Thread %d starting  =====================\n\n", i)
	log.Println("Initial candidates: ", context.candidates)

	for !context.done {
		target, err := context.candidates.GetClosestUnTook()

		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Thread %d target: %v\n", i, target)

		var subType = MSG_FIND_NODES
		if context.lookupData {
			subType = MSG_FIND_VALUE
		}

		channel := make(chan Header)
		kademlia.Channels.Add(target, subType, channel)

		if context.lookupData {
			kademlia.Network.SendFindDataMessage(target, context.target)
		} else {
			kademlia.Network.SendFindContactMessage(target, context.target)
		}

		header := <-channel

		isContacts := true
		if reflect.TypeOf(header.Arg).String() != "[]kademlia.ContactResult" {
			isContacts = false
		}

		target.State = DONE

		if isContacts {
			contacts := ContactResultsToContacts(header.Arg.([]ContactResult))
			CalcDistances(contacts, context.target)

			log.Printf("Thread %d results: [\n", i)
			for _, r := range contacts {
				log.Println("\t", r)
			}
			log.Println("]")

			context.candidates.AppendSorted(contacts)
		} else {
			context.done = true

			context.mutex.Lock()
			context.buffer.Write(header.Arg.(StoreArguments).Data)
			context.mutex.Unlock()
		}

		if !context.lookupData {
			if context.candidates.Finish(int(K)) {
				context.done = true
			}
		}

		close(channel)
	}

	context.wg.Done()

	log.Printf("=====================  Thread %d done  =====================\n", i)
}

func (kademlia *Kademlia) Lookup(key *KademliaID, lookupData bool) (error, []Contact, []byte) {
	log.Printf("Starting lookup on %s\n", key.String())

	context := ThreadContext{
		done:       false,
		lookupData: lookupData,
		target:     key,
		candidates: &ContactCandidates{},
	}

	closestKnown := kademlia.RoutingTable.FindClosestContacts(key, int(K))

	if len(closestKnown) == 0 {
		return errors.New("Can't do lookup without knowing anyone"), nil, nil
	}

	me := NewContact(kademlia.RoutingTable.Me.ID, kademlia.RoutingTable.Me.Address)
	me.State = DONE

	closestKnown = append(closestKnown, me)
	CalcDistances(closestKnown, key)

	context.candidates.AppendSorted(closestKnown)

	context.wg.Add(ALPHA)
	for i := 0; i < ALPHA; i++ {
		go kademlia.lookupThread(i, &context)
	}
	context.wg.Wait()

	if context.lookupData {
		log.Println("Lookup data done")
		return nil, nil, context.buffer.Bytes()
	} else {
		log.Printf("Lookup contacts done. Results: \n%v\n", context.candidates)
		return nil, context.candidates.GetContacts(K), nil
	}
}

func (kademlia *Kademlia) LookupContact(key *KademliaID) (error, []Contact) {
	err, contacts, _ := kademlia.Lookup(key, false)
	return err, contacts
}

func (kademlia *Kademlia) LookupData(hash string) (error, []byte) {
	key := NewKademliaID(hash)
	err, _, data := kademlia.Lookup(key, true)
	return err, data
}

func (kademlia *Kademlia) Store(data []byte) (error, string) {
	key := NewKademliaID(HashBytes(data))

	err, contacts := kademlia.LookupContact(key)

	if err != nil {
		return err, ""
	}

	for i := 0; i < len(contacts); i++ {
		if contacts[i].ID != kademlia.RoutingTable.Me.ID {
			kademlia.Network.SendStoreMessage(&contacts[i], key, data)
		} else {
			kademlia.Storage.Store(key.String(), data, false)
		}
	}

	return nil, key.String()
}

func (node *Kademlia) Bootstrap(contact Contact) {
	node.RoutingTable.AddContact(contact)
	node.LookupContact(node.RoutingTable.Me.ID)
}
