package clocks

import (
	"fmt"
	"math/rand"
)

const (
	LESS = -1
	EQUAL = 0
	GREATER = 1
	CONCURRENT = 2
)
/* ----- Interfaces ------ */
type TotalOrder interface {
	Compare(other TotalOrder) int 
}

type PartialOrder interface {
	Compare(other any) int
}

type Timestamp[T any] interface {
	Clone() T
}

type Clock[T any, K any] interface {
	Inc() T // increments clock, returns clone
	Tick() T  // returns clock inc by 1 (doesn't update clock)
	Merge(other K)
}
/* ----- Structs ------ */
type Lamport struct {
	actorID int
	counter int
}

type VectorTimestamp map[int]int

type VectorClock struct {
	timestamp VectorTimestamp
	actorID int
}
/* ----- New Methods ------ */
// Returns a new Lamport with a random actorID or the given actorID
func NewLamport(ids ...int) *Lamport {
	var actorID int
	if len(ids) > 0 {
		actorID = ids[0]
	} else {
		actorID = rand.Int()
	}

	return &Lamport{actorID: actorID, counter: 0}
}

func NewVectorTimestamp() VectorTimestamp {
	return make(map[int]int)
}
// Returns a new VectorClock with a random actorID or the given actorID 
func NewVectorClock(ids ...int) VectorClock {
	var actorID int
	if len(ids) > 0 {
		actorID = ids[0]
	} else {
		actorID = rand.Int()
	}

	return VectorClock{actorID: actorID, timestamp: NewVectorTimestamp()}
}
/* ----- Lamport Methods ------ */
// Returns either LESS, EQUAL, or GREATER
func (l Lamport) Compare(other Lamport) int {
	if diff := l.counter - other.counter; diff < 0 {
		return LESS
	} else if diff > 0 {
		return GREATER
	} 
	return EQUAL
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

/* ----- VectorTimestamp Methods ------ */
// Returns either LESS, EQUAL, GREATER or CONCURRENT
func (vt VectorTimestamp) Compare(other VectorTimestamp) int {
	isLess := false
	isGreater := false

	for key, value := range vt {
		otherValue := other[key]

		if value < otherValue {
			isLess = true
		} else if value > otherValue {
			isGreater = true
		}

		if isLess && isGreater {
			return CONCURRENT // concurrent
		}
	}

	if !isLess && !isGreater {
		for key := range other {
			if _, exists := vt[key]; !exists {
				isGreater = true
				break
			}
		}
	}

	if isLess {
		return LESS
	} else if isGreater {
		return GREATER
	} else {
		return EQUAL
	}
}

func (vt VectorTimestamp) Inc(key int) {
	vt[key]++
}

/* ----- VectorClock Methods ------ */
func (v *VectorClock) Inc() VectorClock {
	v.timestamp[v.actorID]++
	return VectorClock{actorID: v.actorID, timestamp: v.timestamp}
}

func (v VectorClock) Tick() VectorClock {
	v.timestamp[v.actorID]++
	return v
}

func (v *VectorClock) Merge(other VectorTimestamp) {
	for key, value := range other {
		if value > v.timestamp[key] {
			v.timestamp[key] = value
		}
	}
}

func (v VectorClock) CloneTimestamp() VectorTimestamp {
	timestamp := NewVectorTimestamp()
	for key, value := range v.timestamp {
		timestamp[key] = value
	}
	return timestamp
}

func (v VectorClock) CompareTimestamp(other VectorTimestamp) int {
	return v.timestamp.Compare(other)
}

func (v VectorClock) String() string {
	return fmt.Sprintf("%v: %v", v.actorID, v.timestamp)
}

