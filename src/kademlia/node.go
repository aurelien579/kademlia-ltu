package kademlia

import (
    "strconv"
)

type Node struct {
    routingTable *RoutingTable
}

var Me Node

func NewNode(id string, ip string, port int) Node {
    me := NewContact(NewKademliaID(id), ip + strconv.Itoa(port))
    node := Node{
        routingTable: NewRoutingTable(me),
    }
    
    return node
}

