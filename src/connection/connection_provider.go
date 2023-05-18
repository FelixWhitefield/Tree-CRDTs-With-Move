package connection

import "net"

type ConnectionProvider interface {
	BroadcastChannel() chan []byte
	IncomingOpsChannel() chan []byte
	Connect(addr string)
	HandleBroadcast()
	Listen()
	NumPeers() int
	CloseAll()
	GetPeerAddrs() []net.Addr
}
