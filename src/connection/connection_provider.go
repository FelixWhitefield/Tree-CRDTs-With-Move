package connection 

type ConnectionProvider interface {
	BroadcastChannel() chan []byte
	IncomingOpsChannel() chan []byte
}