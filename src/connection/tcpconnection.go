package connection

import (
	"encoding/binary"
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
	msg, err := c.ReadMessage(lengthBuffer, dataBuffer)
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

	// Read messages from the connection
	for {
		msg, err = c.ReadMessage(lengthBuffer, dataBuffer)
		if err != nil {
			log.Printf("Error reading message: %s", err.Error())
			return
		}

		// Read the length of the message
		// _, err := io.ReadFull(c.conn, lengthBuffer)
		// if err == io.EOF || err == io.ErrUnexpectedEOF {
		// 	log.Println("Connection closed for client:", c.id)
		// 	return
		// } else if err != nil {
		// 	log.Printf("Error reading message length: %s; for client: %v", err.Error(), c.id)
		// 	return
		// }
		// // Decode the length
		// length := binary.BigEndian.Uint32(lengthBuffer)

		// messageBuffer := dataBuffer[:length]
		// // Read the message
		// _, err = io.ReadFull(c.conn, messageBuffer)
		// if err == io.EOF || err == io.ErrUnexpectedEOF {
		// 	log.Println("Connection closed for client:", c.id)
		// 	return
		// } else if err != nil {
		// 	log.Printf("Error reading message: %s; for client: %v", err.Error(), c.id)
		// 	return
		// }

		// Handle the message
		// (Decode into a protobuf message, etc.)

	}
}

func MessageToID(msg []byte) (uuid.UUID, error) {
	peerID := &PeerID{}
	err := proto.Unmarshal(msg, peerID)
	if err != nil {
		return uuid.Nil, err
	}
	peerIDBytes := peerID.Id
	peerIDUUID, err := uuid.FromBytes(peerIDBytes)
	if err != nil {
		return uuid.Nil, err
	}
	return peerIDUUID, nil
}

func (c *TCPConnection) ReadMessage(lengthBuffer, dataBuffer []byte) ([]byte, error) {
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

	return messageBuffer, nil
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
