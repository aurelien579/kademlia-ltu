package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
)

type Kademlia struct {
    RoutingTable *RoutingTable
    Network Network
}

func NewKademlia(id string, ip string, port int) Kademlia {
    me := NewContact(NewKademliaID(id), ip + strconv.Itoa(port))
    kademlia := Kademlia{
        RoutingTable: NewRoutingTable(me),
        Network: NewNetwork(me.ID, ip, port),
    }
    
    return kademlia
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
        length, _ := udpConn.Read(inputBytes)
        buf := bytes.NewBuffer(inputBytes[:length])
        
        decoder := gob.NewDecoder(buf)
        var header Header
        decoder.Decode(&header)

        switch header.SubType{

		case MSG_PING :
			//traiter le msgping
		    kademlia.RoutingTable.AddContact(NewContact(&(header.SrcID),IPToStr(header.SrcIP)))

			addr, _ := net.ResolveUDPAddr("udp", IPToStr(header.SrcIP))
			conn, err := net.DialUDP("udp", nil, addr)

			if err != nil {
				fmt.Println(err)
				return
			}

			var buffer bytes.Buffer
			enc := gob.NewEncoder(&buffer)
			msg := Header{
				SrcID: *kademlia.Network.ID,
				SrcIP: kademlia.Network.IP,
				SrcPort: kademlia.Network.Port,
				Type: MSG_RESPONSE,
				SubType: MSG_PING,
			}

			enc.Encode(msg)
			conn.Write(buffer.Bytes())



		case MSG_FIND_NODES :

		case MSG_FIND_VALUE :

		case MSG_STORE :


		}


        
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
