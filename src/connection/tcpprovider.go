package connection

// 
// 
// When a peer connects, they should exchange their ID (This should always be the first message sent)
// And then after that, they should exchange their peers

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

const (
	Op = "OP"
	OpAck = "OPACK"
	Join = "JOIN"
	Peers = "PEERS"
)

type Operation struct {
	op []byte
}

type InitialMessage struct {
	id uuid.UUID
}

type PeersMessage struct {
	peers map[uuid.UUID]net.TCPAddr
}

type TCPProvider struct {
	id 		   	   uuid.UUID
	numPeers       int
	peersMu        sync.RWMutex
	peers          map[uuid.UUID]*TCPConnection // peerID -> connection
	deliveredMu    sync.RWMutex
	delivered      map[uuid.UUID][]uuid.UUID // opID -> list of peerIDs that have delivered the op
	operations     map[uuid.UUID][]byte		 // opID -> op
	incomingOps    chan []byte
	opsToBroadcast chan Operation
}

func NewTCPProvider[OP []byte, PID comparable](numPeers int, id uuid.UUID) *TCPProvider {
	return &TCPProvider{numPeers: numPeers, id: id}
}

func (p *TCPProvider) Listen(port int) {
	address := net.JoinHostPort("::", strconv.Itoa(port))
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatalf("Error resolving address: %s", err.Error())
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("Error listening: %s", err.Error())
	}
	defer listener.Close()

	log.Println("Listening on:", listener.Addr().String())

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
		}
		log.Println("Accepted connection from:", conn.RemoteAddr())

		go NewTCPConnection(conn, p).handle()
	}
}

func (p *TCPProvider) Connect(addr string) {
	tcp, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcp)
	if err != nil {
		panic(err)
	}

	go NewTCPConnection(conn, p).handle()
}

// Broadcasts an operation to all peers
func (p *TCPProvider) handleBroadcast() {
	for {
		opToSend := <-p.opsToBroadcast

		// Generate a new ID for the operation
		newOpId := uuid.Must(uuid.NewUUID())
		p.AddOperation(opToSend.op, newOpId)
		
		opMsg := OperationMsg{Id: newOpId[:], Op: opToSend.op}
		opData, err := proto.Marshal(&opMsg)
		if err != nil {
			log.Println("Error marshalling operation: ", err.Error())
			continue
		}

		p.peersMu.RLock()
		for _, conn := range p.peers {
			conn.SendMsg(opData)
		}
		p.peersMu.RUnlock()
	}
}

func (p *TCPProvider) AddPeer(id uuid.UUID, conn *TCPConnection) {
	p.peersMu.Lock()
	defer p.peersMu.Unlock()

	p.peers[id] = conn
}

func (p *TCPProvider) RemovePeer(id uuid.UUID) {
	p.peersMu.Lock()
	defer p.peersMu.Unlock()

	p.peers[id] = nil 
}

func (p *TCPProvider) AddOperation(op []byte, id uuid.UUID) {
	p.deliveredMu.Lock()
	defer p.deliveredMu.Unlock()

	p.operations[id] = op
	p.delivered[id] = make([]uuid.UUID, 0, p.numPeers-1)
}

func (p *TCPProvider) AddDelivered(id uuid.UUID, peer uuid.UUID) {
	p.deliveredMu.Lock()
	defer p.deliveredMu.Unlock()

	if len(p.delivered[id]) == p.numPeers-1 {
		delete(p.operations, id)
		delete(p.delivered, id)
	} else {
		p.delivered[id] = append(p.delivered[id], peer)
	}
}

func (p *TCPProvider) GetOperation(id uuid.UUID) []byte {
	p.deliveredMu.RLock()
	defer p.deliveredMu.RUnlock()

	return p.operations[id]
}
