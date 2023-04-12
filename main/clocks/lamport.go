package clocks 

import (
	"fmt"
	"math/rand"
)

type Lamport struct {
	actorID int
	counter int
}

func NewLamport(ids ...int) *Lamport {
	var actorID int
	if len(ids) > 0 {
		actorID = ids[0]
	} else {
		actorID = rand.Int()
	}

	return &Lamport{actorID: actorID, counter: 0}
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