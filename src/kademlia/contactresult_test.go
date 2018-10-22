package kademlia

import(

  "testing"
  "fmt"
)


func TestContactResult (t *testing.T){

  c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"),"127.0.0.1:4000" )
  c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"),"127.0.0.2:4001" )
  c3 := NewContact(NewKademliaID("0000000000000000000000000000000000000003"),"127.0.0.3:4002" )

  cr := NewContactResult(&c1)

  fmt.Printf("c1 : %s\n", c1.String())
  fmt.Printf("cr : %s\n", cr)

  tab := []Contact{c2, c3}
  tabcr := ContactsToContactResults(tab)

  fmt.Printf("cr2 : %s\n", tabcr[0])
  fmt.Printf("cr3 : %s\n", tabcr[1])

  tabc := ContactResultsToContacts(tabcr)

  fmt.Printf("c2 : %s\n", tabc[0].String())
  fmt.Printf("c3 : %s\n", tabc[1].String())

}
