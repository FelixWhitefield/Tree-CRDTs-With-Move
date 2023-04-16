package clocks

import "github.com/google/uuid"

/* ----- Interfaces ------ */
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
	Inc() T // increments clock, returns clone
	Tick() T  // returns clock inc by 1 (doesn't update clock)
	Merge(other T)
	ActorID() uuid.UUID
}