package daemon

import (
	"errors"
	"net"
)

const CMD_GET = 1
const CMD_PUT = 2

type Command struct {
	command int
	arg     string
}

const ERROR = -1
const OK = 1

type Response struct {
	resultCode int
	result     string
}

func SendCommand(conn *net.UDPConn, command int, arg string) {

}

func ReadCommand(conn *net.UDPConn) (Command, error) {
	return Command{}, errors.New("")
}

func SendResponse(conn *net.UDPConn, code int, result string) {

}

func ReadResponse(conn *net.UDPConn) (Response, error) {
	return Response{}, errors.New("")
}
