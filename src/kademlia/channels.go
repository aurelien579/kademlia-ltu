package kademlia

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type ResponseChannel struct {
	Contact *Contact
	ReqType uint8
	Channel chan Header
}

type ChannelList struct {
	mutex sync.Mutex
	*list.List
}

func NewChannelList() *ChannelList {
	return &ChannelList{
		List: list.New(),
	}
}

func (list *ChannelList) Add(contact *Contact, reqType uint8, c chan Header) {
	channel := &ResponseChannel{
		Contact: contact,
		ReqType: reqType,
		Channel: c,
	}

	list.mutex.Lock()
	list.PushBack(channel)
	list.mutex.Unlock()
}

func (list *ChannelList) Find(contact *Contact, reqType uint8) (*ResponseChannel, error) {
	list.mutex.Lock()
	for e := list.Front(); e != nil; e = e.Next() {
		channel := e.Value.(*ResponseChannel)

		if *(channel.Contact.ID) == *contact.ID &&
			channel.ReqType == reqType {
			list.mutex.Unlock()
			return channel, nil
		}
	}
	list.mutex.Unlock()

	return nil, errors.New("Can't find channel")
}

func (list *ChannelList) Delete(contact *Contact, reqType uint8) {
	list.mutex.Lock()
	for e := list.Front(); e != nil; e = e.Next() {
		channel := e.Value.(*ResponseChannel)

		if *channel.Contact.ID == *contact.ID &&
			channel.ReqType == reqType {

			list.Remove(e)
		}
	}
	list.mutex.Unlock()
}

func (list *ChannelList) Send(header *Header) {
	contact := ContactFromHeader(header)
	channel, err := list.Find(&contact, header.SubType)
	if err != nil {
		fmt.Println(err)
		return
	}

	channel.Channel <- *header

	list.Delete(&contact, header.SubType)
}
