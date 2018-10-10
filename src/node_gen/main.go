package main

import (
	"fmt"
	"kademlia"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	fastping "github.com/tatsushid/go-fastping"
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

func findContact(ip string) string {
	p := fastping.NewPinger()

	splitted := strings.Split(ip, ".")
	mySuffix, _ := strconv.Atoi(splitted[3])
	prefix := splitted[0] + "." + splitted[1] + "." + splitted[2] + "."

	c := make(chan string)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		c <- addr.IP.String()
		p.Stop()
	}

	for i := 1; i < 255; i++ {
		if i == mySuffix {
			continue
		}

		ra, err := net.ResolveIPAddr("ip4:icmp", prefix+strconv.Itoa(i))
		fmt.Println("Adding:", ra.String())
		if err != nil {
			log.Println(err)
			continue
		}

		p.AddIPAddr(ra)
	}

	err := p.Run()
	if err != nil {
		fmt.Println(err)
	}

	return <-c
}

func main() {
	var node kademlia.Kademlia
	port := 4000
	ip := getMyIp()

	c := findContact("192.168.0.103")

	log.Println(c)
	return

	log.Println(ip)

	node = kademlia.NewKademlia("000000000000000000000000000000000000000"+string(ip[len(ip)-1]), ip, port)

	go node.Listen(ip, port)

	//time.Sleep(300 * time.Millisecond)

	if ip != "172.17.0.2" {
		node.Bootstrap(kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000002"), "172.17.0.2:4000"))
	}

	for {

	}
}
