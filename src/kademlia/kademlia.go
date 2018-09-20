package kademlia

import (
    "strconv"
)

type Kademlia struct {
    routingTable *RoutingTable
    network Network
}

var Me Kademlia

func KademliaInit(id string, ip string, port int) {
    me := NewContact(NewKademliaID(id), ip + strconv.Itoa(port))
    Me := Kademlia{
        routingTable: NewRoutingTable(me),
    }
    
    return node
}

func (kademlia *Kademlia) Listen(ip string, port int) {
    udpAddr, _ := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(port))
    
    for {
        udpConn, err := net.ListenUDP("udp", udpAddr)
        
        if err != nil {
            fmt.Println(err)
            return
        }
        
        inputBytes := make([]byte, 1024)
        length, _ := c.Read(inputBytes)
        buf := bytes.NewBuffer(inputBytes[:length])
        
        decoder := gob.NewDecoder(buf)
        var header Header
        decoder.Decode(&header)
        
        // TODO: gérer le message dans un nouveau thread + AddContact
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
