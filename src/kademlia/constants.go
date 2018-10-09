package kademlia

/* Kademlia constants */

const REPUBLISH_TIME = 30
const K = 2
const ALPHA = 1

/* Protocol constants */

const MSG_REQUEST uint8 = 1
const MSG_RESPONSE uint8 = 2

const MSG_PING uint8 = 1
const MSG_FIND_NODES uint8 = 2
const MSG_FIND_VALUE uint8 = 3
const MSG_STORE uint8 = 4
