package connection

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"sync"
)

type TCPProvider[OP []byte, OPID comparable] struct {
	peersMu        sync.RWMutex
	peers          map[uuid.UUID]*TCPConnection
	numPeers       int
	deliveredMu    sync.RWMutex
	delivered      map[OPID][]uuid.UUID
	operations     map[OPID]OP
	incomingOps    chan []byte
	opsToBroadcast chan []byte
}

func NewTCPProvider[OP []byte, OPID comparable](numPeers int) *TCPProvider[OP, OPID] {
	return &TCPProvider[OP, OPID]{numPeers: numPeers}
}

func (p *TCPProvider[OP, OPID]) Listen(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	log.Println("Listening on port:", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
		}

		go NewTCPConnection(conn, p).handle()
	}
}

func (p *TCPProvider[OP, OPID]) AddPeer(id uuid.UUID, conn *TCPConnection) {
	p.peersMu.Lock()
	defer p.peersMu.Unlock()

	p.peers[id] = conn
}

func (p *TCPProvider[OP, OPID]) RemovePeer(id uuid.UUID) {
	p.peersMu.Lock()
	defer p.peersMu.Unlock()

	delete(p.peers, id)
}

func (p *TCPProvider[OP, OPID]) AddOperation(op OP, id OPID) {
	p.deliveredMu.Lock()
	defer p.deliveredMu.Unlock()

	p.operations[id] = op
	p.delivered[id] = make([]uuid.UUID, 0, p.numPeers-1)
}

func (p *TCPProvider[OP, OPID]) AddDelivered(id OPID, peer uuid.UUID) {
	p.deliveredMu.Lock()
	defer p.deliveredMu.Unlock()

	if len(p.delivered[id]) == p.numPeers-1 {
		delete(p.operations, id)
		delete(p.delivered, id)
	} else {
		p.delivered[id] = append(p.delivered[id], peer)
	}
}

func (p *TCPProvider[OP, OPID]) GetOperation(id OPID) OP {
	p.deliveredMu.RLock()
	defer p.deliveredMu.RUnlock()

	return p.operations[id]
}
