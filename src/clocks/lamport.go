package clocks

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
)

type Lamport struct {
	ActID   uuid.UUID
	Counter uint64
}

func NewLamport(ids ...uuid.UUID) *Lamport {
	var actorID uuid.UUID
	if len(ids) > 0 {
		actorID = ids[0]
	} else {
		actorID = uuid.New()
	}

	return &Lamport{ActID: actorID, Counter: 0}
}

func (l *Lamport) ActorID() uuid.UUID {
	return l.ActID
}

// will return 0 if a == b, -1 if a < b, 1 if a > b
func (l *Lamport) Compare(other *Lamport) int {
	switch {
	case l.Counter < other.Counter:
		return -1
	case l.Counter > other.Counter:
		return 1
	default:
		return bytes.Compare(l.ActID[:], other.ActID[:])
	}
}

// increments the counter and returns a new Lamport clock with the same actorID
func (l *Lamport) Inc() *Lamport {
	l.Counter++
	return &Lamport{ActID: l.ActID, Counter: l.Counter}
}

// returns a new Lamport clock with the same actorID and a counter incremented by 1 (doesn't update clock)
func (l *Lamport) Tick() *Lamport {
	return &Lamport{ActID: l.ActID, Counter: l.Counter + 1}
}

func (l *Lamport) Merge(other *Lamport) {
	if other.Counter > l.Counter {
		l.Counter = other.Counter
	}
}

func (l *Lamport) Clone() *Lamport {
	return &Lamport{ActID: l.ActID, Counter: l.Counter}
}

func (l *Lamport) Timestamp() *Lamport {
	return l.Clone()
}

func (l *Lamport) String() string {
	return fmt.Sprintf("%v: %v", l.ActID, l.Counter)
}
