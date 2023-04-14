package clocks

import (
	"fmt"
	"github.com/google/uuid"
)

type VectorTimestamp map[uuid.UUID]uint64

type VectorClock struct {
	timestamp VectorTimestamp
	actorID uuid.UUID
}

func NewVectorTimestamp() VectorTimestamp {
	return make(map[uuid.UUID]uint64)
}
// Returns a new VectorClock with a random actorID or the given actorID 
func NewVectorClock(ids ...uuid.UUID) VectorClock {
	var actorID uuid.UUID
	if len(ids) > 0 {
		actorID = ids[0]
	} else {
		actorID = uuid.New()
	}

	return VectorClock{actorID: actorID, timestamp: NewVectorTimestamp()}
}

/* ----- VectorTimestamp ------ */
// Returns either LESS, EQUAL, GREATER or CONCURRENT
// This function makes use of go's default map behaviour (if a key doesn't exist, it returns 0)
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

	switch {
	case isLess: return LESS
	case isGreater: return GREATER
	default: return EQUAL // The two timestamps are the same
	}
}

func (vt VectorTimestamp) Inc(id uuid.UUID) {
	vt[id]++
}

/* ----- VectorClock ------ */
func (v VectorClock) ActorID() uuid.UUID {
	return v.actorID
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

