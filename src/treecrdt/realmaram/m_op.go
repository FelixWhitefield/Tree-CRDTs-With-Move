package realmaram

import (
	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

type opTimestamp[T any] interface {
	c.PartialOrder[T]
	c.Timestamp[T]
}

// Main operation interface
type Operation[T opTimestamp[T]] interface {
	Timestamp() T
}

// Move operations can be either up or down moves
// Each move operation will have a priority which
// is used to determine which move operation should
// in certain conflicts
type OpMove[MD any, T opTimestamp[T]] struct {
	Timestmp T
	ChldID   uuid.UUID
	NewP     *TreeNode[MD]
	Priority c.Lamport
	DownMove bool
}


func (op OpMove[MD, T]) Timestamp() T {
	return op.Timestmp.Clone()
}

func (op OpMove[MD, T]) IsConcurrent(other OpMove[MD, T]) bool {
	return op.Timestmp.Compare(other.Timestmp) == 2
}

func (op OpMove[MD, T]) ComparePriority(other OpMove[MD, T]) int {
	return op.Priority.Compare(&other.Priority)
}


// Represents an add operation
type OpAdd[MD any, T opTimestamp[T]] struct {
	Timestmp T
	ChldID   uuid.UUID
	NewP     *TreeNode[MD]
}

func (op *OpAdd[MD, T]) Timestamp() T {
	return op.Timestmp.Clone()
}

// Represents a remove operation
type OpRemove[T opTimestamp[T]] struct {
	Timestmp T
	ChldID   uuid.UUID
}

func (op *OpRemove[T]) Timestamp() T {
	return op.Timestmp.Clone()
}
