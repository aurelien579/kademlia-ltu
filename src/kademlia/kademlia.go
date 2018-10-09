package kademlia

import (
	"bytes"
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
	Network      Network
	Storage      Storage
	Channels     ChannelList
}

func NewKademlia(id string, ip string, port int) Kademlia {
	me := NewContact(NewKademliaID(id), ip+":"+strconv.Itoa(port))
	kademlia := Kademlia{
		RoutingTable: NewRoutingTable(me),
		Network:      NewNetwork(me.ID, ip, port),
		Storage:      NewStorage(id),
	}

	log.Println("MyID: ", id)

	kademlia.Network.Kademlia = &kademlia

	kademlia.Storage.kademlia = &kademlia

	kademlia.Channels = NewChannelList()

	for i := 0; i < IDLength*8; i++ {
		kademlia.RoutingTable.Buckets[i].node = &kademlia
	}

	log.SetFlags(log.Ltime | log.Lmicroseconds)

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

		Encode(conn, NewHeader(&kademlia.Network, MSG_RESPONSE, MSG_PING))
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

func (kademlia *Kademlia) HandleFindNodes(header Header) {
	switch header.Type {
	case MSG_REQUEST:
		findArguments := header.Arg.(FindArguments)

		var contacts []Contact = kademlia.RoutingTable.FindClosestContacts(&(findArguments.Key), int(findArguments.Count))
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

func (kademlia *Kademlia) HandleFindValue(header Header) {
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
		kademlia.Channels.Send(&header)
	}
}

func (kademlia *Kademlia) HandleStore(header Header) {
	log.Println("HANDLER STORE: ", header)

	if header.Type != MSG_REQUEST {
		return
	}

	args := header.Arg.(StoreArguments)

	kademlia.Storage.Store(args.Key.String(), args.Data)
}

type ThreadContext struct {
	wg         sync.WaitGroup
	candidates ContactCandidates
	target     *KademliaID
	done       bool
	lookupData bool
	buffer     bytes.Buffer
	mutex      sync.Mutex
}

func (kademlia *Kademlia) lookupThread(i int, context *ThreadContext) {
	log.Printf("\n\n\n\n============ Thread %d starting ============\n\n", i)
	log.Println("Initial candidates: ", context.candidates)

	for !context.done {
		target, err := context.candidates.GetClosestUnTook(int(K))

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
			results := header.Arg.([]ContactResult)
			var contacts []Contact

			for _, r := range results {
				contacts = append(contacts, NewContact(NewKademliaID(r.ID), r.Address))
			}

			for i := 0; i < len(contacts); i++ {
				contacts[i].CalcDistance(context.target)
			}

			log.Printf("Thread %d results: [\n", i)
			for _, r := range contacts {
				log.Println("\t", r)
			}
			log.Println("]")

			context.candidates.Append(contacts)
			context.candidates.Sort()
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

	log.Printf("============ Thread %d done ============\n", i)
}

func (kademlia *Kademlia) Lookup(key *KademliaID, lookupData bool) interface{} {
	context := ThreadContext{
		done:       false,
		lookupData: lookupData,
		target:     key,
	}

	closestKnown := kademlia.RoutingTable.FindClosestContacts(key, int(K))

	me := NewContact(kademlia.RoutingTable.Me.ID, kademlia.RoutingTable.Me.Address)
	me.State = DONE

	closestKnown = append(closestKnown, me)
	context.candidates.Append(closestKnown)

	for i := 0; i < context.candidates.Len(); i++ {
		context.candidates.contacts[i].CalcDistance(key)
	}
	context.candidates.Sort()

	context.wg.Add(ALPHA)
	for i := 0; i < ALPHA; i++ {
		go kademlia.lookupThread(i, &context)
	}
	context.wg.Wait()

	if context.lookupData {
		return context.buffer.Bytes()
	} else {
		return context.candidates
	}
}

func (kademlia *Kademlia) LookupContact(key *KademliaID) ContactCandidates {
	return kademlia.Lookup(key, false).(ContactCandidates)
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	key := NewKademliaID(hash)
	return kademlia.Lookup(key, true).([]byte)
}

func (kademlia *Kademlia) Store(data []byte) string {
	key := NewKademliaID(HashBytes(data))

	candidates := kademlia.LookupContact(key)

	for i := 0; i < Min(candidates.Len(), int(K)); i++ {
		if candidates.contacts[i].ID != kademlia.RoutingTable.Me.ID {
			kademlia.Network.SendStoreMessage(&candidates.contacts[i], key, data)
		} else {
			kademlia.Storage.Store(key.String(), data)
		}
	}

	return key.String()
}

func (node *Kademlia) Bootstrap(contact Contact) {
	node.RoutingTable.AddContact(contact)
	node.LookupContact(node.RoutingTable.Me.ID)
}
