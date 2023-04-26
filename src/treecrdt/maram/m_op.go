package maram

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

// Represents a move operation
type OpMove[MD any, T opTimestamp[T]] struct {
	Timestmp T
	ChldID   uuid.UUID
	NewP     *TreeNode[MD]
	Priotity c.Lamport
}

func (op *OpMove[MD, T]) Timestamp() T {
	return op.Timestmp.Clone()
}

func (op *OpMove[MD, T]) CompareOp(other *OpMove[MD, T]) int {
	return op.Timestmp.Compare(other.Timestmp)
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
