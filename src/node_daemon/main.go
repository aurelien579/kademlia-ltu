package main

import (
	"daemon"
	"io/ioutil"
	"kademlia"
	"log"
	"net"
	"strconv"
	"strings"
)

const MY_ID = "000000000000000000000000000000000000FFFF"
const MY_IP = "127.0.0.1"
const MY_PORT = 3333

func getMyIp() string {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if strings.Contains(i.Name, "eth0") {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				return ip.String()
			}
		}
	}

	return ""
}

func main() {
	var node kademlia.Kademlia
	port := 4000
	ip := "127.0.0.1"

	log.Println(ip)

	node = kademlia.NewKademlia("000000000000000000000000000000000000000"+string(ip[len(ip)-1]), ip, port)

	go node.Listen(ip, port)


	for {
		ListenDaemon(&node, 40000)
	}
	
}

func ListenDaemon(node *kademlia.Kademlia, port int) {

	udpAddr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))

	conn, _ := net.ListenUDP("udp", udpAddr)

	command, _ := daemon.ReadCommand(conn)

	log.Printf("Command received : %v\n", command)

	ExecuteCommand(node, command, conn)

}

func ExecuteCommand(node *kademlia.Kademlia, command *daemon.Command, conn *net.UDPConn) {

	switch command.Command {

	case daemon.CMD_GET:

		log.Printf("Launching LookupData : %v\n", command.Arg)

		data := node.LookupData(command.Arg)

		s := string(data[:])

		log.Printf("Response received, telling the daemon\n")

		daemon.SendResponse(conn, daemon.OK, s)

	case daemon.CMD_PUT:

		bytes, _ := ioutil.ReadFile(command.Arg)

		log.Printf("Launching Store\n")

		hash := node.Store(bytes)

		log.Printf("Response received, telling the daemon\n")

		daemon.SendResponse(conn, daemon.OK, hash)

	}

}
