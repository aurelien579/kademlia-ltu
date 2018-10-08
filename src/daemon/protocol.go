package daemon

import (
	"bytes"
	"encoding/gob"
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

func Encode(c *net.UDPConn, value interface{}) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(value)

	if err != nil {
		return err
	}

	c.Write(buffer.Bytes())

	return nil
}

func Decode(c *net.UDPConn, value interface{}) error {
	inputBytes := make([]byte, 1024)

	length, _, err := c.ReadFromUDP(inputBytes)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(inputBytes[:length])
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(value)

	if err != nil {
		return err
	}

	return nil
}

func SendCommand(conn *net.UDPConn, command int, arg string) {
	Encode(conn, Command{
		command: command,
		arg:     arg,
	})
}

func ReadCommand(conn *net.UDPConn) (*Command, error) {
	var command Command

	err := Decode(conn, &command)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func SendResponse(conn *net.UDPConn, code int, result string) {
	Encode(conn, Response{
		resultCode: code,
		result:     result,
	})
}

func ReadResponse(conn *net.UDPConn) (*Response, error) {
	var response Response

	err := Decode(conn, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
