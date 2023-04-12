package clocks 

import (
	"fmt"
	"math/rand"
)

type Lamport struct {
	actorID uint64
	counter uint
}

func NewLamport(ids ...uint64) *Lamport {
	var actorID uint64
	if len(ids) > 0 {
		actorID = ids[0]
	} else {
		actorID = rand.Uint64()
	}

	return &Lamport{actorID: actorID, counter: 0}
}

func (l Lamport) ActorID() uint64 {
	return l.actorID
}

// Returns either LESS, EQUAL, or GREATER
func (l Lamport) Compare(other Lamport) int {
	switch {
	case l.counter < other.counter: return LESS
	case l.counter > other.counter: return GREATER
	default: 
		switch {
		case l.actorID < other.actorID: return LESS
		case l.actorID > other.actorID: return GREATER
		default: return EQUAL // The two timestamps are the same
		}
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