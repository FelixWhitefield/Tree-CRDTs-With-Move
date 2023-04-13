package clocks

const (
	LESS = -1
	EQUAL = 0
	GREATER = 1
	CONCURRENT = 2
)
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

type Clock[T any, U any] interface {
	Inc() T // increments clock, returns clone
	Tick() T  // returns clock inc by 1 (doesn't update clock)
	Merge(other U)
}