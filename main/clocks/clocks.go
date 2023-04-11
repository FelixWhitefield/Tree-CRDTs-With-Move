package clocks

import (
	"fmt"
)

const (
	LESS = -1
	EQUAL = 0
	GREATER = 1
	CONCURRENT = 2
)


type TotalOrder interface {
	Compare(other TotalOrder) int 
}


type PartialOrder interface {
	Compare(other PartialOrder) int
}

type Timestamp interface {
	Clone() Timestamp
}

type Clock interface {
	Inc() Clock // increments clock, returns clone
	Tick() Clock  // returns clock inc by 1 (doesn't update clock)
	Merge(other Timestamp)
}

type Lamport struct {
	actorID int
	counter int
}

type VectorTimestamp map[int]int


type VectorClock struct {
	timestamp VectorTimestamp
	actorID int
}

func NewLamport(actorID int) *Lamport {
	return &Lamport{actorID: actorID, counter: 0}
}

func NewVectorTimestamp() VectorTimestamp {
	return make(map[int]int)
}

func NewVectorClock(actorID int) VectorClock {
	return VectorClock{actorID: actorID, timestamp: NewVectorTimestamp()}
}

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

func (l Lamport) Print() {
	fmt.Println(l.actorID , ": ", l.counter)
}

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
			return 2 // concurrent
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
		return -1
	} else if isGreater {
		return 1
	} else {
		return 0
	}
}

func (vt VectorTimestamp) Inc(key int) {
	vt[key]++
}

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

func (v VectorClock) Clone() VectorClock {
	return VectorClock{actorID: v.actorID, timestamp: v.timestamp}
}

func (v VectorClock) String() string {
	return fmt.Sprintf("%v: %v", v.actorID, v.timestamp)
}


