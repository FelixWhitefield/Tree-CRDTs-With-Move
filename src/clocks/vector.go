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

func NewVectorTimestamp() *VectorTimestamp {
	timestamp := make(VectorTimestamp)
	return &timestamp
}
// Returns a new VectorClock with a random actorID or the given actorID 
func NewVectorClock(ids ...uuid.UUID) *VectorClock {
	var actorID uuid.UUID
	if len(ids) > 0 {
		actorID = ids[0]
	} else {
		actorID = uuid.New()
	}

	return &VectorClock{timestamp: *NewVectorTimestamp(), actorID: actorID}
}

/* ----- VectorTimestamp ------ */
// will return 0 if a == b, -1 if a < b, 1 if a > b, 2 if a || b (concurrent)
// This function makes use of go's default map behaviour (if a key doesn't exist, it returns 0)
func (vt *VectorTimestamp) Compare(other *VectorTimestamp) int {
	isLess := false
	isGreater := false

	for key, value := range *vt {
		otherValue := (*other)[key]

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
		for key := range *other {
			if _, exists := (*vt)[key]; !exists {
				isGreater = true
				break
			}
		}
	}

	switch {
	case isLess: return -1
	case isGreater: return 1
	default: return 0 // The two timestamps are the same
	}
}

func (vt *VectorTimestamp) Inc(id uuid.UUID) {
	(*vt)[id]++
}

func (vt *VectorTimestamp) Clone() *VectorTimestamp {
	newTimestamp := NewVectorTimestamp()
	for key, value := range *vt {
		(*newTimestamp)[key] = value
	}
	return newTimestamp
}

/* ----- VectorClock ------ */
func (v *VectorClock) ActorID() uuid.UUID {
	return v.actorID
}

func (v *VectorClock) Inc() *VectorTimestamp {
	v.timestamp[v.actorID]++
	return v.CloneTimestamp()
}

func (v *VectorClock) Tick() *VectorTimestamp {
	vc := v.CloneTimestamp()
	vc.Inc(v.actorID)
	return vc
}

func (v *VectorClock) Merge(other *VectorTimestamp) {
	for key, value := range *other {
		if value > v.timestamp[key] {
			v.timestamp[key] = value
		}
	}
}

func (v *VectorClock) CloneTimestamp() *VectorTimestamp {
	return v.timestamp.Clone()
}

func (v *VectorClock) CompareTimestamp(other *VectorTimestamp) int {
	return v.timestamp.Compare(other)
}

func (v *VectorClock) String() string {
	return fmt.Sprintf("%v: %v", v.actorID, v.timestamp)
}

