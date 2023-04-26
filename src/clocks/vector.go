package clocks

import (
	"fmt"
	"github.com/google/uuid"
)

type VectorTimestamp map[uuid.UUID]uint64

type VectorClock struct {
	Vector VectorTimestamp
	ActID  uuid.UUID
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

	return &VectorClock{Vector: *NewVectorTimestamp(), ActID: actorID}
}

/* ----- VectorTimestamp ------ */
// will return 0 if a == b, -1 if a < b, 1 if a > b, 2 if a || b (concurrent)
// This function makes use of go's default map behaviour (if a key doesn't exist, it returns 0)
func (vt *VectorTimestamp) Compare(other *VectorTimestamp) int {
	isLess := false
	isGreater := false

	// loop through both maps and compare the values
	for key := range *vt {
		if (*vt)[key] < (*other)[key] {
			isLess = true
		} else if (*vt)[key] > (*other)[key] {
			isGreater = true
		}
	}

	for key := range *other {
		if (*vt)[key] < (*other)[key] {
			isLess = true
		} else if (*vt)[key] > (*other)[key] {
			isGreater = true
		}
	}

	switch {
	case isLess && isGreater:
		return 2
	case isLess:
		return -1
	case isGreater:
		return 1
	default:
		return 0 // The two timestamps are the same
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
	return v.ActID
}

func (v *VectorClock) Inc() *VectorTimestamp {
	v.Vector[v.ActID]++
	return v.Timestamp()
}

func (v *VectorClock) Tick() *VectorTimestamp {
	vc := v.Timestamp()
	vc.Inc(v.ActID)
	return vc
}

func (v *VectorClock) Merge(other *VectorTimestamp) {
	for key, value := range *other {
		if value > v.Vector[key] {
			v.Vector[key] = value
		}
	}
}

func (v *VectorClock) Timestamp() *VectorTimestamp {
	return v.Vector.Clone()
}

func (v *VectorClock) CompareTimestamp(other *VectorTimestamp) int {
	return v.Vector.Compare(other)
}

func (v *VectorClock) String() string {
	return fmt.Sprintf("%v: %v", v.ActID, v.Vector)
}

func (v *VectorTimestamp) String() string {
	return fmt.Sprintf("%v", *v)
}

func (v *VectorTimestamp) CausallyReady(other *VectorTimestamp) bool {
	oneLarger := false
	for key, value := range *v {
		if value <= (*other)[key] {
			continue
		} else {
			if oneLarger {
				return false
			}
			if value == (*other)[key]+1 {
				oneLarger = true
			} else {
				return false
			}
		}
	}
	return oneLarger 
}

func (v *VectorTimestamp) Same(other *VectorTimestamp) bool {
	for key, value := range *v {
		otherValue, ok := (*other)[key]
		if !ok || value != otherValue {
			return false
		}
	}
	return true
}

func (v *VectorTimestamp) Less(other *VectorTimestamp) bool {
	for key, value := range *v {
		otherValue := (*other)[key]
		if value > otherValue {
			return false
		}
	}
	return true
}