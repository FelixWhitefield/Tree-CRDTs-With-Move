package connection

import (
	"net"
)

type TCPConnection struct {

}

func NewTCPConnection[OP []byte, OPID comparable](conn net.Conn, p *TCPProvider[OP, OPID]) *TCPConnection {
	return &TCPConnection{}
}

func (c *TCPConnection) handle() {

}