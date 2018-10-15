package main

import (
	"kademlia"
	"net"
	"strings"
)

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
	ip := getMyIp()

	node = kademlia.NewKademlia("0000000000000000000000000000000000000001", ip, port)

	go node.Listen(ip, port)

	for {

	}
}
