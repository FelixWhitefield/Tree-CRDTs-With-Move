package connection

// This is the TCP connection handler
//
// This represents a connection to a peer, and handles the receiving of messages

import (
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
)

const (
	MAX_MSG_SIZE = 1024 * 10 // 10KB
)

type TCPConnection struct {
	conn             net.Conn
	remoteListenAddr net.Addr
	peerId           uuid.UUID
	tcpProv          *TCPProvider
}

func NewTCPConnection(conn net.Conn, p *TCPProvider) *TCPConnection {
	return &TCPConnection{conn: conn, tcpProv: p}
}

func (c *TCPConnection) handle() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Connection Panicked:", err)
		}
	}()
	defer c.conn.Close()

	// Send the ID to the client ----------------------------------
	peerIdMsg := &Message{Message: &Message_PeerID{PeerID: &PeerID{Id: c.tcpProv.id[:], Addr: c.tcpProv.laddr.String()}}}
	peerIDBytes, err := proto.Marshal(peerIdMsg)
	if err != nil {
		log.Printf("Error marshalling peer ID: %s", err.Error())
		return
	}
	c.SendMsg(peerIDBytes)
	// -------------------------------------------------------------

	// Create buffer outside of loop to avoid creating a new one each time (Reusing the buffers)
	lengthBuffer := make([]byte, 4)
	dataBuffer := make([]byte, MAX_MSG_SIZE)

	// First message should be the ID of the peer
	msg, err := c.readConnMsg(lengthBuffer, dataBuffer)
	if err != nil {
		log.Printf("Error reading initial message: %s", err.Error())
		return
	}

	newPeerId, remoteListenStr, err := MessageToID(msg) // Format the message into a UUID
	c.peerId = newPeerId
	if err != nil {
		log.Printf("Error reading peer ID: %s", err.Error())
		return
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", remoteListenStr)
	if err != nil {
		log.Printf("Error resolving remote address: %s", err.Error())
		return
	}

	c.remoteListenAddr = remoteAddr
	err = c.tcpProv.addPeer(c) // Add the peer to the list of peers
	if err != nil {
		log.Printf("Error adding peer: %s", err.Error())
		return
	}
	log.Printf("Connected to: %v; at %v", c.peerId.String(), c.conn.RemoteAddr())
	defer c.tcpProv.removePeer(c) // Remove the peer from the list of peers when the connection is closed

	c.tcpProv.sendMissingOps(c.peerId) // Send any missing operations to the new peer
	// Read messages from the connection
	for {
		msg, err = c.readConnMsg(lengthBuffer, dataBuffer)
		if err == io.EOF || errors.Is(err, net.ErrClosed) {
			log.Printf("Connection closed by peer: %s", c.peerId.String())
			return
		}
		if err != nil {
			log.Printf("Malformed message: %s. Connection closed for %s", err.Error(), c.peerId.String())
			return
		}
		// Handle the message
		switch msg.Message.(type) {
		case *Message_PeerAddresses:
			//log.Println("Received peer addresses")
			// connect to peers who are not already connected
			peers := msg.GetPeerAddresses().PeerAddrs
			go c.tcpProv.ConnectMany(peers)
		case *Message_Operation:
			//log.Println("Received Operation")
			opMsg := msg.GetOperation()
			opAck := &Message{Message: &Message_OperationAck{OperationAck: &OperationAck{Id: opMsg.GetId(), Ack: true}}}
			opAckBytes, err := proto.Marshal(opAck)
			if err != nil {
				log.Printf("Error marshalling operation ack: %s", err.Error())
			}

			c.SendMsg(opAckBytes) // Send the operation ack to the client
			c.tcpProv.incomingOps <- opMsg.GetOp()
		case *Message_OperationAck:
			//log.Println("Received Op ACK")
			opAck := msg.GetOperationAck()
			ackId, err := uuid.FromBytes(opAck.GetId())
			if err != nil {
				log.Printf("Error converting ack ID to UUID: %s", err.Error())
				continue
			}
			if opAck.GetAck() {
				c.tcpProv.addDelivered(ackId, c.peerId)
			}
		default:
			log.Printf("Unknown message type: %s", msg.String())
		}
	}
}

func (c *TCPConnection) SharePeers() {
	peerAddrs := c.tcpProv.GetPeerAddrs()
	peerAddrsStr := make([]string, 0, len(peerAddrs))

	for _, addr := range peerAddrs {
		if addr == c.remoteListenAddr { // Don't send the peer their own address
			continue
		}
		peerAddrsStr = append(peerAddrsStr, addr.String())
	}

	if len(peerAddrsStr) == 0 { // Don't send an empty list of peers
		return
	}

	peerAddrsMsg := &Message{Message: &Message_PeerAddresses{PeerAddresses: &PeerAddresses{PeerAddrs: peerAddrsStr}}}
	peerAddrsBytes, err := proto.Marshal(peerAddrsMsg)
	if err != nil {
		log.Printf("Error marshalling peer addresses: %s", err.Error())
		return
	}
	//log.Println("Sharing peers")
	c.SendMsg(peerAddrsBytes)
}

// Takes a protobuf message and converts it into a uuid
// This is used to convert the ID of a peer into a uuid
func MessageToID(msg *Message) (uuid.UUID, string, error) {
	peerId, ok := msg.Message.(*Message_PeerID)
	if !ok {
		return uuid.Nil, "", errors.New("Message is not a PeerID")
	}
	peerID := peerId.PeerID

	peerIDBytes := peerID.Id
	peerIDUUID, err := uuid.FromBytes(peerIDBytes)
	if err != nil {
		return uuid.Nil, "", err
	}
	peerAddr := peerID.Addr
	return peerIDUUID, peerAddr, nil
}

// Reads in a message from the connection
// The message is prefixed with a 4 byte length header
// The message is then read in and unmarshalled
// This function takes in two buffers to be used for reading in the message
func (c *TCPConnection) readConnMsg(lengthBuffer, dataBuffer []byte) (*Message, error) {
	_, err := io.ReadFull(c.conn, lengthBuffer)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	// Decode the length
	length := binary.BigEndian.Uint32(lengthBuffer)

	messageBuffer := dataBuffer[:length]
	// Read the message
	_, err = io.ReadFull(c.conn, messageBuffer)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	message := &Message{}
	err = proto.Unmarshal(messageBuffer, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Sends a byte message to the connection
// The message is prefixed with a 4 byte length header
// which describes the length of the message
// This could use a buffer pool to reduce allocations (TODO)
func (c *TCPConnection) SendMsg(data []byte) {
	if c == nil {
		return
	}
	message := make([]byte, 4+len(data)) // 4 bytes for length (up to 4GB, max length for protobuf)
	binary.BigEndian.PutUint32(message, uint32(len(data)))

	//Write the length and then the data, using a single write call (to ensure they are sent together)
	copy(message[4:], data)

	_, err := c.conn.Write(message)
	if err != nil {
		log.Printf("Error: %s; sending message to client: %v", err.Error(), c.peerId)
	}
}

func (c *TCPConnection) SendOperation(encodedOp []byte) {
	if c == nil {
		return
	}
	_, err := c.conn.Write(encodedOp)
	if err != nil {
		log.Printf("Error: %s; sending message to client: %v", err.Error(), c.peerId)
	}
}
