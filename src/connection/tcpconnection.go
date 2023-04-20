package connection

import (
	"net"
)

type TCPConnection[OPID comparable] struct {
	conn net.Conn
}

func NewTCPConnection[OP []byte, OPID comparable](conn net.Conn, p *TCPProvider[OP, OPID]) *TCPConnection[OPID] {
	return &TCPConnection[OPID]{}
}

func (c *TCPConnection[OPID]) handle() {

}

func (c *TCPConnection[OPID]) sendOp(opid OPID, op []byte) {
	c.conn.Write(op)
}