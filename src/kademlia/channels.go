package kademlia

import (
	"container/list"
	"errors"
	"fmt"
)

type ResponseChannel struct {
	Contact *Contact
	ReqType uint8
	Channel chan Header
}

type ChannelList struct {
	*list.List
}

func NewChannelList() ChannelList {
	return ChannelList{
		List: list.New(),
	}
}

func (list *ChannelList) Add(contact *Contact, reqType uint8, c chan Header) {
	channel := &ResponseChannel{
		Contact: contact,
		ReqType: reqType,
		Channel: c,
	}

	list.PushBack(channel)
}

func (list *ChannelList) Find(contact *Contact, reqType uint8) (*ResponseChannel, error) {
	for e := list.Front(); e != nil; e = e.Next() {
		channel := e.Value.(*ResponseChannel)

		if *(channel.Contact.ID) == *contact.ID &&
			channel.ReqType == reqType {
			return channel, nil
		}
	}

	return nil, errors.New("Can't find channel")
}

func (list *ChannelList) Delete(contact *Contact, reqType uint8) {
	for e := list.Front(); e != nil; e = e.Next() {
		channel := e.Value.(*ResponseChannel)

		if *channel.Contact.ID == *contact.ID &&
			channel.ReqType == reqType {

			list.Remove(e)
		}
	}
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
