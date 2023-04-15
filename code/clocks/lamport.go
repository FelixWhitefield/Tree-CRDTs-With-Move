package clocks 

import (
	"fmt"
	"bytes"
	"github.com/google/uuid"
)

type Lamport struct {
	actorID uuid.UUID
	counter uint64
}

func NewLamport(ids ...uuid.UUID) *Lamport {
	var actorID uuid.UUID
	if len(ids) > 0 {
		actorID = ids[0]
	} else {
		actorID = uuid.New()
	}

	return &Lamport{actorID: actorID, counter: 0}
}

func (l Lamport) ActorID() uuid.UUID {
	return l.actorID
}

// will return 0 if a == b, -1 if a < b, 1 if a > b
func (l Lamport) Compare(other Lamport) int {
	switch {
	case l.counter < other.counter: return -1
	case l.counter > other.counter: return 1
	default: return bytes.Compare(l.actorID[:], other.actorID[:]) 
	}
}

func (l *Lamport) Inc() Lamport {
	l.counter++
	return Lamport{actorID: l.actorID, counter: l.counter}
}

func (l Lamport) Tick() Lamport {
	l.counter++
	return l
}

func (l *Lamport) Merge(other Lamport) {
	if other.counter > l.counter {
		l.counter = other.counter
	} 
}

func (l Lamport) Clone() Lamport {
	return Lamport{actorID: l.actorID, counter: l.counter}
}

func (l Lamport) String() string {
	return fmt.Sprintf("%v: %v", l.actorID, l.counter)
}