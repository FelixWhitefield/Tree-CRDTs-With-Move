package connection

import (
	"testing"
	"time"
)

func TestTCPProviderConnecting(t *testing.T) {
	tcpProvider1 := NewTCPProvider(2, 1511)
	go tcpProvider1.Listen()
	tcpProvider2 := NewTCPProvider(2, 1512)
	go tcpProvider2.Listen()

	time.Sleep(200 * time.Millisecond) // Time for nodes to start listening

	tcpProvider1.Connect("localhost:1512")

	time.Sleep(1 * time.Second) // Time for nodes to connect

	// Test two peers can connect
	if len(tcpProvider1.peers) != 1 && len(tcpProvider2.peers) != 1 {
		t.Errorf("Expected 1 peer, got %d and %d", len(tcpProvider1.peers), len(tcpProvider2.peers))
	}

	// Test peer sharing
	tcpProvider3 := NewTCPProvider(2, 1513)
	go tcpProvider3.Listen()

	time.Sleep(200 * time.Millisecond) // Time for nodes to start listening

	tcpProvider1.Connect("localhost:1513")

	time.Sleep(1 * time.Second) // Time for nodes to connect

	if len(tcpProvider1.peers) != 2 && len(tcpProvider2.peers) != 2 && len(tcpProvider3.peers) != 2 {
		t.Errorf("Expected 2 peers, got %d, %d and %d", len(tcpProvider1.peers), len(tcpProvider2.peers), len(tcpProvider3.peers))
	}
}
