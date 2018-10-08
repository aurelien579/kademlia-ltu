package main

import (
	"daemon"
	"errors"
	"fmt"
	"net"
	"os"
)

func parseCommand(str string) (int, error) {
	if str == "get" {
		return daemon.CMD_GET, nil
	} else if str == "put" {
		return daemon.CMD_PUT, nil
	} else {
		return -1, errors.New("Invalid command")
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <command> <filename>\n", os.Args[0])
		return
	}

	command, err := parseCommand(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid command entered: %s\n", os.Args[1])
		fmt.Printf("Valid commands are: get, put\n")
		return
	}

	arg := os.Args[2]

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:40000")
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	daemon.SendCommand(conn, command, arg)

	response, err := daemon.ReadResponse(conn)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("%d, %s\n", response.ResultCode, response.Result)
}
