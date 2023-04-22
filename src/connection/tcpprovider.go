package connection

//
//
// When a peer connects, they should exchange their ID (This should always be the first message sent)
// And then after that, they should exchange their peers

import (
	"errors"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Operation struct {
	op []byte
}

type TCPProvider struct {
	id             uuid.UUID
	numPeers       int
	peersMu        sync.RWMutex
	peers          map[uuid.UUID]*TCPConnection // peerID -> connection : When deleted, will set value to nil (The total peer set should not change)
	deliveredMu    sync.RWMutex                 // Mutex for the delivered map (Also locked when accessing the operations map, removes need for a separate mutex)
	delivered      map[uuid.UUID][]uuid.UUID    // opID -> list of peerIDs that have been delivered the op + acked
	operations     map[uuid.UUID][]byte         // opID -> op
	incomingOps    chan []byte
	opsToBroadcast chan Operation
}

func NewTCPProvider(numPeers int, id uuid.UUID) *TCPProvider {
	return &TCPProvider{
		numPeers:       numPeers,
		id:             id,
		peers:          make(map[uuid.UUID]*TCPConnection, numPeers),
		delivered:      make(map[uuid.UUID][]uuid.UUID),
		operations:     make(map[uuid.UUID][]byte),
		incomingOps:    make(chan []byte, 10),
		opsToBroadcast: make(chan Operation, 10),
	}
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

// Errors if the peer already exists or the peer map is full
func (p *TCPProvider) AddPeer(id uuid.UUID, tcpConn *TCPConnection) error {
	p.peersMu.Lock()
	defer p.peersMu.Unlock()

	if val, ok := p.peers[id]; ok && val != nil { // Check if the peer already exists
		return errors.New("peer already exists (not nil)")
	} else if len(p.peers) == p.numPeers { // Check if the peer map is full
		return errors.New("peer map is full")
	}

	p.peers[id] = tcpConn
	return nil
}

// Sets the peer in map to nil
func (p *TCPProvider) RemovePeer(id uuid.UUID) {
	p.peersMu.Lock()
	defer p.peersMu.Unlock()

	p.peers[id] = nil
}

func (p *TCPProvider) AddOperation(op []byte, id uuid.UUID) {
	p.deliveredMu.Lock()
	defer p.deliveredMu.Unlock()

	p.operations[id] = op
	p.delivered[id] = make([]uuid.UUID, 0, p.numPeers-1) // Size is numPeers-1 because final peer won't store the operation in the delivered map
}

func (p *TCPProvider) AddDelivered(id uuid.UUID, peer uuid.UUID) {
	p.deliveredMu.Lock()
	defer p.deliveredMu.Unlock()

	if len(p.delivered[id]) == p.numPeers-1 { // This would be the last peer to receive the operation (No need to add it to the delivered map)
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
