package connection

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
	conn    net.Conn
	id      uuid.UUID
	tcpProv *TCPProvider
}

func NewTCPConnection(conn net.Conn, p *TCPProvider) *TCPConnection {
	return &TCPConnection{conn: conn, tcpProv: p}
}

func (c *TCPConnection) handle() {
	defer c.conn.Close()   

	// Send the ID to the client ----------------------------------
	peerID := &PeerID{Id: c.tcpProv.id[:]}
	peerIDBytes, err := proto.Marshal(peerID)
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
	msg, err := c.ReadConnMsg(lengthBuffer, dataBuffer)
	if err != nil {
		log.Printf("Error reading initial message: %s", err.Error())
		return
	}

	peerIDUUID, err := MessageToID(msg) // Format the message into a UUID
	if err != nil {
		log.Printf("Error reading peer ID: %s", err.Error())
		return
	}
	c.id = peerIDUUID // Add the ID to the connection

	err = c.tcpProv.AddPeer(peerIDUUID, c) // Add the peer to the list of peers
	if err != nil {
		log.Printf("Error adding peer: %s", err.Error())
		return
	}
	defer c.tcpProv.RemovePeer(peerIDUUID) // Remove the peer from the list of peers when the connection is closed

	// Read messages from the connection
	for {
		msg, err = c.ReadConnMsg(lengthBuffer, dataBuffer)
		if err != nil {
			log.Printf("Error reading message: %s", err.Error())
			return
		}

		// Handle the message
		switch msg.Message.(type) {
		case *Message_PeerAddresses:
			// connect to peers who are not already connected
		case *Message_Operation:
			// Add the operation to the list of incoming operations
		}

	}
}

// Takes a protobuf message and converts it into a uuid
// This is used to convert the ID of a peer into a uuid
func MessageToID(msg *Message) (uuid.UUID, error) {
	peerId, ok := msg.Message.(*Message_PeerID)
	if !ok {
		return uuid.Nil, errors.New("Message is not a PeerID")
	}
	peerID := peerId.PeerID

	peerIDBytes := peerID.Id
	peerIDUUID, err := uuid.FromBytes(peerIDBytes)
	if err != nil {
		return uuid.Nil, err
	}
	return peerIDUUID, nil
}

// Reads in a message from the connection
// The message is prefixed with a 4 byte length header
// The message is then read in and unmarshalled
// This function takes in two buffers to be used for reading in the message
func (c *TCPConnection) ReadConnMsg(lengthBuffer, dataBuffer []byte) (*Message, error) {
	_, err := io.ReadFull(c.conn, lengthBuffer)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		log.Println("Connection closed for client:", c.id)
		return nil, err
	} else if err != nil {
		log.Printf("Error reading message length: %s; for client: %v", err.Error(), c.id)
		return nil, err
	}
	// Decode the length
	length := binary.BigEndian.Uint32(lengthBuffer)

	messageBuffer := dataBuffer[:length]
	// Read the message
	_, err = io.ReadFull(c.conn, messageBuffer)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		log.Println("Connection closed for client:", c.id)
		return nil, err
	} else if err != nil {
		log.Printf("Error reading message: %s; for client: %v", err.Error(), c.id)
		return nil, err
	}

	message := &Message{}
	err = proto.Unmarshal(messageBuffer, message)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err.Error())
		return nil, err
	}

	return message, nil
}

// Sends a byte message to the connection
// The message is prefixed with a 4 byte length header
// which describes the length of the message
func (c *TCPConnection) SendMsg(data []byte) {
	length := make([]byte, 4) // 4 bytes for length (up to 4GB, max length for protobuf)
	binary.BigEndian.PutUint32(length, uint32(len(data)))

	// Write the length and then the data, using a single write call (to ensure they are sent together)
	_, err := c.conn.Write(append(length, data...))
	if err != nil {
		log.Printf("Error: %s; sending message to client: %v", err.Error(), c.id)
	}
}
