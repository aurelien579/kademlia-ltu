package kademlia

type Network struct {
    ID KademliaID
    IP uint32
    Port uint16
}

const MSG_REQUEST uint8 = 1
const MSG_RESPONSE uint8 = 2

const MSG_PING uint8 = 1
const MSG_FIND_NODES uint8 = 2
const MSG_FIND_VALUE uint8 = 3
const MSG_STORE uint8 = 4

type Header struct {
    SrcID KademliaID
    SrcIP uint32
    SrcPort uint16
    Type uint8      /* Request/Response */
    SubType uint8
}

func IPToLong(ip string) uint32 {
    var long uint32
    binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
    return long
}

func IPToStr(ipInt uint32) string {
    // need to do two bit shifting and “0xff” masking
    b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
    b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
    b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
    b3 := strconv.FormatInt((ipInt & 0xff), 10)
    return b0 + "." + b1 + "." + b2 + "." + b3
}

func NewNetwork(id KademliaID, ip string, port int) Network {
    return Network{
        ID: id,
        IP: ip2Long(ip),
        Port: port,
    }
}

func (network *Network) SendPingMessage(contact *Contact) {
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, addr)
    
    if err != nil {
        fmt.Println(err)
        return
    }
    
    var buffer bytes.Buffer
    enc := gob.NewEncoder(&buffer)
    msg := Header{
        SrcID: network.ID,
        SrcIP: network.IP,
        SrcPort: network.Port,
        Type: MSG_REQUEST,
        SubType: MSG_PING,
    }
    
    enc.Encode(msg)
    conn.Write(buffer.Bytes())
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

