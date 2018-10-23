package kademlia

import (
	"fmt"
	"sync"
	"testing"
)

func TestContact(t *testing.T) {

	c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1:4000")
	c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.2:4001")
	c3 := NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.3:4002")

	c4 := NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "127.0.0.4:4003")

	c1.CalcDistance(c3.ID)
	c2.CalcDistance(c3.ID)

	fmt.Printf("c1.distance : %s\n", c1.distance)
	fmt.Printf("c2.distance : %s\n", c2.distance)

	b1 := c1.Closer(&c2)
	b2 := c2.Closer(&c1)

	fmt.Printf("c1 closer than c2 : %t\n", b1)
	fmt.Printf("c2 closer than c1 : %t\n", b2)

	s := c4.String()

	fmt.Printf("c4.string : %s\n", s)

	tab := []Contact{c1, c2, c3}

	CalcDistances(tab, c4.ID)

	fmt.Printf("c1 distance : %s\n", tab[0].distance)
	fmt.Printf("c2 distance : %s\n", tab[1].distance)
	fmt.Printf("c3 distance : %s\n\n\n", tab[2].distance)

}

func TestContactCandidates(t *testing.T) {

	c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1:4000")
	c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.2:4001")
	c3 := NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.3:4002")

	c4 := NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "127.0.0.4:4003")
	c5 := NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "127.0.0.5:4004")

	c6 := NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "127.0.0.6:4005")

	tab := []Contact{c1, c2, c3}

	tab2 := []Contact{c4, c5}

	var mutex = sync.Mutex{}

	cc := ContactCandidates{tab, mutex}

	s := cc.String()

	fmt.Printf("cc string1 : %s\n", s)

	//function Append is tested in AppendSorted
	//function Sort is tested in AppendSorted
	cc.Append(tab2)

	s = cc.String()

	fmt.Printf("cc string2 : %s\n", s)

	b1 := c5.IsIn(cc)
	b2 := c6.IsIn(cc)

	fmt.Printf("c5 Is In cc : %t\n", b1)
	fmt.Printf("c6 Is In cc : %t\n", b2)

	var tab3 []Contact

	tab3 = cc.GetContacts(3)

	fmt.Printf("get 3 contacts : %s\n", tab3[0].String())
	fmt.Printf("get 3 contacts : %s\n", tab3[1].String())
	fmt.Printf("get 3 contacts : %s\n", tab3[2].String())

	l := cc.Len()

	fmt.Printf("cc len : %d\n", l)

	cc.Swap(1, 2)

	fmt.Printf("cc swap indice 1 2 : %s\n", cc.String())

	//cc.contacts[0].CalcDistance(c6.ID)
	CalcDistances(cc.contacts, c6.ID)
	s = cc.String()

	fmt.Printf("cc string3 : %s\n", s)

	b3 := cc.Less(1, 2)

	fmt.Printf("cc less 1 2 : %t\n", b3)

	b1 = cc.Finish(4)

	fmt.Printf("cc finish when every one UNDONE : %t\n", b1)

	cc.contacts[0].State = DONE
	cc.contacts[1].State = DONE
	cc.contacts[2].State = DONE
	cc.contacts[3].State = DONE
	cc.contacts[4].State = DONE
	b2 = cc.Finish(4)

	fmt.Printf("cc finish when every one DONE : %t\n", b2)

	cc.contacts[2].State = UNDONE
	cc.contacts[3].State = UNDONE
	cc.contacts[4].State = UNDONE

	c, _ := cc.GetClosestUnTook()

	fmt.Printf("cc get closest untook when 3 4 5 UNDONE : %s\n", c)

	fmt.Printf("cc when one is taken : %s\n", cc.String())

}
