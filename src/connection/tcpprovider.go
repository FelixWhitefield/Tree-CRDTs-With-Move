package connection

//
//
// When a peer connects, they should exchange their ID (This should always be the first message sent)
// And then after that, they should exchange their peers
//
// This will share peers with other peers. Whenever a new connection is made, peer addresses are shared with the new peer.

import (
	"errors"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type TCPProvider struct {
	port           int
	id             uuid.UUID
	numPeers       int
	peersMu        sync.RWMutex
	peers          map[uuid.UUID]*TCPConnection // peerID -> connection : When deleted, will set value to nil (The total peer set should not change)
	peerAddrs      map[net.Addr]bool            // Set of peer addresses
	deliveredMu    sync.RWMutex                 // Mutex for the delivered map (Also locked when accessing the operations map, removes need for a separate mutex)
	delivered      map[uuid.UUID]map[uuid.UUID]bool    // opID -> set of peerIDs that have been delivered the op + acked
	operations     map[uuid.UUID][]byte         // opID -> op
	incomingOps    chan []byte
	opsToBroadcast chan []byte
}

func NewTCPProvider(numPeers int, port int) *TCPProvider {
	return NewTCPProviderWID(numPeers, port, uuid.New())
}

func NewTCPProviderWID(numPeers int, port int, id uuid.UUID) *TCPProvider {
	return &TCPProvider{
		port:           port,
		numPeers:       numPeers,
		id:             id,
		peers:          make(map[uuid.UUID]*TCPConnection, numPeers),
		peerAddrs:      make(map[net.Addr]bool, numPeers),
		delivered:      make(map[uuid.UUID]map[uuid.UUID]bool),
		operations:     make(map[uuid.UUID][]byte),
		incomingOps:    make(chan []byte, 100),
		opsToBroadcast: make(chan []byte, 100),
	}
}

func (p *TCPProvider) CloseAll() {
	for _, peer := range p.peers {
		peer.conn.Close()
	}
}

// This should be called in a goroutine after appropriate setup - Setting up channels
func (p *TCPProvider) Listen() {
	address := net.JoinHostPort("::", strconv.Itoa(p.port))
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

// Broadcasts an operation to all peers
func (p *TCPProvider) HandleBroadcast() {
	for {
		opToSend := <-p.opsToBroadcast

		// Generate a new ID for the operation
		newOpId := uuid.New()

		opMsg := Message{Message: &Message_Operation{Operation: &OperationMsg{Id: newOpId[:], Op: opToSend}}}
		opData, err := proto.Marshal(&opMsg)
		if err != nil {
			log.Println("Error marshalling operation: ", err.Error())
			continue
		}

		p.AddOperation(opToSend, newOpId)
		//log.Println("Broadcasting operation:", newOpId.String())

		p.peersMu.RLock()
		for _, conn := range p.peers {
			conn.SendMsg(opData)
		}
		p.peersMu.RUnlock()
	}
}

func (p *TCPProvider) IncomingOpsChannel() chan []byte {
	return p.incomingOps
}

func (p *TCPProvider) BroadcastChannel() chan []byte {
	return p.opsToBroadcast
}

// Attempts to connect to a peer
func (p *TCPProvider) Connect(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Printf("Error resolving address: %s", err.Error())
	}
	// If we already have this peer, don't connect again
	if p.peerAddrs[tcpAddr] {
		return
	}

	p.ConnectToPeer(tcpAddr)
}

// Connects to many peers
func (p *TCPProvider) ConnectMany(addrs []string) {
	for _, addr := range addrs {
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Printf("Error resolving address: %s", err.Error())
		}
		// If we already have this peer, don't connect again
		if p.peerAddrs[tcpAddr] {
			continue
		}

		p.ConnectToPeer(tcpAddr)
	}
}

// Connects to a peer and adds it to the peer map
// And starts new goroutine to handle the connection
func (p *TCPProvider) ConnectToPeer(tcpAddr *net.TCPAddr) {
	// Check if addr is local address and port is the same as ours
	if tcpAddr.IP.IsLoopback() && tcpAddr.Port == p.port {
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Error connecting to peer: %s", err.Error())
		return
	}

	go NewTCPConnection(conn, p).handle()
}

func (p *TCPProvider) GetPeerAddrs() []net.Addr {
	p.peersMu.RLock()
	defer p.peersMu.RUnlock()

	peerAddrs := make([]net.Addr, 0, len(p.peers))
	for addr := range p.peerAddrs {
		peerAddrs = append(peerAddrs, addr)
	}

	return peerAddrs
}

// This should only be called after after peer id has been received
// Errors if the peer already exists or the peer map is full
func (p *TCPProvider) AddPeer(tcpConn *TCPConnection) error {
	p.peersMu.Lock()
	defer p.peersMu.Unlock()

	if val, ok := p.peers[tcpConn.peerId]; ok && val != nil { // Check if the peer already exists
		return errors.New("peer already exists (not nil)")
	} else if len(p.peers) == p.numPeers { // Check if the peer map is full
		return errors.New("peer map is full")
	}

	p.peerAddrs[tcpConn.conn.RemoteAddr()] = true
	p.peers[tcpConn.peerId] = tcpConn
	return nil
}

// Sets the peer in map to nil
func (p *TCPProvider) RemovePeer(tcpConn *TCPConnection) {
	p.peersMu.Lock()
	defer p.peersMu.Unlock()

	delete(p.peerAddrs, tcpConn.conn.RemoteAddr())
	p.peers[tcpConn.peerId] = nil
}

func (p *TCPProvider) AddOperation(op []byte, opId uuid.UUID) {
	p.deliveredMu.Lock()
	defer p.deliveredMu.Unlock()

	p.operations[opId] = op
	p.delivered[opId] = make(map[uuid.UUID]bool, p.numPeers-1) // Size is numPeers-1 because final peer won't store the operation in the delivered map
}

func (p *TCPProvider) AddDelivered(opId uuid.UUID, peerId uuid.UUID) {
	p.deliveredMu.Lock()
	defer p.deliveredMu.Unlock()

	// This would be the last peer to receive the operation (No need to add it to the delivered map)
	// If final peer and not in delivered map, delete entries
	if _, exists := p.delivered[opId][peerId]; len(p.delivered[opId]) == p.numPeers-1 && !exists {
		delete(p.operations, opId)
		delete(p.delivered, opId)
	} else {
		p.delivered[opId][peerId] = true
	}
}

func (p *TCPProvider) GetOperation(opId uuid.UUID) []byte {
	p.deliveredMu.RLock()
	defer p.deliveredMu.RUnlock()

	return p.operations[opId]
}
