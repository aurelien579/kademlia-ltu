package daemon

import (
	"bytes"
	"encoding/gob"
	"net"
)

const CMD_GET = 1
const CMD_PUT = 2
const CMD_UNPIN = 3
const CMD_PIN = 4

type Command struct {
	Command int
	Arg     string
}

const ERROR = -1
const OK = 1

type Response struct {
	ResultCode int
	Result     string
}

func Encode(c *net.UDPConn, addr *net.UDPAddr, value interface{}) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(value)

	if err != nil {
		return err
	}

	if addr == nil {
		c.Write(buffer.Bytes())
	} else {
		c.WriteToUDP(buffer.Bytes(), addr)
	}

	return nil
}

func Decode(c *net.UDPConn, value interface{}) (*net.UDPAddr, error) {
	inputBytes := make([]byte, 1024)

	length, addr, err := c.ReadFromUDP(inputBytes)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(inputBytes[:length])
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(value)

	if err != nil {
		return nil, err
	}

	return addr, nil
}

func SendCommand(conn *net.UDPConn, addr *net.UDPAddr, command int, arg string) {
	Encode(conn, addr, Command{
		Command: command,
		Arg:     arg,
	})
}

func ReadCommand(conn *net.UDPConn) (*Command, *net.UDPAddr, error) {
	var command Command

	addr, err := Decode(conn, &command)
	if err != nil {
		return nil, nil, err
	}

	return &command, addr, nil
}

func SendResponse(conn *net.UDPConn, addr *net.UDPAddr, code int, result string) error {
	return Encode(conn, addr, Response{
		ResultCode: code,
		Result:     result,
	})
}

func ReadResponse(conn *net.UDPConn) (*Response, error) {
	var response Response

	_, err := Decode(conn, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
