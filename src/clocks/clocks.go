package clocks

import "github.com/google/uuid"

/* ----- Interfaces ------ */
// T should be the struct implementing this interface
type TotalOrder[T any] interface {
	Compare(other T) int
}

type PartialOrder[T any] interface {
	Compare(other T) int
}

type Timestamp[T any] interface {
	Clone() T
}

// Clock with timestamp of type T
type Clock[T Timestamp[T]] interface {
	Timestamp() T // returns a clone of the current timestamp
	Inc() T       // increments clock, returns clone
	Tick() T      // returns clock inc by 1 (doesn't update clock)
	Merge(other T)
	ActorID() uuid.UUID
}
