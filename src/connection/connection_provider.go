package connection 

type ConnectionProvider interface {
	broadcast(opID []byte, message []byte)
	OutputChannel() chan []byte
}