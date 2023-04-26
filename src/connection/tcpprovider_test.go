package connection 

import (
	"testing"
	"time"
)

func TestTCPProviderConnecting(t *testing.T) {
	tcpProvider1 := NewTCPProvider(1, 5555)
	tcpProvider2 := NewTCPProvider(1, 5556)

	tcpProvider1.Connect("localhost:5556")

	


}