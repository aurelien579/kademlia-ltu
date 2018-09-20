package kademlia

type Network struct {
    
}

const MSG_REQUEST uint8 = 1
const MSG_RESPONSE uint8 = 2

const MSG_PING uint8 = 1
const MSG_FIND_NODES uint8 = 2
const MSG_FIND_VALUE uint8 = 3
const MSG_STORE uint8 = 4

type Header struct {
    SrcID kademliaID
    SrcIP uint32
    SrcPort uint8
    Type uint8      /* Request/Response */
    SubType uint8
}

func (network *Network) SendPingMessage(contact *Contact) {	
	
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

