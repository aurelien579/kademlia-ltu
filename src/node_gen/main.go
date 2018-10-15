package main

import (
	"fmt"
	"kademlia"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const MY_ID = "000000000000000000000000000000000000FFFF"
const MY_IP = "127.0.0.1"
const MY_PORT = 3333

func getFreePort(start int) int {
	for {
		if start >= 65555 {
			return 0
		}

		addr, _ := net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(start))
		conn, err := net.ListenUDP("udp", addr)

		if err != nil {
			start++
		} else {
			conn.Close()
			return start
		}
	}
}

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

func findContact(name string) string {
	addrs, err := net.LookupHost(name)
	if err != nil {
		time.Sleep(1 * time.Second)
		fmt.Println("Second attempt...")
		addrs, err = net.LookupHost(name)

		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	return addrs[0]
}

func main() {
	var node kademlia.Kademlia
	port := 4000
	ip := getMyIp()

	contactIp := findContact("bootstrap_node")

	node = kademlia.NewKademlia("000000000000000000000000000000000000000"+string(ip[len(ip)-1]), ip, port)

	go node.Listen(ip, port)

	node.Bootstrap(kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), contactIp+":4000"))

	for {

	}
}
